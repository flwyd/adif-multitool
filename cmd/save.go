// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/flwyd/adif-multitool/adif"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var Save = Command{Name: "save", Run: runSave, Help: helpSave,
	Description: "Save standard input to file(s) with format inferred by extension"}

type SaveContext struct {
	OverwriteExisting bool
	WriteIfEmpty      bool
	CreateDirectory   bool
	Quiet             bool
	ShardFileCount    int
	ShardMaxRecords   int
}

func (c *SaveContext) wantShards() bool {
	return c.ShardFileCount > 0 || c.ShardMaxRecords > 0
}

func (c *SaveContext) shardEnds(records int) ([]int, error) {
	if c.ShardFileCount <= 0 && c.ShardMaxRecords <= 0 {
		return []int{records}, nil
	}
	if records <= 0 {
		if c.ShardFileCount <= 0 {
			return []int{0}, nil
		}
		return make([]int, c.ShardFileCount), nil
	}
	var num int
	if c.ShardFileCount > 0 {
		if c.ShardMaxRecords > 0 && records > c.ShardMaxRecords*c.ShardFileCount {
			return nil, fmt.Errorf("record count %d with %d shards would put more than %d records in a file",
				records, c.ShardFileCount, c.ShardMaxRecords)
		}
		num = c.ShardFileCount
	} else {
		num = records / c.ShardMaxRecords
		if records%c.ShardMaxRecords != 0 {
			num++
		}
	}
	size := records / num
	res := make([]int, num)
	prev := 0
	for i := 0; i < num; i++ {
		cur := size
		if records%num > i {
			cur++
		}
		res[i] = prev + cur
		prev = res[i]
	}
	return res, nil
}

func helpSave() string {
	return `Unless options are set explicitly, existing files will not be overwritten and
logfiles without any records will not be saved (useful if validate failed).

File name may be a template with {FIELD} placeholders replaced by field values.
For example, '{QSO_DATE}_{BAND}.adi' will create a separate file for each
contact date + band combination.  Quote the name to avoid shell expansion.

If shard options are set, records will be split across multiple files with a
file number and total file count in the filename, e.g. "mylog+1-of-3.adi".
The number of records in each shard is unspecified and may change.
`
}

func runSave(ctx *Context, args []string) error {
	cctx := ctx.CommandCtx.(*SaveContext)
	if len(args) != 1 {
		return fmt.Errorf("save expects 1 output file or template, got %v", args)
	}
	fname := args[0]
	fs := ctx.fs
	format := ctx.OutputFormat
	if !format.IsValid() {
		f, err := adif.GuessFormatFromName(fname)
		if err != nil {
			if strings.ToLower(path.Ext(fname)) == "adif" {
				f = adif.FormatADI
			} else {
				return fmt.Errorf("unknown output format, set --output: %w", err)
			}
		}
		format = f
	}
	st := newSaveTemplate(fname)
	if fs == nil {
		fs = osFilesystem{}
	}
	l, err := readFile(ctx, os.Stdin.Name())
	if err != nil {
		return err
	}
	if len(ctx.FieldOrder) > 0 {
		ro := l.FieldOrder
		l.FieldOrder = slices.Clone(ctx.FieldOrder)
		updateFieldOrder(l, ro)
	}
	for _, u := range ctx.UserdefFields {
		l.AddUserdef(u)
	}

	saveLog := func(l *adif.Logfile, file string) error {
		if !cctx.OverwriteExisting && fs.Exists(file) {
			return fmt.Errorf("output file %s already exists", file)
		}
		dir := path.Dir(file)
		if cctx.CreateDirectory {
			if err := fs.MkdirAll(dir); err != nil && !errors.Is(err, os.ErrExist) {
				return err
			}
		}
		out, err := fs.Create(file)
		if err != nil {
			if !fs.Exists(dir) {
				return fmt.Errorf("missing directory %s, use the create-dirs option: %w", dir, err)
			}
			return err
		}
		defer out.Close()
		ctx.Out = out
		ctx.OutputFormat = format
		err = write(ctx, l)
		if err == nil && !cctx.Quiet {
			fmt.Fprintf(os.Stderr, "Wrote %d records to %s\n", len(l.Records), file)
		}
		return err
	}

	logs := make(map[string]*accumulator)
	if len(l.Records) == 0 {
		if !cctx.WriteIfEmpty {
			return fmt.Errorf("no records in input, not saving to %s", fname)
		}
		if !cctx.Quiet {
			fmt.Fprintf(os.Stderr, "Warning: saving %s with no records\n", fname)
		}
		emptyfile := st.format(adif.NewRecord())
		a, err := newAccumulator(ctx)
		if err != nil {
			return err
		}
		if err := a.initFromHeader(l); err != nil {
			return err
		}
		a.Out.Filename = emptyfile
		logs[emptyfile] = a
	}
	for _, r := range l.Records {
		file := st.format(r)
		if logs[file] == nil {
			if !cctx.OverwriteExisting && fs.Exists(file) {
				return fmt.Errorf("output file %s already exists", file)
			}
			if cctx.CreateDirectory {
				dir := path.Dir(file)
				if err := fs.MkdirAll(dir); err != nil && !errors.Is(err, os.ErrExist) {
					return err
				}
			}
			a, err := newAccumulator(ctx)
			if err != nil {
				return err
			}
			if err := a.initFromHeader(l); err != nil {
				return err
			}
			a.Out.Filename = file
			logs[file] = a
		}
		logs[file].Out.AddRecord(r)
	}

	if cctx.wantShards() {
		sharded := make(map[string]*accumulator)
		for k, v := range logs {
			ends, err := cctx.shardEnds(len(v.Out.Records))
			if err != nil {
				return fmt.Errorf("could not split %s: %v", k, err)
			}
			ext := filepath.Ext(k)
			base := k[:len(k)-len(ext)]
			glob := base + "+*-of-*" + ext
			matches, err := filepath.Glob(glob)
			if err != nil {
				return fmt.Errorf("could not scan for existing files matching %s: %v", glob, err)
			}
			for _, m := range matches {
				b := filepath.Base(m)
				b = b[:len(b)-len(ext)]
				s := strings.Split(b, "+")
				var mi, mcount int
				if n, err := fmt.Sscanf(s[len(s)-1], "%d-of-%d", &mi, &mcount); err == nil && n == 2 {
					if mcount != len(ends) {
						return fmt.Errorf("existing file with a different shard count; delete matching files or save with a different name: %s", m)
					}
				}
			}
			cstr := fmt.Sprintf("%d", len(ends))
			ifmt := fmt.Sprintf("%%0%dd", len(cstr))
			for i, end := range ends {
				istr := fmt.Sprintf(ifmt, i+1)
				fname := fmt.Sprintf("%s+%s-of-%s%s", base, istr, cstr, ext)
				a, err := newAccumulator(ctx)
				if err != nil {
					return err
				}
				if err := a.initFromHeader(l); err != nil {
					return err
				}
				a.Out.Filename = fname
				sharded[fname] = a
				prev := 0
				if i > 0 {
					prev = ends[i-1]
				}
				a.Out.Records = v.Out.Records[prev:end]
			}
		}
		logs = sharded
	}

	for _, a := range logs {
		if err := a.prepare(); err != nil {
			return err
		}
	}
	errs := make([]error, len(logs))
	files := maps.Keys(logs)
	sort.Strings(files)
	for i, f := range files {
		errs[i] = saveLog(logs[f].Out, f)
	}
	return errorsJoin(errs...)
}

type saveTemplate struct {
	pieces []func(r *adif.Record) string
	static bool
}

func (t saveTemplate) format(r *adif.Record) string {
	var s strings.Builder
	for _, p := range t.pieces {
		s.WriteString(p(r))
	}
	return s.String()
}

var templateFieldPat = regexp.MustCompile(`\{\w+\}`)

func newSaveTemplate(s string) saveTemplate {
	fields := templateFieldPat.FindAllString(s, -1)
	if len(fields) == 0 {
		return saveTemplate{
			pieces: []func(*adif.Record) string{
				func(_ *adif.Record) string { return s },
			},
		}
	}
	for i, f := range fields {
		fields[i] = f[1 : len(f)-1] // strip braces
	}
	t := saveTemplate{
		pieces: make([]func(r *adif.Record) string, len(fields)*2+1),
	}
	literals := templateFieldPat.Split(s, -1)
	for i, field := range fields {
		f := field
		l := literals[i]
		t.pieces[i*2] = func(r *adif.Record) string { return l }
		t.pieces[i*2+1] = func(r *adif.Record) string {
			ff, _ := r.Get(f)
			v := strings.Map(func(c rune) rune {
				if c == ' ' || !unicode.IsPrint(c) {
					return '_'
				}
				if c == '-' || c == '_' || unicode.IsDigit(c) {
					return c
				}
				if unicode.IsLetter(c) {
					return unicode.ToUpper(c)
				}
				return '-' // replace marks that are awkward in filenames
			}, ff.Value)
			if v == "" {
				return strings.ToUpper(f) + "-EMPTY"
			}
			return strings.ToUpper(v)
		}
	}
	l := literals[len(literals)-1]
	t.pieces[len(fields)*2] = func(r *adif.Record) string { return l }
	return t
}

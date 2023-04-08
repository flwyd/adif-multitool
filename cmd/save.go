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
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/flwyd/adif-multitool/adif"
	"golang.org/x/exp/maps"
)

var Save = Command{Name: "save", Run: runSave, Help: helpSave,
	Description: "Save standard input to file(s) with format inferred by extension"}

type SaveContext struct {
	OverwriteExisting bool
	WriteIfEmpty      bool
	CreateDirectory   bool
	Quiet             bool
}

func helpSave() string {
	return `Unless options are set explicitly, existing files will not be overwritten and
logfiles withoutout any records will not be saved (useful if validate failed).

File name may be a template with {FIELD} placeholders replaced by field values.
For example, '{QSO_DATE}_{BAND}.adi' will create a separate file for each
contact date + band combination.  Quote the name to avoid shell expansion.
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
	for _, u := range ctx.UserdefFields {
		l.AddUserdef(u)
	}

	saveLog := func(l *adif.Logfile, file string) error {
		if !cctx.OverwriteExisting && fs.Exists(file) {
			return fmt.Errorf("output file %s already exists", file)
		}
		if len(l.Records) == 0 {
			if !cctx.WriteIfEmpty {
				return fmt.Errorf("no records in input, not saving to %s", file)
			}
			if !cctx.Quiet {
				fmt.Fprintf(os.Stderr, "Warning: saving %s with no records", file)
			}
		}
		if cctx.CreateDirectory {
			dir := path.Dir(file)
			if err := fs.MkdirAll(dir); err != nil && !errors.Is(err, os.ErrExist) {
				return err
			}
		}
		out, err := fs.Create(file)
		if err != nil {
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

	if len(l.Records) == 0 {
		if st.static && cctx.WriteIfEmpty {
			return saveLog(l, fname)
		}
		return fmt.Errorf("no records in input, not saving to %s", fname)
	}
	logs := make(map[string]*adif.Logfile)
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
			logs[file] = adif.NewLogfile()
			for _, f := range l.Header.Fields() {
				logs[file].Header.Set(f)
			}
			for _, u := range l.Userdef {
				logs[file].AddUserdef(u)
			}
		}
		logs[file].AddRecord(r)
	}

	errs := make([]error, len(logs))
	files := maps.Keys(logs)
	sort.Strings(files)
	for i, f := range files {
		errs[i] = saveLog(logs[f], f)
	}
	return errors.Join(errs...)
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
			pieces: []func(*adif.Record) string{func(_ *adif.Record) string { return s }}, static: true,
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
				if !unicode.IsPrint(c) {
					return '_'
				}
				if unicode.IsLetter(c) {
					return unicode.ToUpper(c)
				}
				// replace marks that are awkward in filenames
				switch c {
				case '/', '\\', ':', ';', '*', '?', '"', '\'', '`',
					'<', '>', '[', ']', '{', '}', '(', ')':
					return '-'
				}
				return c
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

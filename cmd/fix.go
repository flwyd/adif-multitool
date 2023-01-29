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
	"regexp"
	"strings"
	"time"

	"github.com/flwyd/adif-multitool/adif"
	"github.com/flwyd/adif-multitool/adif/spec"
)

var Fix = Command{Name: "fix", Run: runFix,
	Description: "Correct field formats to match the ADIF specification"}

var allNumeric = regexp.MustCompile("^[0-9]+$")

func runFix(ctx *Context, args []string) error {
	// TODO add any needed flags
	out := adif.NewLogfile()
	for _, f := range filesOrStdin(args) {
		l, err := readFile(ctx, f)
		if err != nil {
			return err
		}
		updateFieldOrder(out, l.FieldOrder)
		for _, rec := range l.Records {
			out.Records = append(out.Records, fixRecord(rec))
		}
	}
	return write(ctx, out)
}

func fixRecord(r *adif.Record) *adif.Record {
	fields := r.Fields()
	for i, f := range fields {
		if fieldType(f) == spec.DateDataType {
			f.Value = fixDate(f.Value)
		} else if fieldType(f) == spec.TimeDataType {
			f.Value = fixTime(f.Value)
		}
		fields[i] = f
	}
	return adif.NewRecord(fields...)
}

func fieldType(f adif.Field) spec.DataType {
	if fs, ok := spec.Fields[strings.ToUpper(f.Name)]; ok {
		return fs.Type
	}
	if f.Type.Identifier() != "" {
		for _, dt := range spec.DataTypes {
			if dt.Indicator == f.Type.Identifier() {
				return dt
			}
		}
	}
	return spec.StringDataType // reasonable default
}

var dateFormats = []string{"2006-1-2", "2006/1/2", "2006.1.2"}

func fixDate(d string) string {
	d = strings.TrimSpace(d)
	// TODO take date format as flag
	if allNumeric.MatchString(d) {
		return d
	}
	for _, pat := range dateFormats {
		if p, err := time.Parse(pat, d); err == nil {
			return p.Format("20060102")
		}
	}
	return d
}

var (
	timeWithSecs    = []string{"15:04:05", "3:04:05 PM", "3:04:05 pm", "3:04:05PM", "3:04:05pm"}
	timeWithoutSecs = []string{"15:04", "3:04 PM", "3:04 pm", "3:04PM", "3:04pm"}
)

func fixTime(t string) string {
	t = strings.TrimSpace(t)
	// TODO take date format as flag
	if allNumeric.MatchString(t) {
		switch len(t) {
		case 6, 4:
			return t
		case 5, 3:
			return "0" + t
		}
	}
	for _, pat := range timeWithSecs {
		if p, err := time.Parse(pat, t); err == nil {
			return p.Format("150405")
		}
	}
	for _, pat := range timeWithoutSecs {
		if p, err := time.Parse(pat, t); err == nil {
			return p.Format("1504")
		}
	}
	return t
}

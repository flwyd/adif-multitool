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
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/flwyd/adif-multitool/adif"
)

var Fix = Command{Name: "fix", Run: runFix,
	Description: "Correct field formats to match the ADIF specification"}

var allNumeric = regexp.MustCompile("^[0-9]+$")

func runFix(ctx *Context, args []string) error {
	// TODO add any needed flags
	srcs := argSources(args...)
	out := adif.NewLogfile("")
	for _, src := range srcs {
		l, err := readSource(ctx, src)
		if err != nil {
			return fmt.Errorf("error reading %s: %v", src, err)
		}
		for _, rec := range l.Records {
			out.Records = append(out.Records, fixRecord(rec))
		}
	}
	return write(ctx, out)
}

var (
	// TODO type registry from ADIF spec so these aren't needed
	dateFields = map[string]bool{
		"CLUBLOG_QSO_UPLOAD_DATE":  true,
		"EQSL_QSLRDATE":            true,
		"EQSL_QSLSDATE":            true,
		"HAMLOGEU_QSO_UPLOAD_DATE": true,
		"HAMQTH_QSO_UPLOAD_DATE":   true,
		"HRDLOG_QSO_UPLOAD_DATE":   true,
		"LOTW_QSLRDATE":            true,
		"LOTW_QSLSDATE":            true,
		"QRZCOM_QSO_UPLOAD_DATE":   true,
		"QSLRDATE":                 true,
		"QSLSDATE":                 true,
		"QSO_DATE":                 true,
		"QSO_DATE_OFF":             true,
	}
	timeFields = map[string]bool{
		"TIME_OFF": true,
		"TIME_ON":  true,
	}
)

func fixRecord(r *adif.Record) *adif.Record {
	fields := r.Fields()
	for i, f := range fields {
		if f.Type == adif.Date || dateFields[strings.ToUpper(f.Name)] {
			f.Value = fixDate(f.Value)
		} else if f.Type == adif.Time || timeFields[strings.ToUpper(f.Name)] {
			f.Value = fixTime(f.Value)
		}
		fields[i] = f
	}
	return adif.NewRecord(fields...)
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

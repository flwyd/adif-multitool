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
	"math"
	"regexp"
	"strconv"
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
	acc := accumulator{Out: out, Ctx: ctx}
	for _, f := range filesOrStdin(args) {
		l, err := acc.read(f)
		if err != nil {
			return err
		}
		updateFieldOrder(out, l.FieldOrder)
		for _, rec := range l.Records {
			out.Records = append(out.Records, fixRecord(rec, l))
		}
	}
	if err := acc.prepare(); err != nil {
		return err
	}
	// fix again in case userdef fields were added
	for _, r := range out.Records {
		for _, f := range r.Fields() {
			ff := fixField(f, out)
			if f != ff {
				r.Set(ff)
			}
		}
	}
	return write(ctx, out)
}

func fixRecord(r *adif.Record, l *adif.Logfile) *adif.Record {
	fields := r.Fields()
	for i, f := range fields {
		fields[i] = fixField(f, l)
	}
	return adif.NewRecord(fields...)
}

func fixField(f adif.Field, l *adif.Logfile) adif.Field {
	t := fieldType(f, l)
	if t == spec.DateDataType {
		f.Value = fixDate(f.Value)
	} else if t == spec.TimeDataType {
		f.Value = fixTime(f.Value)
	} else if t == spec.LocationDataType {
		f.Value = fixLocation(f.Value, f.Name)
	}
	return f
}

func fieldType(f adif.Field, l *adif.Logfile) spec.DataType {
	if fs, ok := spec.Fields[strings.ToUpper(f.Name)]; ok {
		return fs.Type
	}
	t := f.Type
	if u, ok := l.GetUserdef(f.Name); ok {
		t = u.Type
	}
	if t.Indicator() != "" {
		for _, dt := range spec.DataTypes {
			if dt.Indicator == t.Indicator() {
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

var gpsPattern = regexp.MustCompile(`^[-+]?\d{1,3}\.\d+$`)

func fixLocation(l, name string) string {
	l = strings.TrimSpace(l)
	if l == "" {
		return l
	}
	if gpsPattern.MatchString(l) {
		f, err := strconv.ParseFloat(l, 64)
		if err != nil {
			return l
		}
		var dir rune
		name = strings.ToUpper(name)
		if strings.Contains(name, "LATITUDE") {
			name = "LAT"
		} else if strings.Contains(name, "LONGITUDE") {
			name = "LON"
		}
		if strings.Contains(name, "LAT") {
			if math.Abs(f) > 90.0 {
				return l
			}
			if f >= 0.0 {
				dir = 'N'
			} else {
				dir = 'S'
			}
		} else if strings.Contains(name, "LON") {
			if math.Abs(f) > 180.0 {
				return l
			}
			if f >= 0.0 {
				dir = 'E'
			} else {
				dir = 'W'
			}
		} else {
			return l // can't tell if it's latitude or longitude, so can't set dir
		}
		f = math.Abs(f)
		deg := int(f)
		min := (f - float64(deg)) * 60.0
		return fmt.Sprintf("%c%03d %06.3f", dir, deg, min)
	}
	return l
}

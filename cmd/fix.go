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

var Fix = Command{Name: "fix", Run: runFix, Help: helpFix,
	Description: "Correct field formats to match the ADIF specification"}

var allNumeric = regexp.MustCompile("^[0-9]+$")

func helpFix() string {
	return `Fixable data formats:
  Date fields: 2006-01-02, 2006/01/02, 2006.01.02 (or without zero padding)
  Time fields (seconds): 15:04:05, 3:04:05 PM, 3:04:05pm
  Time fields (no seconds): 15:04, 3:04 PM, 3:04pm
  Location fields: decimal degrees (GPS coordinates)
  Country fields: ISO 3166-1 alpha-2 and alpha-3 codes
`
}

func runFix(ctx *Context, args []string) error {
	acc, err := newAccumulator(ctx)
	if err != nil {
		return err
	}
	for _, f := range filesOrStdin(args) {
		l, err := acc.read(f)
		if err != nil {
			return err
		}
		updateFieldOrder(acc.Out, l.FieldOrder)
		for _, rec := range l.Records {
			acc.Out.AddRecord(fixRecord(rec, l))
		}
	}
	if err := acc.prepare(); err != nil {
		return err
	}
	// fix again in case userdef fields were added
	for _, r := range acc.Out.Records {
		for _, f := range r.Fields() {
			ff := fixField(f, r, acc.Out)
			if f != ff {
				r.Set(ff)
			}
		}
	}
	return write(ctx, acc.Out)
}

func fixRecord(r *adif.Record, l *adif.Logfile) *adif.Record {
	fields := r.Fields()
	for i, f := range fields {
		fields[i] = fixField(f, r, l)
	}
	return adif.NewRecord(fields...)
}

func fixField(f adif.Field, r *adif.Record, l *adif.Logfile) adif.Field {
	f.Value = strings.TrimSpace(f.Value)
	t := fieldType(f, l)
	if t == spec.DateDataType {
		f.Value = fixDate(f.Value)
	} else if t == spec.TimeDataType {
		f.Value = fixTime(f.Value)
	} else if t == spec.LocationDataType {
		f.Value = fixLocation(f.Value, f.Name)
	} else if f.Name == spec.CountryField.Name || f.Name == spec.MyCountryField.Name {
		var state string
		if f.Name == spec.CountryField.Name {
			if s, ok := r.Get(spec.StateField.Name); ok {
				state = s.Value
			}
		} else if f.Name == spec.MyCountryField.Name {
			if s, ok := r.Get(spec.MyStateField.Name); ok {
				state = s.Value
			}
		}
		f.Value = fixCountry(f.Value, state)
	}
	return f
}

func fieldType(f adif.Field, l *adif.Logfile) spec.DataType {
	if fs, ok := spec.FieldNamed(f.Name); ok {
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

func fixCountry(c, state string) string {
	if e := spec.CountryEnumeration.Value(c); len(e) > 0 {
		return c
	}
	if cc, ok := spec.ISO3166Alpha[strings.ToUpper(c)]; ok {
		if len(cc.DXCC) == 1 {
			return cc.DXCC[0].EntityName
		}
		if state != "" {
			if e, ok := cc.Subdivisions[state]; ok {
				return e.EntityName
			}
			if len(cc.Subdivisions) > 0 {
				// if state isn't associated with a sub-national entity, use main one if the country has subdivisions defined
				return cc.DXCC[0].EntityName
			}
		}
	}
	return c
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
		name = strings.ToUpper(name)
		if strings.Contains(name, "LATITUDE") {
			name = "LAT"
		} else if strings.Contains(name, "LONGITUDE") {
			name = "LON"
		}
		if strings.Contains(name, "LAT") {
			s, err := formatLatitude(f)
			if err != nil {
				return l
			}
			return s
		} else if strings.Contains(name, "LON") {
			s, err := formatLongitude(f)
			if err != nil {
				return l
			}
			return s
		} else {
			return l // can't tell if it's latitude or longitude, so can't set dir
		}
	}
	return l
}

func formatLatitude(deg float64) (string, error) {
	if math.Abs(deg) > 90 {
		return "", fmt.Errorf("out of latitude range: %f", deg)
	}
	var dir rune
	if deg >= 0 {
		dir = 'N'
	} else {
		dir = 'S'
	}
	return formatLocation(deg, dir), nil
}

func formatLongitude(deg float64) (string, error) {
	if math.Abs(deg) > 180 {
		return "", fmt.Errorf("out of longitude range: %f", deg)
	}
	var dir rune
	if deg >= 0 {
		dir = 'E'
	} else {
		dir = 'W'
	}
	return formatLocation(deg, dir), nil
}

func formatLocation(degrees float64, dir rune) string {
	f := math.Abs(degrees)
	deg := int(f)
	min := (f - float64(deg)) * 60.0
	return fmt.Sprintf("%c%03d %06.3f", dir, deg, min)
}

// Copyright 2022 Google LLC
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

package adif

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var errNotSet = errors.New("not set")
var errEmpty = errors.New("empty value")

type Record struct {
	fields  []Field
	named   map[string]int // map of field name to index in fields slice
	comment string
}

func NewRecord(fs ...Field) *Record {
	r := &Record{fields: make([]Field, 0, len(fs)), named: make(map[string]int)}
	for _, f := range fs {
		r.Set(f)
	}
	return r
}

func (r *Record) Fields() []Field {
	f := make([]Field, len(r.fields))
	copy(f, r.fields)
	return f
}

func (r *Record) Get(name string) (f Field, ok bool) {
	name = strings.ToUpper(name)
	i, ok := r.named[name]
	if !ok {
		return Field{}, false
	}
	return r.fields[i], true
}

func (r *Record) ParseBool(name string) (bool, error) {
	f, ok := r.Get(name)
	if !ok {
		return false, errNotSet
	}
	switch f.Value {
	case "Y", "y":
		return true, nil
	case "N", "n":
		return false, nil
	case "":
		return false, errEmpty
	default:
		return false, fmt.Errorf("invalid boolean value %q", f.Value)
	}
}

var adifNumberPat = regexp.MustCompile(`^-?(\d+|\d+\.\d*|\.\d+)$`)

func (r *Record) ParseFloat(name string) (float64, error) {
	f, ok := r.Get(name)
	if !ok {
		return 0.0, errNotSet
	}
	if f.Value == "" {
		return 0.0, errEmpty
	}
	if !adifNumberPat.MatchString(f.Value) { // ParseFloat is more broad than ADIF
		return 0.0, fmt.Errorf("invalid number format %q", f.Value)
	}
	v, err := strconv.ParseFloat(f.Value, 64)
	if err != nil {
		return 0.0, err
	}
	return v, nil
}

func (r *Record) ParseInt(name string) (int, error) {
	f, ok := r.Get(name)
	if !ok {
		return 0, errNotSet
	}
	if f.Value == "" {
		return 0, errEmpty
	}
	if !adifNumberPat.MatchString(f.Value) { // ParseFloat is more broad than ADIF
		return 0, fmt.Errorf("invalid number format %q", f.Value)
	}
	v, err := strconv.Atoi(f.Value)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (r *Record) ParseDate(name string) (time.Time, error) {
	f, ok := r.Get(name)
	if !ok {
		return time.Time{}, errNotSet
	}
	if f.Value == "" {
		return time.Time{}, errEmpty
	}
	return time.ParseInLocation("20060102", f.Value, time.UTC)
}

func (r *Record) ParseTime(name string) (time.Time, error) {
	f, ok := r.Get(name)
	if !ok {
		return time.Time{}, errNotSet
	}
	if f.Value == "" {
		return time.Time{}, errEmpty
	}
	switch len(f.Value) {
	case 4:
		return time.ParseInLocation("1504", f.Value, time.UTC)
	case 6:
		return time.ParseInLocation("150405", f.Value, time.UTC)
	default:
		return time.Time{}, fmt.Errorf("invalid time format %q", f.Value)
	}
}

func (r *Record) Set(f Field) error {
	f.Name = strings.ToUpper(f.Name)
	if len(f.Name) == 0 {
		return fmt.Errorf("empty field name in %s", f)
	}
	if i, ok := r.named[f.Name]; ok {
		r.fields[i] = f
	} else {
		r.named[f.Name] = len(r.fields)
		r.fields = append(r.fields, f)
	}
	return nil
}

func (r *Record) GetComment() string {
	return r.comment
}

func (r *Record) SetComment(c string) {
	r.comment = c
}

func (r *Record) String() string {
	return fmt.Sprint(r.fields)
}

// Equal compares two records for equality of fields, ignoring order and comments.
// Records are considered equal even if one has assigned empty fields while the
// other does not have a field of that name set.
func (r *Record) Equal(o *Record) bool {
	for _, f := range r.fields {
		ff, ok := o.Get(f.Name)
		if ok {
			if ff != f { // TODO add Field.Equal with handling for unknown vs. known type
				return false
			}
		} else if f.Value != "" { // treat empty and missing as equal
			return false
		}
	}
	// check both directions in case o has fields r lacks
	for _, ff := range o.fields {
		f, ok := r.Get(ff.Name)
		if ok {
			if f != ff {
				return false
			}
		} else if ff.Value != "" {
			return false
		}
	}
	// ignore comment
	return true
}

func (r *Record) Empty() bool {
	return len(r.fields) == 0
}

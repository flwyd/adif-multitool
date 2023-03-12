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
	"fmt"
	"strings"
)

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

// Equal compares two records for equality of fields, ignoring order.
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

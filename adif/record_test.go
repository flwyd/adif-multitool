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

package adif

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNew(t *testing.T) {
	tests := []struct{ fields, want []Field }{
		{
			fields: []Field{},
			want:   []Field{},
		},
		{
			fields: []Field{{Name: "FOO", Value: "123", Type: Number}},
			want:   []Field{{Name: "FOO", Value: "123", Type: Number}},
		},
		{
			fields: []Field{
				{Name: "freq", Value: "14.250", Type: Number},
				{Name: "QsO_dAtE", Value: "20151124", Type: Date},
				{Name: "not an ADIF field", Value: "'sup?"}},
			want: []Field{
				{Name: "FREQ", Value: "14.250", Type: Number},
				{Name: "QSO_DATE", Value: "20151124", Type: Date},
				{Name: "NOT AN ADIF FIELD", Value: "'sup?"}},
		},
	}
	for _, tc := range tests {
		r := NewRecord(tc.fields...)
		if diff := cmp.Diff(tc.want, r.Fields()); diff != "" {
			t.Errorf("NewRecord(%v) got unexpected result, diff:\n%s", tc.fields, diff)
		}
		for _, f := range tc.want {
			if got, ok := r.Get(f.Name); !ok || got != f {
				t.Errorf("Get(%s) got %v, want %v", f.Name, got, f)
			}
		}
	}
}

func TestSet(t *testing.T) {
	r := NewRecord()
	if got, ok := r.Get("FOO"); (ok || got != Field{}) {
		t.Errorf(`Get("FOO") on empty record got %v`, got)
	}
	want := Field{Name: "FOO", Value: "Bar"}
	r.Set(want)
	if got, ok := r.Get("FOO"); !ok || got != want {
		t.Errorf("Get(%q) got %v, want %v", "FOO", got, want)
	}
	want = Field{Name: "A_FIELD", Value: "123", Type: Number}
	r.Set(want)
	if got, ok := r.Get("A_FIELD"); !ok || got != want {
		t.Errorf("Get(%q) got %v, want %v", "A_FIELD", got, want)
	}
	if got, ok := r.Get("a_FielD"); !ok || got != want {
		t.Errorf("Get(%q) got %v, want %v", "a_FielD", got, want)
	}
	r.Set(Field{Name: "foo", Value: "1234", Type: Time})
	want = Field{Name: "FOO", Value: "1234", Type: Time}
	if got, ok := r.Get("FOO"); !ok || got != want {
		t.Errorf("Get(%q) got %v, want %v", "FOO", got, want)
	}
	if got, ok := r.Get(""); (ok || got != Field{}) {
		t.Errorf(`Get("") got %v, want nothing`, got)
	}
	if got, ok := r.Get("BAR"); (ok || got != Field{}) {
		t.Errorf(`Get("BAR") got %v, want nothing`, got)
	}
}

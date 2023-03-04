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

// Package adif defines basic types (Field, Record, Logfile) for working with
// amateur radio logs in ADIF format.  It also defines I/O types which can read
// and write ADIF data in several formats, including the ADIF-specified ADI and
// ADX formats as well as CSV and JSON.  The adif package does not have any
// code dealing with specific fields or enumerations defined by specification;
// see the adif/spec package for such code.
package adif

import (
	"fmt"
	"math"
	"strings"
)

type Logfile struct {
	Records    []*Record
	Header     *Record
	Userdef    []UserdefField
	Comment    string
	Filename   string
	FieldOrder []string
}

func NewLogfile() *Logfile {
	return &Logfile{
		Records: make([]*Record, 0),
		Header:  NewRecord(),
	}
}

func (f *Logfile) String() string {
	if f.Filename == "" {
		return "(standard input)"
	}
	return f.Filename
}

func (f *Logfile) AddRecord(r *Record) *Logfile {
	f.Records = append(f.Records, r)
	return f
}

func (f *Logfile) AddUserdef(u UserdefField) error {
	for _, x := range f.Userdef {
		if strings.EqualFold(u.Name, x.Name) {
			if u.Type != x.Type && (u.Type != TypeUnspecified || x.Type != TypeUnspecified) {
				return fmt.Errorf("userdef field %s type conflicts with existing %s", u, x)
			}
			if len(u.EnumValues) > 0 || len(x.EnumValues) > 0 {
				seenOld := make(map[string]bool)
				seenNew := make(map[string]bool)
				for _, v := range u.EnumValues {
					seenNew[strings.ToUpper(v)] = true
				}
				for _, v := range x.EnumValues {
					seenOld[strings.ToUpper(v)] = true
				}
				var missOld, missNew bool
				for v := range seenOld {
					if !seenNew[v] {
						missOld = true
					}
				}
				for v := range seenNew {
					if !seenOld[v] {
						missNew = true
					}
				}
				if missOld && missNew { // neither is a subset of the other
					return fmt.Errorf("userdef field %s enumeration conflicts with existing %s", u, x)
				}
				for _, v := range u.EnumValues {
					if !seenOld[strings.ToUpper(v)] {
						x.EnumValues = append(x.EnumValues, v)
					}
				}
			}
			x.Min = math.Min(x.Min, u.Min)
			x.Max = math.Max(x.Max, u.Max)
			if x.Type == TypeUnspecified {
				x.Type = u.Type
			}
			return nil
		}
	}
	// doesn't match an existing field
	f.Userdef = append(f.Userdef, u)
	return nil
}

func (f *Logfile) GetUserdef(name string) (UserdefField, bool) {
	for _, u := range f.Userdef {
		if strings.EqualFold(name, u.Name) {
			return u, true
		}
	}
	return UserdefField{}, false
}

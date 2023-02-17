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
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type jsonRecord map[string]any

func newJsonRecord(r *Record, typed bool) jsonRecord {
	j := make(jsonRecord)
	for _, f := range r.Fields() {
		if typed {
			switch f.Type {
			default:
				// TODO look up field type in spec
				j[f.Name] = f.Value
			case TypeBoolean:
				if f.Value == "Y" || f.Value == "y" {
					j[f.Name] = true
				} else if f.Value == "N" || f.Value == "n" {
					j[f.Name] = false
				} else {
					j[f.Name] = f.Value
				}
			case TypeNumber:
				if v, err := strconv.ParseFloat(f.Value, 64); err != nil {
					j[f.Name] = f.Value
				} else {
					j[f.Name] = v
				}
			}
		} else {
			j[f.Name] = f.Value
		}
	}
	return j
}

func (j jsonRecord) toRecord() (*Record, error) {
	r := NewRecord()
	for k, v := range j {
		switch vv := v.(type) {
		default:
			// TODO handle USERDEF fields
			return nil, fmt.Errorf("unsupported JSON field type %q: %v", k, v)
		case string:
			r.Set(Field{Name: k, Value: vv})
		case bool:
			if vv {
				r.Set(Field{Name: k, Value: "Y", Type: TypeBoolean})
			} else {
				r.Set(Field{Name: k, Value: "N", Type: TypeBoolean})
			}
		case json.Number:
			r.Set(Field{Name: k, Value: vv.String(), Type: TypeNumber})
		case nil:
			r.Set(Field{Name: k, Value: ""})
		}
	}
	return r, nil
}

type jsonFile struct {
	Header  jsonRecord   `json:"HEADER"`
	Records []jsonRecord `json:"RECORDS"`
}

type JSONIO struct {
	HTMLSafe    bool
	Indent      int
	TypedOutput bool
}

func NewJSONIO() *JSONIO { return &JSONIO{} }

func (_ *JSONIO) String() string { return "json" }

func (o *JSONIO) Read(in io.Reader) (*Logfile, error) {
	d := json.NewDecoder(in)
	d.UseNumber()
	var f jsonFile
	if err := d.Decode(&f); err != nil {
		return nil, fmt.Errorf("JSON decoding error: %w", err)
	}
	l := NewLogfile()
	if f.Header != nil {
		h, err := f.Header.toRecord()
		if err != nil {
			return nil, err
		}
		l.Header = h
	}
	l.Records = make([]*Record, len(f.Records))
	for i, r := range f.Records {
		rec, err := r.toRecord()
		if err != nil {
			return nil, err
		}
		l.Records[i] = rec
	}
	return l, nil
}

func (o *JSONIO) Write(l *Logfile, out io.Writer) error {
	var f jsonFile
	f.Header = newJsonRecord(l.Header, o.TypedOutput)
	f.Records = make([]jsonRecord, len(l.Records))
	for i, r := range l.Records {
		f.Records[i] = newJsonRecord(r, o.TypedOutput)
	}
	e := json.NewEncoder(out)
	e.SetIndent("", strings.Repeat(" ", o.Indent))
	e.SetEscapeHTML(o.HTMLSafe)
	if err := e.Encode(f); err != nil {
		return fmt.Errorf("JSON encoding error: %w", err)
	}
	return nil
}

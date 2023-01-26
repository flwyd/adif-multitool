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
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

type CSVIO struct {
	Comma            rune
	Comment          rune
	LazyQuotes       bool
	TrimLeadingSpace bool
	UseCRLF          bool
}

func NewCSVIO() *CSVIO {
	return &CSVIO{Comma: ','}
}

func (o *CSVIO) String() string { return "csv" }

func (o *CSVIO) Read(in NamedReader) (*Logfile, error) {
	l := NewLogfile(in.Name())
	c := csv.NewReader(in)
	c.Comma = o.Comma
	c.Comment = o.Comment
	c.LazyQuotes = o.LazyQuotes
	c.TrimLeadingSpace = o.TrimLeadingSpace
	h, err := c.Read()
	if err == io.EOF {
		return nil, fmt.Errorf("got EOF reading CSV header row from %s", in.Name())
	}
	if err != nil {
		return nil, fmt.Errorf("error reading CSV header row from %s: %v", in.Name(), err)
	}
	l.FieldOrder = make([]string, len(h))
	for i, n := range h {
		l.FieldOrder[i] = strings.ToUpper(n)
	}
	// TODO if there are any USERDEF fields, add them to l.Header
	for line, err := c.Read(); err != io.EOF; line, err = c.Read() {
		lnum, _ := c.FieldPos(0)
		if err != nil {
			return nil, fmt.Errorf("error reading CSV %s line %d: %v", in.Name(), lnum, err)
		}
		r := NewRecord()
		for i, v := range line {
			if i >= len(h) {
				// CONSIDER ignoring this column and logging a warning
				return nil, fmt.Errorf("extra field value %s at line %d column %d of %s", v, lnum, i, in.Name())
			}
			if err := r.Set(Field{Name: h[i], Value: v}); err != nil {
				return nil, fmt.Errorf("could not set field %s to %s: %v", h[i], v, err)
			}
		}
		l.Records = append(l.Records, r)
	}
	return l, nil
}

func (o *CSVIO) Write(l *Logfile, out io.Writer) error {
	order := make([]string, len(l.FieldOrder))
	seen := make(map[string]bool)
	for i, n := range l.FieldOrder {
		n = strings.ToUpper(n)
		order[i] = n
		seen[n] = true
	}
	for _, r := range l.Records {
		for _, f := range r.Fields() {
			n := strings.ToUpper(f.Name)
			if !seen[n] {
				order = append(order, n)
				seen[n] = true
			}
		}
	}
	if len(order) == 0 {
		return nil
	}
	c := csv.NewWriter(out)
	defer c.Flush()
	c.Comma = o.Comma
	c.UseCRLF = o.UseCRLF
	// CSV header row
	if err := c.Write(order); err != nil {
		return fmt.Errorf("error writing CSV header to %s: %v", l, err)
	}
	for i, r := range l.Records {
		row := make([]string, 0, len(r.fields))
		for _, n := range order {
			if f, ok := r.Get(n); !ok {
				row = append(row, "")
			} else {
				row = append(row, f.Value)
			}
		}
		if err := c.Write(row); err != nil {
			return fmt.Errorf("error writing CSV record %d to %s: %v", i, l, err)
		}
	}
	return nil
}

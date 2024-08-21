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
	"errors"
	"fmt"
	"io"
	"strings"
)

type CSVIO struct {
	Comma             rune
	Comment           rune
	CRLF              bool
	LazyQuotes        bool
	RequireFullRecord bool
	TrimLeadingSpace  bool
	OmitHeader        bool
}

func NewCSVIO() *CSVIO {
	return &CSVIO{Comma: ','}
}

func (o *CSVIO) String() string { return "csv" }

func (o *CSVIO) Read(in io.Reader) (*Logfile, error) {
	l := NewLogfile()
	c := csv.NewReader(in)
	c.ReuseRecord = true
	c.Comma = o.Comma
	c.Comment = o.Comment
	c.LazyQuotes = o.LazyQuotes
	c.TrimLeadingSpace = o.TrimLeadingSpace
	if o.RequireFullRecord {
		c.FieldsPerRecord = 0
	} else {
		c.FieldsPerRecord = -1
	}
	header, err := c.Read()
	if errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("got EOF reading CSV header row")
	}
	if err != nil {
		return nil, fmt.Errorf("error reading CSV header row: %w", err)
	}
	// copy header, since array will be reused
	h := make([]string, len(header))
	copy(h, header)
	l.FieldOrder = make([]string, len(h))
	for i, n := range h {
		l.FieldOrder[i] = strings.ToUpper(n)
	}
	// TODO if there are any USERDEF fields, add them to l.Header
	for line, err := c.Read(); err != io.EOF; line, err = c.Read() {
		lnum, _ := c.FieldPos(0)
		if err != nil {
			return nil, fmt.Errorf("error reading CSV: %w", err)
		}
		r := NewRecord()
		for i, v := range line {
			// if RequreFullRecord, c.Read already returned an error for under/overflow, but catch overflow even if it's false
			if i >= len(h) {
				return nil, fmt.Errorf("extra field value %q at line %d field %d", v, lnum, i)
			}
			if err := r.Set(Field{Name: h[i], Value: v}); err != nil {
				return nil, fmt.Errorf("could not set field %s to %q: %w", h[i], v, err)
			}
		}
		for i := len(line); i < len(h); i++ {
			if err := r.Set(Field{Name: h[i], Value: ""}); err != nil {
				return nil, fmt.Errorf("could not set field %s to empty: %w", h[i], err)
			}
		}
		l.AddRecord(r)
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
	c.Comma = o.Comma
	c.UseCRLF = o.CRLF
	// CSV header row
	if !o.OmitHeader {
		if err := c.Write(order); err != nil {
			return fmt.Errorf("writing CSV header to %s: %w", l, err)
		}
	}
	row := make([]string, len(order))
	for i, r := range l.Records {
		for i, n := range order {
			if f, ok := r.Get(n); !ok {
				row[i] = ""
			} else {
				row[i] = f.Value
			}
		}
		if err := c.Write(row); err != nil {
			return fmt.Errorf("writing CSV record %d to %s: %w", i+1, l, err)
		}
	}
	c.Flush()
	return c.Error()
}

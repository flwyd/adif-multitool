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
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

type TSVIO struct {
	CRLF               bool
	EscapeSpecial      bool
	IgnoreEmptyHeaders bool
	OmitHeader         bool
}

func NewTSVIO() *TSVIO { return &TSVIO{} }

func (_ *TSVIO) String() string { return "tsv" }

func (o *TSVIO) Read(r io.Reader) (*Logfile, error) {
	scan := bufio.NewScanner(r)
	if !scan.Scan() {
		return nil, errors.New("no TSV header row")
	}
	headRow := scan.Text()
	if strings.TrimSpace(headRow) == "" {
		return nil, errors.New("empty TSV header row")
	}
	head := strings.Split(headRow, "\t")
	seen := make(map[string]bool)
	for i, h := range head {
		head[i] = strings.ToUpper(o.unescape(strings.TrimSpace(h)))
		if head[i] == "" {
			if !o.IgnoreEmptyHeaders {
				return nil, fmt.Errorf("empty TSV header field position %d", i)
			}
		} else {
			if seen[head[i]] {
				return nil, fmt.Errorf("duplicate TSV header %q position %d", head[i], i)
			}
			seen[head[i]] = true
		}
	}
	line := 1 // after reading header
	l := NewLogfile()
	l.FieldOrder = head
	for scan.Scan() {
		line++
		fs := strings.Split(scan.Text(), "\t")
		if len(fs) > len(head) {
			return nil, fmt.Errorf("line %d has %d fields, more than %d in TSV header", len(l.Records)+2, len(fs), len(head))
		}
		if len(fs) == 1 && fs[0] == "" {
			continue // skip blank lines
		}
		fields := make([]Field, len(fs))
		for i, f := range fs {
			fields[i] = Field{Name: head[i], Value: o.unescape(f)}
		}
		l.AddRecord(NewRecord(fields...))
	}
	if err := scan.Err(); err != nil {
		return nil, fmt.Errorf("reading TSV line %d: %w", line, err)
	}
	return l, nil
}

func (o *TSVIO) Write(l *Logfile, w io.Writer) error {
	order := make([]string, len(l.FieldOrder))
	seen := make(map[string]bool)
	for i, n := range l.FieldOrder {
		n = strings.ToUpper(n)
		if !o.EscapeSpecial && strings.ContainsAny(n, "\t\r\n") {
			return fmt.Errorf("invalid TSV field name %q", n)
		}
		order[i] = n
		seen[n] = true
	}
	for _, r := range l.Records {
		for _, f := range r.Fields() {
			n := strings.ToUpper(f.Name)
			if !o.EscapeSpecial && strings.ContainsAny(n, "\t\r\n") {
				return fmt.Errorf("invalid TSV field name %q", n)
			}
			if !o.EscapeSpecial && strings.ContainsAny(f.Value, "\t\r\n") {
				return fmt.Errorf("invalid TSV field value %q = %q", n, f.Value)
			}
			if !seen[n] {
				order = append(order, n)
				seen[n] = true
			}
		}
	}
	if len(order) == 0 {
		return nil
	}
	out := bufio.NewWriter(w)
	writeDelim := func(i int) error {
		if i < len(order)-1 {
			if err := out.WriteByte('\t'); err != nil {
				return err
			}
		} else {
			if o.CRLF {
				if err := out.WriteByte('\r'); err != nil {
					return err
				}
			}
			if err := out.WriteByte('\n'); err != nil {
				return err
			}
		}
		return nil
	}
	if !o.OmitHeader {
		for i, h := range order {
			if _, err := out.WriteString(o.escape(h)); err != nil {
				return fmt.Errorf("writing TSV header: %w", err)
			}
			if err := writeDelim(i); err != nil {
				return fmt.Errorf("writing TSV header: %w", err)
			}
		}
	}
	for _, r := range l.Records {
		for i, h := range order {
			f, _ := r.Get(h)
			if _, err := out.WriteString(o.escape(f.Value)); err != nil {
				return fmt.Errorf("writing TSV record: %w", err)
			}
			if err := writeDelim(i); err != nil {
				return fmt.Errorf("writing TSV header: %w", err)
			}
		}
	}
	return out.Flush()
}

func (o *TSVIO) escape(s string) string {
	if !o.EscapeSpecial {
		return s
	}
	if !strings.ContainsAny(s, "\t\r\n\\") {
		return s
	}
	var out strings.Builder
	for _, r := range s {
		switch r {
		case '\t':
			out.WriteString(`\t`)
		case '\r':
			out.WriteString(`\r`)
		case '\n':
			out.WriteString(`\n`)
		case '\\':
			out.WriteString(`\\`)
		default:
			out.WriteRune(r)
		}
	}
	return out.String()
}

func (o *TSVIO) unescape(s string) string {
	if !o.EscapeSpecial {
		return s
	}
	if !strings.ContainsRune(s, '\\') {
		return s
	}
	var out strings.Builder
	r := bufio.NewReader(strings.NewReader(s))
	for {
		chunk, err := r.ReadString('\\')
		if err != nil {
			out.WriteString(chunk)
			break
		}
		out.WriteString(chunk[0 : len(chunk)-1]) // don't write backslash
		next, _, err := r.ReadRune()
		if err != nil {
			out.WriteRune('\\')
			break
		}
		switch next {
		case '\\':
			out.WriteRune('\\')
		case 't':
			out.WriteRune('\t')
		case 'r':
			out.WriteRune('\r')
		case 'n':
			out.WriteRune('\n')
		default:
			out.WriteRune('\\') // leave unknown escapes unmolested
			out.WriteRune(next)
		}
	}
	return out.String()
}

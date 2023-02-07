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
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type ADIIO struct {
	LowerCase bool // TODO consider a case enum: keep, upper, lower, or just get rid of this option
	FieldSep  Separator
	RecordSep Separator
}

func NewADIIO() *ADIIO {
	return &ADIIO{FieldSep: SeparatorSpace, RecordSep: SeparatorNewline}
}

func (_ *ADIIO) String() string { return "adi" }

func (o *ADIIO) Read(in io.Reader) (*Logfile, error) {
	var comments []string
	l := NewLogfile()
	r := bufio.NewReader(in)
	s, err := r.ReadString('<')
	if err == io.EOF {
		if s != "" {
			l.Comment = s
		}
		return l, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error reading to first tag: %w", err)
	}
	// final byte is '<'
	comments = append(comments, s[0:len(s)-1])
	// ADIF specification seems to imply that without a comment at the start
	// of a file, and thus the first character is '<', then there is no header
	// and the < starts the first record.  This invariant may not hold for all
	// software though, so allow an <EOH> even if we didn't get a comment.
	cur := NewRecord()
	var sawHeader, sawRecord bool
	for { // invariant: last byte read was '<'
		s, err = r.ReadString('>')
		if err == io.EOF {
			return nil, fmt.Errorf("unfinished ADI tag at end: %q", s)
		}
		if err != nil {
			return nil, fmt.Errorf("error reading ADI tag %q: %w", s, err)
		}
		if s == ">" {
			return nil, fmt.Errorf("invalid ADI tag <>")
		}
		tag := strings.Split(s[0:len(s)-1], ":")
		switch len(tag) {
		case 1:
			switch strings.ToUpper(tag[0]) {
			case "EOH":
				if sawHeader {
					return nil, fmt.Errorf("invalid ADI file with two <EOH> headers")
				}
				if sawRecord {
					return nil, fmt.Errorf("invalid ADI file with <EOH> header after first <EOR> record")
				}
				sawHeader = true
				cur.SetComment(strings.Join(comments, o.RecordSep.Val()))
				l.Header = cur
				cur = NewRecord()
				comments = nil
			case "EOR":
				sawRecord = true
				cur.SetComment(strings.Join(comments, o.RecordSep.Val()))
				l.Records = append(l.Records, cur)
				cur = NewRecord()
				comments = nil
			default:
				return nil, fmt.Errorf("invalid ADI field without length <%s", s)
			}
		case 2, 3:
			length, err := strconv.Atoi(tag[1])
			if err != nil || length < 0 {
				return nil, fmt.Errorf("invalid ADI field length <%s", s)
			}
			v := make([]byte, length)
			if _, err = io.ReadFull(r, v); err != nil {
				return nil, fmt.Errorf("error reading ADI field value <%s got %q: %w", s, v, err)
			}
			// spec says everything is ASCII, but this accepts UTF-8
			// as long as the tag length is accurate in bytes
			f := Field{Name: tag[0], Value: string(v)}
			if len(tag) == 3 {
				f.Type, err = DataTypeFromIdentifier(tag[2])
				if err != nil {
					return nil, fmt.Errorf("%v from <%s", err, s)
				}
			}
			cur.Set(f)
		default:
			return nil, fmt.Errorf("invalid ADI tag format <%s", s)
		}
		// arbitrary text between one field or record and the next
		c, err := r.ReadString('<')
		c = strings.TrimFunc(strings.TrimSuffix(c, "<"), unicode.IsSpace)
		if c != "" {
			comments = append(comments, c)
		}
		if err == io.EOF {
			if len(cur.fields) != 0 {
				return nil, fmt.Errorf("final record missing <EOR>: %s", cur)
			}
			if len(comments) > 0 {
				l.Comment = strings.Join(comments, o.RecordSep.Val())
			}
			return l, nil
		}
		if err != nil {
			return nil, fmt.Errorf("error reading ADI intra-field text: %w", err)
		}
	}
}

func (o *ADIIO) Write(l *Logfile, out io.Writer) error {
	b := bufio.NewWriter(out)
	defer b.Flush()
	if !l.Header.Empty() {
		c := l.Header.GetComment()
		if c == "" {
			c = "ADI format, see https://adif.org.uk/"
		}
		if err := o.writeComment(b, c, o.RecordSep.Val()); err != nil {
			return fmt.Errorf("error writing ADI header: %w", err)
		}
		for _, f := range l.Header.Fields() {
			if err := o.writeField(f, b); err != nil {
				return fmt.Errorf("error writing ADI header: %w", err)
			}
		}
		if _, err := b.WriteString(fmt.Sprintf("<%s>%s", o.fixCase("EOH"), o.RecordSep.Val())); err != nil {
			return fmt.Errorf("error writing ADI header: %w", err)
		}
	}
	for i, r := range l.Records {
		if c := r.GetComment(); c != "" {
			if err := o.writeComment(b, c, o.FieldSep.Val()); err != nil {
				return fmt.Errorf("error writing ADI record comment: %w", err)
			}
		}
		seen := make(map[string]bool)
		for _, n := range l.FieldOrder {
			if f, ok := r.Get(n); ok {
				if err := o.writeField(f, b); err != nil {
					return fmt.Errorf("error writing ADI record #%d: %w", i, err)
				}
				seen[f.Name] = true
			}
		}
		for _, f := range r.Fields() {
			if !seen[f.Name] {
				if err := o.writeField(f, b); err != nil {
					return fmt.Errorf("error writing ADI record #%d: %w", i, err)
				}
			}
		}
		if _, err := b.WriteString(fmt.Sprintf("<%s>%s", o.fixCase("EOR"), o.RecordSep.Val())); err != nil {
			return fmt.Errorf("error writing ADI record #%d: %w", i, err)
		}
	}
	if l.Comment != "" {
		if err := o.writeComment(b, l.Comment, "\n"); err != nil {
			return fmt.Errorf("error writing ADI comment: %w", err)
		}
	}
	return nil
}

func (o *ADIIO) writeField(f Field, b *bufio.Writer) error {
	var tag string
	// TODO error if IntlString/non-ASCII, unless a flag allows
	if f.Type == Unspecified {
		tag = fmt.Sprintf("<%s:%d>", o.fixCase(f.Name), len(f.Value))
	} else {
		tag = fmt.Sprintf("<%s:%d:%s>", o.fixCase(f.Name), len(f.Value), o.fixCase(f.Type.Identifier()))
	}
	if _, err := b.WriteString(fmt.Sprintf("%s%s%s", tag, f.Value, o.FieldSep.Val())); err != nil {
		return fmt.Errorf("error writing %s: %w", f, err)
	}
	return nil
}

var escapeAngleBrackets = strings.NewReplacer("<", "&lt;", ">", "&gt;")

func (o *ADIIO) writeComment(w io.Writer, comment, suffix string) error {
	if _, err := escapeAngleBrackets.WriteString(w, comment); err != nil {
		return err
	}
	if _, err := escapeAngleBrackets.WriteString(w, suffix); err != nil {
		return err
	}
	return nil
}

func (o *ADIIO) fixCase(s string) string {
	if o.LowerCase {
		return strings.ToLower(s)
	}
	return strings.ToUpper(s)
}

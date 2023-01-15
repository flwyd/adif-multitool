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
	"time"
)

// TODO options like space/newline field separators, lf/crlf for records, case
type ADIIO struct {
	LowerCase bool // TODO consider a case enum: keep, upper, lower
	FieldSep  Separator
	RecordSep Separator
	// TODO add Comment string to Record, just print that
	HeaderCommentFn func(*Logfile) string
}

func NewADIIO() *ADIIO {
	return &ADIIO{FieldSep: SeparatorSpace, RecordSep: SeparatorNewline,
		HeaderCommentFn: func(l *Logfile) string {
			return fmt.Sprintf("Generated %s with %d records", time.Now().Format(time.RFC1123Z), len(l.Records))
		},
	}
}

func (_ *ADIIO) Read(in Source) (*Logfile, error) {
	l := NewLogfile(in.Name())
	r := bufio.NewReader(in)
	s, err := r.ReadString('<')
	if err == io.EOF {
		if s != "" {
			l.Comment = s
		}
		return l, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error reading to first tag of %s: %v", in.Name(), err)
	}
	if len(s) > 1 { // final byte is '<'
		// TODO ADIF specification says the comment is part of the header, maybe
		// make a Header type that contains zero-or-more comments.
		l.Comment = s[0 : len(s)-1]
	}
	// ADIF specification seems to imply that without a comment at the start
	// of a file, and thus the first character is '<', then there is no header
	// and the < starts the first record.  This invariant may not hold for all
	// software though, so allow an <EOH> even if we didn't get a comment.
	cur := NewRecord()
	var sawHeader, sawRecord bool
	for { // invariant: last byte read was '<'
		s, err = r.ReadString('>')
		if err == io.EOF {
			return nil, fmt.Errorf("unfinished ADI tag at end of file %s: %q", in.Name(), s)
		}
		if err != nil {
			return nil, fmt.Errorf("error reading ADI tag %q in %s: %v", s, in.Name(), err)
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
					return nil, fmt.Errorf("invalid ADI file with two <EOH> headers in %s", in.Name())
				}
				if sawRecord {
					return nil, fmt.Errorf("invalid ADI file with <EOH> header after first <EOR> record in %s", in.Name())
				}
				sawHeader = true
				l.Header = cur
				cur = NewRecord()
			case "EOR":
				sawRecord = true
				l.Records = append(l.Records, cur)
				cur = NewRecord()
			default:
				return nil, fmt.Errorf("invalid ADI field without length <%s in %s", s, in.Name())
			}
		case 2, 3:
			length, err := strconv.Atoi(tag[1])
			if err != nil || length < 0 {
				return nil, fmt.Errorf("invalid ADI field length <%s in %s", s, in.Name())
			}
			v := make([]byte, length)
			if _, err = io.ReadFull(r, v); err != nil {
				return nil, fmt.Errorf("error reading ADI field value <%s got %q in %s: %v", s, v, in.Name(), err)
			}
			// spec says everything is ASCII, but this accepts UTF-8
			// as long as the tag length is accurate in bytes
			f := Field{Name: tag[0], Value: string(v)}
			if len(tag) == 3 {
				f.Type, err = DataTypeFromIdentifier(tag[2])
				if err != nil {
					return nil, fmt.Errorf("%v from <%s in %s", err, s, in.Name())
				}
			}
			cur.Set(f)
		default:
			return nil, fmt.Errorf("invalid ADI tag format <%s in %s", s, in.Name())
		}
		// arbitrary text between one field or record and the next
		if _, err := r.ReadString('<'); err == io.EOF {
			return l, nil
		}
		if err != nil {
			return nil, fmt.Errorf("error reading ADI intra-field text in %s: %v", in.Name(), err)
		}
	}
}

func (o *ADIIO) Write(l *Logfile, out io.Writer) error {
	b := bufio.NewWriter(out)
	defer b.Flush()
	if _, err := b.WriteString(o.HeaderCommentFn(l) + o.RecordSep.Val()); err != nil {
		return fmt.Errorf("error writing ADI header: %v", err)
	}
	for _, f := range l.Header.Fields() {
		if err := o.writeField(f, b); err != nil {
			return fmt.Errorf("error writing ADI header: %v", err)
		}
	}
	if _, err := b.WriteString(fmt.Sprintf("<%s>%s", o.fixCase("EOH"), o.RecordSep.Val())); err != nil {
		return fmt.Errorf("error writing ADI header: %v", err)
	}
	for i, r := range l.Records {
		seen := make(map[string]bool)
		for _, n := range l.FieldOrder {
			if f, ok := r.Get(n); ok {
				if err := o.writeField(f, b); err != nil {
					return fmt.Errorf("error writing ADI record #%d: %v", i, err)
				}
				seen[f.Name] = true
			}
		}
		for _, f := range r.Fields() {
			if !seen[f.Name] {
				if err := o.writeField(f, b); err != nil {
					return fmt.Errorf("error writing ADI record #%d: %v", i, err)
				}
			}
		}
		if _, err := b.WriteString(fmt.Sprintf("<%s>%s", o.fixCase("EOR"), o.RecordSep.Val())); err != nil {
			return fmt.Errorf("error writing ADI record #%d: %v", i, err)
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
		return fmt.Errorf("error writing %s: %v", f, err)
	}
	return nil
}

func (o *ADIIO) fixCase(s string) string {
	if o.LowerCase {
		return strings.ToLower(s)
	}
	return strings.ToUpper(s)
}

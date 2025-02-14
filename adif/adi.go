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
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type ADIIO struct {
	LowerCase                  bool // TODO consider a case enum: keep, upper, lower, or just get rid of this option
	ASCIIOnly, AllowUnknownTag bool
	FieldSep, RecordSep        Separator
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
	if errors.Is(err, io.EOF) {
		if s != "" {
			l.Comment = s
		}
		return l, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error reading to first tag: %w", err)
	}
	// final byte is '<'
	comments = append(comments, strings.TrimSpace(s[0:len(s)-1]))
	// ADIF specification seems to imply that without a comment at the start
	// of a file, and thus the first character is '<', then there is no header
	// and the < starts the first record.  This invariant may not hold for all
	// software though, so allow an <EOH> even if we didn't get a comment.
	cur := NewRecord()
	var sawHeader, sawRecord bool
	for { // invariant: last byte read was '<'
		s, err = r.ReadString('>')
		if errors.Is(err, io.EOF) {
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
				l.AddRecord(cur)
				cur = NewRecord()
				comments = nil
			default:
				if o.AllowUnknownTag {
					comments = append(comments, fmt.Sprintf("[%s]", tag[0]))
				} else {
					return nil, fmt.Errorf("invalid ADI field without length <%s, see adi-allow-unknown-tags option", s)
				}
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
			if strings.HasPrefix(strings.ToUpper(tag[0]), "USERDEF") {
				if len(tag) != 3 {
					return nil, fmt.Errorf("missing type for %s field %q", tag[0], v)
				}
				fname, extra, hasextra := strings.Cut(string(v), ",")
				u := UserdefField{Name: fname}
				u.Type, err = DataTypeFromIndicator(tag[2])
				if err != nil {
					return nil, fmt.Errorf("%v from <%s", err, s)
				}
				if hasextra {
					if u.Type == TypeNumber {
						if n, err := fmt.Sscanf(extra, "{%f:%f}", &u.Min, &u.Max); err != nil || n != 2 {
							return nil, fmt.Errorf("invalid %s range %q", tag[0], extra)
						}
					} else {
						if !strings.HasPrefix(extra, "{") || !strings.HasSuffix(extra, "}") {
							return nil, fmt.Errorf("invalid %s enumeration list %q", tag[0], extra)
						}
						u.EnumValues = strings.Split(extra[1:len(extra)-1], ",")
					}
				}
				l.AddUserdef(u)
			} else {
				// spec says everything is ASCII, but this accepts UTF-8
				// as long as the tag length is accurate in bytes
				f := Field{Name: tag[0], Value: string(v)}
				if len(tag) == 3 {
					f.Type, err = DataTypeFromIndicator(tag[2])
					if err != nil {
						return nil, fmt.Errorf("%v from <%s", err, s)
					}
				}
				cur.Set(f)
			}
		default:
			return nil, fmt.Errorf("invalid ADI tag format <%s", s)
		}
		// arbitrary text between one field or record and the next
		c, err := r.ReadString('<')
		c = strings.TrimSpace(strings.TrimSuffix(c, "<"))
		if c != "" {
			comments = append(comments, c)
		}
		if errors.Is(err, io.EOF) {
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

const defaultAdiComment = "ADI format, see https://adif.org.uk/"

func (o *ADIIO) Write(l *Logfile, out io.Writer) error {
	if err := o.validate(l); err != nil {
		return err
	}
	b := bufio.NewWriter(out)
	if !l.Header.Empty() {
		c := l.Header.GetComment()
		if c == "" {
			c = defaultAdiComment
		}
		if err := o.writeComment(b, c, o.RecordSep.Val()); err != nil {
			return fmt.Errorf("writing ADI header: %w", err)
		}
		for _, f := range l.Header.Fields() {
			if err := o.writeField(f, b); err != nil {
				return fmt.Errorf("writing ADI header: %w", err)
			}
		}
		if err := o.writeUserdef(l.Userdef, b); err != nil {
			return fmt.Errorf("writing ADI header: %w", err)
		}
		if _, err := b.WriteString(fmt.Sprintf("<%s>%s", o.fixCase("EOH"), o.RecordSep.Val())); err != nil {
			return fmt.Errorf("writing ADI header: %w", err)
		}
	} else if len(l.Userdef) > 0 { // add a header for the userdef fields
		if err := o.writeComment(b, defaultAdiComment, o.RecordSep.Val()); err != nil {
			return fmt.Errorf("writing ADI header: %w", err)
		}
		if err := o.writeUserdef(l.Userdef, b); err != nil {
			return fmt.Errorf("writing ADI header: %w", err)
		}
		if _, err := b.WriteString(fmt.Sprintf("<%s>%s", o.fixCase("EOH"), o.RecordSep.Val())); err != nil {
			return fmt.Errorf("writing ADI header: %w", err)
		}
	}

	for i, r := range l.Records {
		if c := r.GetComment(); c != "" {
			if err := o.writeComment(b, c, o.FieldSep.Val()); err != nil {
				return fmt.Errorf("writing ADI record comment: %w", err)
			}
		}
		seen := make(map[string]bool)
		for _, n := range l.FieldOrder {
			if f, ok := r.Get(n); ok {
				if err := o.writeField(f, b); err != nil {
					return fmt.Errorf("writing ADI record #%d: %w", i, err)
				}
				seen[f.Name] = true
			}
		}
		for _, f := range r.Fields() {
			if !seen[f.Name] {
				if err := o.writeField(f, b); err != nil {
					return fmt.Errorf("writing ADI record #%d: %w", i, err)
				}
			}
		}
		if _, err := b.WriteString(fmt.Sprintf("<%s>%s", o.fixCase("EOR"), o.RecordSep.Val())); err != nil {
			return fmt.Errorf("writing ADI record #%d: %w", i, err)
		}
	}

	if l.Comment != "" {
		if err := o.writeComment(b, l.Comment, "\n"); err != nil {
			return fmt.Errorf("writing ADI comment: %w", err)
		}
	}
	return b.Flush()
}

func (o *ADIIO) writeField(f Field, b *bufio.Writer) error {
	val := ensureCRLF(f.Value)
	var tag string
	// TODO error if IntlString/non-ASCII, unless a flag allows
	if f.Type == TypeUnspecified {
		tag = fmt.Sprintf("<%s:%d>", o.fixCase(f.Name), len(val))
	} else {
		tag = fmt.Sprintf("<%s:%d:%s>", o.fixCase(f.Name), len(val), o.fixCase(f.Type.Indicator()))
	}
	if _, err := b.WriteString(fmt.Sprintf("%s%s%s", tag, val, o.FieldSep.Val())); err != nil {
		return fmt.Errorf("error writing %s: %w", f, err)
	}
	return nil
}

func (o *ADIIO) writeUserdef(us []UserdefField, b *bufio.Writer) error {
	for i, u := range us {
		if err := u.ValidateSelf(); err != nil {
			return err
		}
		f := Field{Name: fmt.Sprintf("USERDEF%d", i+1), Type: u.Type}
		if len(u.EnumValues) > 0 {
			f.Value = fmt.Sprintf("%s,{%s}", u.Name, strings.Join(u.EnumValues, ","))
		} else if u.Min != 0.0 || u.Max != 0.0 {
			f.Value = fmt.Sprintf("%s,%s", u.Name, formatRange(u))
		} else {
			f.Value = u.Name
		}
		if err := o.writeField(f, b); err != nil {
			return err
		}
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

func (o *ADIIO) validate(l *Logfile) error {
	check := func(r *Record) error {
		if r == nil { // in case l.Header is nil
			return nil
		}
		for _, f := range r.Fields() {
			if o.ASCIIOnly {
				for _, r := range f.Value {
					// spec limits Character to ASCII 32 to 126; MultilineString also allows CR/LF
					if (r < 32 || r > 126) && (r != '\r' && r != '\n') {
						return fmt.Errorf("non-ASCII character in %v", f)
					}
				}
			}
		}
		return nil
	}
	if err := check(l.Header); err != nil {
		return err
	}
	for _, r := range l.Records {
		if err := check(r); err != nil {
			return err
		}
	}
	return nil
}

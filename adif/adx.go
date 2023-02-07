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
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type adxField struct {
	XMLName   xml.Name
	Value     string `xml:",chardata"`
	Type      string `xml:"TYPE,attr,omitempty"`
	FieldID   string `xml:"FIELDID,attr,omitempty"`
	FieldName string `xml:"FIELDNAME,attr,omitempty"`
	ProgramID string `xml:"PROGRAMID,attr,omitempty"`
	Enum      string `xml:"ENUM,attr,omitempty"`
	Range     string `xml:"RANGE,attr,omitempty"`
	Comment   string `xml:",comment"`
}

func (f adxField) Field() Field {
	// TODO figure out user defined and app-specific fields
	return Field{Name: f.XMLName.Local, Value: f.Value, Type: typeIdentifiers[f.Type]}
}

func newAdxField(f Field) adxField {
	// TODO figure out user defined and app-specific fields
	return adxField{XMLName: xml.Name{Local: f.Name}, Value: f.Value, Type: f.Type.Identifier()}
}

type adxRecord struct {
	Comment string     `xml:",comment"`
	Fields  []adxField `xml:",any"`
}

func (r adxRecord) Record() *Record {
	res := NewRecord()
	fcs := make([]string, 0, len(r.Fields))
	if r.Comment != "" {
		fcs = append(fcs, r.Comment)
	}
	for _, f := range r.Fields {
		res.Set(f.Field())
		if f.Comment != "" {
			fcs = append(fcs, f.Comment)
		}
	}
	res.SetComment(strings.Join(fcs, "\n"))
	return res
}

func newAdxRecord(r *Record) adxRecord {
	res := adxRecord{}
	for _, f := range r.Fields() {
		res.Fields = append(res.Fields, newAdxField(f))
	}
	if c := r.GetComment(); c != "" {
		res.Comment = r.GetComment()
	}
	return res
}

type adxFile struct {
	Header  adxRecord   `xml:"HEADER"`
	Records []adxRecord `xml:"RECORDS>RECORD"`
	Comment string      `xml:",comment"`
	// The ADIF test QSO file has comments between <RECORD> tags (rather than
	// inside the tag), but a xml:",comment" field just gets a concatenated
	// string of comments with context to associate them back to a record.
	// TODO A custom xml.Unmarshaler interface for a recordList wrapper struct
	// might be able to handle this, if it's worth doing.
}

type ADXIO struct {
	Indent int
}

func NewADXIO() *ADXIO {
	return &ADXIO{}
}

func (_ *ADXIO) String() string { return "adx" }

func (o *ADXIO) Read(in io.Reader) (*Logfile, error) {
	l := NewLogfile()
	f := adxFile{}
	d := xml.NewDecoder(in)
	if err := d.Decode(&f); err != nil {
		return nil, fmt.Errorf("could not decode ADX file: %w", err)
	}
	l.Comment = f.Comment
	l.Header = f.Header.Record()
	for _, r := range f.Records {
		l.Records = append(l.Records, r.Record())
	}
	return l, nil
}

func (o *ADXIO) Write(l *Logfile, out io.Writer) error {
	f := adxFile{}
	f.Header = newAdxRecord(l.Header)
	for _, r := range l.Records {
		f.Records = append(f.Records, newAdxRecord(r))
	}
	if l.Comment != "" {
		f.Comment = l.Comment
	}
	if n, err := out.Write([]byte(xml.Header)); err != nil {
		return fmt.Errorf("could not write XML header to %s: %w", out, err)
	} else if n != len(xml.Header) {
		return fmt.Errorf("could not write XML header to %s: only wrote %d bytes", out, n)
	}
	e := xml.NewEncoder(out)
	e.Indent("", strings.Repeat(" ", o.Indent))
	start := xml.StartElement{Name: xml.Name{Local: "ADX"}}
	if err := e.EncodeElement(&f, start); err != nil {
		return fmt.Errorf("error writing ADX file %s: %w", l.Filename, err)
	}
	return nil
}

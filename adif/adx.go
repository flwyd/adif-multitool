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
	"strconv"
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

func (f adxField) IsUserdef() bool { return f.XMLName.Local == "USERDEF" }

func (f adxField) IsAppDefined() bool { return f.XMLName.Local == "APP" }

func (f adxField) Field() Field {
	dt, err := DataTypeFromIndicator(f.Type)
	if err != nil {
		dt = TypeUnspecified
	}
	if f.IsUserdef() {
		return Field{Name: f.FieldName, Value: f.Value, Type: dt}
	}
	if f.IsAppDefined() {
		return Field{Name: fmt.Sprintf("APP_%s_%s", f.ProgramID, f.FieldName), Value: f.Value, Type: dt}
	}
	return Field{Name: f.XMLName.Local, Value: f.Value, Type: dt}
}

func (f adxField) UserdefField() (UserdefField, error) {
	u := UserdefField{}
	if f.XMLName.Local != "USERDEF" {
		return u, fmt.Errorf("not a USERDEF field: %+v", f)
	}
	if f.Value == "" {
		return u, fmt.Errorf("USERDEF field without a name: %+v", f)
	}
	u.Name = f.Value
	if f.Type != "" {
		t, err := DataTypeFromIndicator(f.Type)
		if err != nil {
			return u, fmt.Errorf("%s: unknown type identifier %s: %w", f.FieldName, f.Type, err)
		}
		u.Type = t
	}
	if f.Range != "" {
		if n, err := fmt.Sscanf(f.Range, "{%f:%f}", &u.Min, &u.Max); n != 2 || err != nil {
			return u, fmt.Errorf("%s: invalid range %q", f.FieldName, f.Range)
		}
	}
	if f.Enum != "" {
		if !strings.HasPrefix(f.Enum, "{") || !strings.HasSuffix(f.Enum, "}") {
			return u, fmt.Errorf("%s: invalid enumeration %q", f.FieldName, f.Enum)
		}
		u.EnumValues = strings.Split(f.Enum[1:len(f.Enum)-1], ",")
	}
	return u, nil
}

func newAdxField(f Field) adxField {
	if f.IsAppDefined() {
		s := strings.SplitN(f.Name, "_", 3)
		return adxField{
			XMLName:   xml.Name{Local: "APP"},
			ProgramID: s[1],
			FieldName: s[2],
			Type:      f.Type.Indicator(),
			Value:     ensureCRLF(f.Value),
		}
	}
	return adxField{
		XMLName: xml.Name{Local: f.Name},
		Value:   ensureCRLF(f.Value),
		Type:    f.Type.Indicator(),
	}
}

func newAdxUserdefField(f Field) adxField { // for use in records
	return adxField{
		XMLName:   xml.Name{Local: "USERDEF"},
		FieldName: f.Name,
		Value:     ensureCRLF(f.Value),
		Type:      f.Type.Indicator()}
}

func newAdxUserdef(u UserdefField, id int) adxField { // for use in header
	f := adxField{
		XMLName: xml.Name{Local: "USERDEF"},
		FieldID: strconv.FormatInt(int64(id), 10),
		Value:   strings.ToUpper(u.Name),
		Type:    u.Type.Indicator()}
	if u.Min != 0.0 || u.Max != 0.0 {
		f.Range = formatRange(u)
	}
	if len(u.EnumValues) > 0 {
		f.Enum = fmt.Sprintf("{%s}", strings.Join(u.EnumValues, ","))
	}
	return f
}

type adxRecord struct {
	Comment string     `xml:",comment"`
	Fields  []adxField `xml:",any"`
}

func (r adxRecord) Record(header bool) *Record {
	res := NewRecord()
	fcs := make([]string, 0, len(r.Fields))
	if c := strings.TrimSpace(r.Comment); c != "" {
		fcs = append(fcs, c)
	}
	for _, f := range r.Fields {
		if header && f.IsUserdef() {
			continue
		}
		res.Set(f.Field())
		if c := strings.TrimSpace(f.Comment); c != "" {
			fcs = append(fcs, c)
		}
	}
	res.SetComment(strings.Join(fcs, "\n"))
	return res
}

func (r adxRecord) UserdefFields() ([]UserdefField, error) {
	var res []UserdefField
	for _, f := range r.Fields {
		if f.IsUserdef() {
			u, err := f.UserdefField()
			if err != nil {
				return nil, err
			}
			res = append(res, u)
		}
	}
	return res, nil
}

func newAdxRecord(r *Record, l *Logfile) adxRecord {
	res := adxRecord{}
	for _, f := range r.Fields() {
		var x adxField
		if _, ok := l.GetUserdef(f.Name); ok {
			x = newAdxUserdefField(f)
		} else {
			x = newAdxField(f)
		}
		res.Fields = append(res.Fields, x)
	}
	if c := strings.TrimSpace(r.GetComment()); c != "" {
		res.Comment = c
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
	l.Comment = strings.TrimSpace(f.Comment)
	l.Header = f.Header.Record(true)
	us, err := f.Header.UserdefFields()
	if err != nil {
		return nil, err
	}
	for _, u := range us {
		l.AddUserdef(u)
	}
	for _, r := range f.Records {
		l.AddRecord(r.Record(false))
	}
	return l, nil
}

func (o *ADXIO) Write(l *Logfile, out io.Writer) error {
	f := adxFile{}
	f.Header = newAdxRecord(l.Header, l)
	for i, u := range l.Userdef {
		if err := u.ValidateSelf(); err != nil {
			return err
		}
		f.Header.Fields = append(f.Header.Fields, newAdxUserdef(u, i+1))
	}
	for _, r := range l.Records {
		f.Records = append(f.Records, newAdxRecord(r, l))
	}
	if c := strings.TrimSpace(l.Comment); c != "" {
		f.Comment = c
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
	if _, err := out.Write([]byte("\n")); err != nil {
		return fmt.Errorf("error writing ADX file %s: %w", l.Filename, err)
	}
	return nil
}

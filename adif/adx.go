package adif

import (
	"encoding/xml"
	"fmt"
	"io"
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
	Fields  []adxField `xml:",any"`
	Comment string     `xml:",comment"`
}

func (r adxRecord) Record() *Record {
	fields := make([]Field, len(r.Fields))
	for i, f := range r.Fields {
		fields[i] = f.Field()
	}
	return NewRecord(fields...)
	// TODO set comment
}

func newAdxRecord(r *Record) adxRecord {
	res := adxRecord{}
	// TODO set comment
	for _, f := range r.Fields() {
		res.Fields = append(res.Fields, newAdxField(f))
	}
	return res
}

type adxFile struct {
	Header  adxRecord   `xml:"HEADER"`
	Records []adxRecord `xml:"RECORDS>RECORD"`
	Comment string      `xml:",comment"`
}

type ADXIO struct {
	// TODO
}

func NewADXIO() *ADXIO {
	return &ADXIO{}
}

func (_ *ADXIO) String() string { return "adx" }

func (_ *ADXIO) Read(in NamedReader) (*Logfile, error) {
	l := NewLogfile(in.Name())
	f := adxFile{}
	d := xml.NewDecoder(in)
	if err := d.Decode(&f); err != nil {
		return nil, fmt.Errorf("could not decode ADX file %s: %w", in.Name(), err)
	}
	l.Comment = f.Comment
	l.Header = f.Header.Record()
	for _, r := range f.Records {
		l.Records = append(l.Records, r.Record())
	}
	return l, nil
}

func (_ *ADXIO) Write(l *Logfile, out io.Writer) error {
	f := adxFile{}
	f.Header = newAdxRecord(l.Header)
	for _, r := range l.Records {
		f.Records = append(f.Records, newAdxRecord(r))
	}
	if n, err := out.Write([]byte(xml.Header)); err != nil {
		return fmt.Errorf("could not write XML header to %s: %w", out, err)
	} else if n != len(xml.Header) {
		return fmt.Errorf("could not write XML header to %s: only wrote %d bytes", out, n)
	}
	e := xml.NewEncoder(out)
	e.Indent("", " ")
	start := xml.StartElement{Name: xml.Name{Local: "ADX"}}
	if err := e.EncodeElement(&f, start); err != nil {
		return fmt.Errorf("error writing ADX file %s: %w", l.Filename, err)
	}
	return nil
}

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
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TODO test a file without a header (first character is <)

func TestEmptyADI(t *testing.T) {
	input := "Created on 2001-02-03 at 4:05\n"
	adi := NewADIIO()
	if parsed, err := adi.Read(strings.NewReader(input)); err != nil {
		t.Errorf("Read(%q) got error %v", input, err)
	} else {
		if got := len(parsed.Records); got != 0 {
			t.Errorf("Read(%q) got %d records, want %d", input, got, 0)
		}
	}
}

func TestReadADI(t *testing.T) {
	input := `Generated today <ADIF_VER:5>3.1.4 <CREATED_TIMESTAMP:15>20220102 153456 <PROGRAMID:11>adi_test <USERDEF1:8:S>My Field <PROGRAMVERSION:5>1.2.3
<USERDEF2:19:E>SweaterSize,{S,M,L} <userdef3:15:N>shoesize,{5:20} <EOH>
<QSO_DATE:8>19901031 <TIME_ON:4>1234  <BAND:3>40M<CALLSIGN:4>W1AW
<NAME:17>Hiram Percy Maxim <APP_MONOLOG_birth_day:8:D>18690902 <EOR>

<qso_datE:8:D>20221224
Field comment #1 <time_ON:6:T>095846
<band:6:E>1.25cm
<callsign:3:S>N0P
Field comment #2 <name:11:S>Santa Claus
<MY field:12>{!@#}, ($%^)
<eor>
<QSO_date:8>19190219
<APP_monolog_BIRTH_DAY:8>18960815
<APP_adifmt_BIRTH_DAY:15:S>August 15, 1896
<RIG:82:M>100 watt C.W.
Armstrong regenerative circuit
Inverted L antenna, 70' above ground
<FREQ:5:N>7.654
<CaLlSiGn:3:S>1AY
This is a random comment
<name:12:s>"C.G." Tuska
<SWEATERSIZE:1>L<shoeSize:2>12<eOr>
Comment at &lt;end&gt; of file.
`
	wantFields := [][]Field{
		{
			{Name: "QSO_DATE", Value: "19901031"},
			{Name: "TIME_ON", Value: "1234"},
			{Name: "BAND", Value: "40M"},
			{Name: "CALLSIGN", Value: "W1AW"},
			{Name: "NAME", Value: "Hiram Percy Maxim"},
			{Name: "APP_MONOLOG_BIRTH_DAY", Value: "18690902", Type: TypeDate},
		},
		{
			{Name: "QSO_DATE", Value: "20221224", Type: TypeDate},
			{Name: "TIME_ON", Value: "095846", Type: TypeTime},
			{Name: "BAND", Value: "1.25cm", Type: TypeEnumeration},
			{Name: "CALLSIGN", Value: "N0P", Type: TypeString},
			{Name: "NAME", Value: "Santa Claus", Type: TypeString},
			{Name: "MY FIELD", Value: "{!@#}, ($%^)"},
		},
		{
			{Name: "QSO_DATE", Value: "19190219"},
			{Name: "APP_MONOLOG_BIRTH_DAY", Value: "18960815"},
			{Name: "APP_ADIFMT_BIRTH_DAY", Value: "August 15, 1896", Type: TypeString},
			{Name: "RIG", Value: `100 watt C.W.
Armstrong regenerative circuit
Inverted L antenna, 70' above ground
`, Type: TypeMultilineString},
			{Name: "FREQ", Value: "7.654", Type: TypeNumber},
			{Name: "CALLSIGN", Value: "1AY", Type: TypeString},
			{Name: "NAME", Value: `"C.G." Tuska`, Type: TypeString},
			{Name: "SWEATERSIZE", Value: "L"},
			{Name: "SHOESIZE", Value: "12"},
		},
	}

	wantComments := []string{
		"",
		"Field comment #1\nField comment #2",
		"This is a random comment",
	}
	adi := NewADIIO()
	if parsed, err := adi.Read(strings.NewReader(input)); err != nil {
		t.Errorf("Read(%q) got error %v", input, err)
	} else {
		for i, r := range parsed.Records {
			fields := r.Fields()
			if diff := cmp.Diff(wantFields[i], fields); diff != "" {
				t.Errorf("Read(%q) record %d did not match expected, diff:\n%s", input, i+1, diff)
			}
			if got := r.GetComment(); got != wantComments[i] {
				t.Errorf("Read(%q) record %d got comment %q want %q", input, i+1, got, wantComments[i])
			}
		}
		if gotlen := len(parsed.Records); gotlen != len(wantFields) {
			t.Errorf("Read(%q) got %d records:\n%v\nwant %d\n%v", input, gotlen, parsed.Records[len(wantFields):], len(wantFields), wantFields)
		}
		if want := "Comment at &lt;end&gt; of file."; parsed.Comment != want {
			t.Errorf("Read(%q) got logfile comment %q, want %q", input, parsed.Comment, want)
		}
	}
}

func TestWriteADI(t *testing.T) {
	l := NewLogfile()
	l.Comment = "The <last> word."
	l.AddRecord(NewRecord(
		Field{Name: "QSO_DATE", Value: "19901031", Type: TypeDate},
		Field{Name: "TIME_ON", Value: "1234", Type: TypeTime},
		Field{Name: "BAND", Value: "40M"},
		Field{Name: "CALLSIGN", Value: "W1AW"},
		Field{Name: "NAME", Value: "Hiram Percy Maxim", Type: TypeString},
		Field{Name: "APP_MONOLOG_BIRTH_DAY", Value: "18690902", Type: TypeDate},
	)).AddRecord(NewRecord(
		Field{Name: "QSO_DATE", Value: "20221224"},
		Field{Name: "TIME_ON", Value: "095846"},
		Field{Name: "BAND", Value: "1.25cm", Type: TypeEnumeration},
		Field{Name: "CALLSIGN", Value: "N0P", Type: TypeString},
		Field{Name: "NAME", Value: "Santa Claus"},
		Field{Name: "My Field", Value: "{!@#}, ($%^)"},
	)).AddRecord(NewRecord(
		Field{Name: "QSO_DATE", Value: "19190219", Type: TypeDate},
		Field{Name: "APP_MONOLOG_BIRTH_DAY", Value: "18960815"},
		Field{Name: "APP_ADIFMT_BIRTH_DAY", Value: "August 15, 1896", Type: TypeString},
		Field{Name: "RIG", Value: `100 watt C.W.
Armstrong regenerative circuit
Inverted L antenna, 70' above ground
`, Type: TypeMultilineString},
		Field{Name: "FREQ", Value: "7.654", Type: TypeNumber},
		Field{Name: "CALLSIGN", Value: "1AY", Type: TypeString},
		Field{Name: "NAME", Value: `"C.G." Tuska`, Type: TypeString},
		Field{Name: "SweaterSize", Value: "L"},
		Field{Name: "SHOESIZE", Value: "12"},
	))
	l.Records[1].SetComment("Record comment")
	l.Header.Set(Field{Name: "ADIF_VER", Value: "3.1.4"})
	l.Header.Set(Field{Name: "PROGRAMID", Value: "adi_test"})
	l.Header.Set(Field{Name: "PROGRAMVERSION", Value: "1.2.3"})
	l.Header.Set(Field{Name: "CREATED_TIMESTAMP", Value: "20220102 153456"})
	l.Userdef = []UserdefField{
		{Name: "MY FIELD", Type: TypeString},
		{Name: "sweatersize", Type: TypeEnumeration, EnumValues: []string{"S", "M", "L"}},
		{Name: "ShoeSize", Type: TypeNumber, Min: 5, Max: 20},
	}
	want := `ADI format, see https://adif.org.uk/
<ADIF_VER:5>3.1.4 <PROGRAMID:8>adi_test <PROGRAMVERSION:5>1.2.3 <CREATED_TIMESTAMP:15>20220102 153456 <USERDEF1:8:S>MY FIELD <USERDEF2:19:E>sweatersize,{S,M,L} <USERDEF3:15:N>ShoeSize,{5:20} <EOH>
<QSO_DATE:8:D>19901031 <TIME_ON:4:T>1234 <BAND:3>40M <CALLSIGN:4>W1AW <NAME:17:S>Hiram Percy Maxim <APP_MONOLOG_BIRTH_DAY:8:D>18690902 <EOR>
Record comment <QSO_DATE:8>20221224 <TIME_ON:6>095846 <BAND:6:E>1.25cm <CALLSIGN:3:S>N0P <NAME:11>Santa Claus <MY FIELD:12>{!@#}, ($%^) <EOR>
<QSO_DATE:8:D>19190219 <APP_MONOLOG_BIRTH_DAY:8>18960815 <APP_ADIFMT_BIRTH_DAY:15:S>August 15, 1896 <RIG:85:M>100 watt C.W.
Armstrong regenerative circuit
Inverted L antenna, 70' above ground
 <FREQ:5:N>7.654 <CALLSIGN:3:S>1AY <NAME:12:S>"C.G." Tuska <SWEATERSIZE:1>L <SHOESIZE:2>12 <EOR>
The &lt;last&gt; word.
`
	for _, r := range []string{"watt C.W.", "regenerative circuit", "above ground"} {
		want = strings.Replace(want, r+"\n", r+"\r\n", 1)
	}
	adi := NewADIIO()
	adi.RecordSep = SeparatorNewline
	adi.FieldSep = SeparatorSpace
	out := &strings.Builder{}
	if err := adi.Write(l, out); err != nil {
		t.Errorf("Write(%v) got error %v", l, err)
	} else {
		got := out.String()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Write(%v) had diff with expected:\n%s", l, diff)
		}
	}
}

func TestADIASCIIOnly(t *testing.T) {
	adi := NewADIIO()
	adi.ASCIIOnly = true
	tests := []string{"\t", "\x00", "\x7F", "\x80", "Baño", "Baþ", "नहाना", "욕조"}
	for _, s := range tests {
		l := NewLogfile()
		l.AddRecord(NewRecord(Field{Name: "CALL", Value: "N0P"},
			Field{Name: "APP_TEST_ASCII", Value: s, Type: TypeString}))
		out := &strings.Builder{}
		if err := adi.Write(l, out); err == nil {
			t.Errorf("adi.Write with field value %q got %q, want error", s, out)
		} else if out.String() != "" {
			t.Errorf("adi.Write with non-ASCII character wrote partial output %q", out)
		}
		lh := NewLogfile()
		lh.Header = NewRecord(Field{Name: "APP_TEST_HEADER", Value: s})
		out = &strings.Builder{}
		if err := adi.Write(l, out); err == nil {
			t.Errorf("adi.Write with header value %q got %s, want error", s, out)
		} else if out.String() != "" {
			t.Errorf("adi.Write with non-ASCII character in header wrote partial output %q", out)
		}
	}
}

func TestADIUnknownTag(t *testing.T) {
	// See https://groups.io/g/adifdev/topic/angle_brackets_outside_of/111067202 for context
	input := `Header comment with <a tag> inside
<ADIF_VER:5>3.1.5 <CREATED_TIMESTAMP:15>20220102 153456 <PROGRAMID:8>adi_test <EOH>
Record comment <QSO_DATE:8>19901031 <TIME_ON:4>1234 Field <tag> comment <BAND:3>40M <CALLSIGN:4>W1AW <EOR>
<unknown_tag_1> <QSO_DATE:8>20210203 <TIME_ON:4>0405 <BAND:2>2m <CALLSIGN:3>N0X End record comment<EOR>

<APP_LoTW_EOF>
`

	for _, b := range []bool{false, true} {
		t.Run(fmt.Sprintf("allow_%t", b), func(t *testing.T) {
			adi := NewADIIO()
			adi.AllowUnknownTag = b
			parsed, err := adi.Read(strings.NewReader(input))
			if b {
				if err != nil {
					t.Errorf("Read(%q) with AllowUnknownTag got error %v", input, err)
				} else {
					if want := "[APP_LoTW_EOF]"; parsed.Comment != want {
						t.Errorf("tag not handled, want %q got %q", want, parsed.Comment)
					}
					if want := "Header comment with\n[a tag]\ninside"; parsed.Header.GetComment() != want {
						t.Errorf("tag not handled, want %q got %q", want, parsed.Header.GetComment())
					}
					if want := "Record comment\nField\n[tag]\ncomment"; parsed.Records[0].GetComment() != want {
						t.Errorf("tag not handled, want %q got %q", want, parsed.Records[0].GetComment())
					}
					if want := "[unknown_tag_1]\nEnd record comment"; parsed.Records[1].GetComment() != want {
						t.Errorf("tag not handled, want %q got %q", want, parsed.Records[1].GetComment())
					}
					r := NewRecord(Field{Name: "QSO_DATE", Value: "19901031"}, Field{Name: "TIME_ON", Value: "1234"}, Field{Name: "BAND", Value: "40M"}, Field{Name: "CALLSIGN", Value: "W1AW"})
					if !r.Equal(parsed.Records[0]) {
						t.Errorf("wrong record 0 parse, want %v got %v", r, parsed.Records[0])
					}
					r = NewRecord(Field{Name: "QSO_DATE", Value: "20210203"}, Field{Name: "TIME_ON", Value: "0405"}, Field{Name: "BAND", Value: "2m"}, Field{Name: "CALLSIGN", Value: "N0X"})
					if !r.Equal(parsed.Records[1]) {
						t.Errorf("wrong record 1 parse, want %v got %v", r, parsed.Records[1])
					}
				}
			} else {
				if err == nil {
					t.Errorf("Read(%q) with AllowUnknownTag=false did not get error:\n%v", input, parsed)
				}
			}
		})
	}
}

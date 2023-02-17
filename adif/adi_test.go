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
<NAME:17>Hiram Percy Maxim <EOR>

<qso_datE:8:D>20221224
Field comment #1 <time_ON:6:T>095846
<band:6:E>1.25cm
<callsign:3:S>N0P
Field comment #2 <name:11:S>Santa Claus
<MY field:12>{!@#}, ($%^)
<eor>
<QSO_date:8>19190219
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
				t.Errorf("Read(%q) record %d did not match expected, diff:\n%s", input, i, diff)
			}
			if got := r.GetComment(); got != wantComments[i] {
				t.Errorf("Read(%q) record %d got comment %q want %q", input, i, got, wantComments[i])
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
	l.Records = append(l.Records, NewRecord(
		Field{Name: "QSO_DATE", Value: "19901031", Type: TypeDate},
		Field{Name: "TIME_ON", Value: "1234", Type: TypeTime},
		Field{Name: "BAND", Value: "40M"},
		Field{Name: "CALLSIGN", Value: "W1AW"},
		Field{Name: "NAME", Value: "Hiram Percy Maxim", Type: TypeString},
	))
	l.Records = append(l.Records, NewRecord(
		Field{Name: "QSO_DATE", Value: "20221224"},
		Field{Name: "TIME_ON", Value: "095846"},
		Field{Name: "BAND", Value: "1.25cm", Type: TypeEnumeration},
		Field{Name: "CALLSIGN", Value: "N0P", Type: TypeString},
		Field{Name: "NAME", Value: "Santa Claus"},
		Field{Name: "My Field", Value: "{!@#}, ($%^)"},
	))
	l.Records[len(l.Records)-1].SetComment("Record comment")
	l.Records = append(l.Records, NewRecord(
		Field{Name: "QSO_DATE", Value: "19190219", Type: TypeDate},
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
<QSO_DATE:8:D>19901031 <TIME_ON:4:T>1234 <BAND:3>40M <CALLSIGN:4>W1AW <NAME:17:S>Hiram Percy Maxim <EOR>
Record comment <QSO_DATE:8>20221224 <TIME_ON:6>095846 <BAND:6:E>1.25cm <CALLSIGN:3:S>N0P <NAME:11>Santa Claus <MY FIELD:12>{!@#}, ($%^) <EOR>
<QSO_DATE:8:D>19190219 <RIG:82:M>100 watt C.W.
Armstrong regenerative circuit
Inverted L antenna, 70' above ground
 <FREQ:5:N>7.654 <CALLSIGN:3:S>1AY <NAME:12:S>"C.G." Tuska <SWEATERSIZE:1>L <SHOESIZE:2>12 <EOR>
The &lt;last&gt; word.
`
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

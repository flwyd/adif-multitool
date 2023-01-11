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
	input := StringSource{Filename: "empty.adi",
		Reader: strings.NewReader("Created on 2001-02-03 at 4:05\n")}
	adi := NewADIIO()
	if parsed, err := adi.Read(input); err != nil {
		t.Errorf("Read(%v) got error %v", input, err)
	} else {
		if got := len(parsed.Records); got != 0 {
			t.Errorf("Read(%v) got %d records, want %d", input, got, 0)
		}
	}
}

func TestReadADI(t *testing.T) {
	input := StringSource{Filename: "test.adi", Reader: strings.NewReader(
		`Generated today <ADIF_VER:5>3.1.4 <CREATED_TIMESTAMP:15>20220102 153456 <PROGRAMID:11>adi_test <PROGRAMVERSION:5>1.2.3 <EOH>
<QSO_DATE:8>19901031 <TIME_ON:4>1234  <BAND:3>40M<CALLSIGN:4>W1AW
<NAME:18>Hiram Percey Maxim <EOR>

<qso_datE:8:D>20221224
<time_ON:6:T>095846
<band:6:E>1.25cm
<callsign:3:S>N0P
<name:11:S>Santa Claus
<eor>
<QSO_date:8>19190219
<RIG:82:M>100 watt C.W.
Armstrong regenerative circuit
Inverted L antenna, 70' above ground
<FREQ:5:N>7.654
<CaLlSiGn:3:S>1AY
This is a random comment
<name:12:s>"C.G." Tuska<eOr>
`)}
	wantFields := [][]Field{
		{
			{Name: "QSO_DATE", Value: "19901031"},
			{Name: "TIME_ON", Value: "1234"},
			{Name: "BAND", Value: "40M"},
			{Name: "CALLSIGN", Value: "W1AW"},
			{Name: "NAME", Value: "Hiram Percey Maxim"},
		},
		{
			{Name: "QSO_DATE", Value: "20221224", Type: Date},
			{Name: "TIME_ON", Value: "095846", Type: Time},
			{Name: "BAND", Value: "1.25cm", Type: Enumeration},
			{Name: "CALLSIGN", Value: "N0P", Type: String},
			{Name: "NAME", Value: "Santa Claus", Type: String},
		},
		{
			{Name: "QSO_DATE", Value: "19190219"},
			{Name: "RIG", Value: `100 watt C.W.
Armstrong regenerative circuit
Inverted L antenna, 70' above ground
`, Type: MultilineString},
			{Name: "FREQ", Value: "7.654", Type: Number},
			{Name: "CALLSIGN", Value: "1AY", Type: String},
			{Name: "NAME", Value: `"C.G." Tuska`, Type: String},
		},
	}
	adi := NewADIIO()
	if parsed, err := adi.Read(input); err != nil {
		t.Errorf("Read(%v) got error %v", input, err)
	} else {
		for i, r := range parsed.Records {
			fields := r.Fields()
			if diff := cmp.Diff(wantFields[i], fields); diff != "" {
				t.Errorf("Read(%v) record %d did not match expected, diff:\n%s", input, i, diff)
			}
		}
		if gotlen := len(parsed.Records); gotlen != len(wantFields) {
			t.Errorf("Read(%v) got %d records:\n%v\nwant %d\n%v", input, gotlen, parsed.Records[len(wantFields):], len(wantFields), wantFields)
		}
		if parsed.Filename != input.Filename {
			t.Errorf("Read(%v) got Filename %q, want %q", input, parsed.Filename, input.Filename)
		}
	}
}

func TestWriteADI(t *testing.T) {
	l := NewLogfile("test-logfile")
	l.Records = append(l.Records, NewRecord(
		Field{Name: "QSO_DATE", Value: "19901031", Type: Date},
		Field{Name: "TIME_ON", Value: "1234", Type: Time},
		Field{Name: "BAND", Value: "40M"},
		Field{Name: "CALLSIGN", Value: "W1AW"},
		Field{Name: "NAME", Value: "Hiram Percey Maxim", Type: String},
	))
	l.Records = append(l.Records, NewRecord(
		Field{Name: "QSO_DATE", Value: "20221224"},
		Field{Name: "TIME_ON", Value: "095846"},
		Field{Name: "BAND", Value: "1.25cm", Type: Enumeration},
		Field{Name: "CALLSIGN", Value: "N0P", Type: String},
		Field{Name: "NAME", Value: "Santa Claus"},
	))
	l.Records = append(l.Records, NewRecord(
		Field{Name: "QSO_DATE", Value: "19190219", Type: Date},
		Field{Name: "RIG", Value: `100 watt C.W.
Armstrong regenerative circuit
Inverted L antenna, 70' above ground
`, Type: MultilineString},
		Field{Name: "FREQ", Value: "7.654", Type: Number},
		Field{Name: "CALLSIGN", Value: "1AY", Type: String},
		Field{Name: "NAME", Value: `"C.G." Tuska`, Type: String},
	))
	l.Header.Set(Field{Name: "ADIF_VER", Value: "3.1.4"})
	l.Header.Set(Field{Name: "PROGRAMID", Value: "adi_test"})
	l.Header.Set(Field{Name: "PROGRAMVERSION", Value: "1.2.3"})
	l.Header.Set(Field{Name: "CREATED_TIMESTAMP", Value: "20220102 153456"})
	want := `ADI comment at the top of the file
<ADIF_VER:5>3.1.4 <PROGRAMID:8>adi_test <PROGRAMVERSION:5>1.2.3 <CREATED_TIMESTAMP:15>20220102 153456 <EOH>
<QSO_DATE:8:D>19901031 <TIME_ON:4:T>1234 <BAND:3>40M <CALLSIGN:4>W1AW <NAME:18:S>Hiram Percey Maxim <EOR>
<QSO_DATE:8>20221224 <TIME_ON:6>095846 <BAND:6:E>1.25cm <CALLSIGN:3:S>N0P <NAME:11>Santa Claus <EOR>
<QSO_DATE:8:D>19190219 <RIG:82:M>100 watt C.W.
Armstrong regenerative circuit
Inverted L antenna, 70' above ground
 <FREQ:5:N>7.654 <CALLSIGN:3:S>1AY <NAME:12:S>"C.G." Tuska <EOR>
`
	adi := NewADIIO()
	adi.RecordSep = SeparatorNewline
	adi.FieldSep = SeparatorSpace
	adi.HeaderCommentFn = func(l *Logfile) string { return "ADI comment at the top of the file" }
	out := &strings.Builder{}
	if err := adi.Write(l, out); err != nil {
		t.Errorf("Read(%v) got error %v", l, err)
	} else {
		got := out.String()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Write(%v) had diff with expected:\n%s", l, diff)
		}
	}
}

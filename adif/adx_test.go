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
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEmptyADX(t *testing.T) {
	tests := []string{
		xml.Header + "<ADX><HEADER></HEADER><RECORDS></RECORDS></ADX>",
		xml.Header + "<ADX><RECORDS></RECORDS></ADX>",
		xml.Header + "<ADX></ADX>",
		"<ADX></ADX>",
	}
	for _, tc := range tests {
		input := StringReader{Filename: "empty.adx",
			Reader: strings.NewReader(tc)}
		adx := NewADXIO()
		if parsed, err := adx.Read(input); err != nil {
			t.Errorf("Read(%v) got error %v", input, err)
		} else {
			if got := len(parsed.Records); got != 0 {
				t.Errorf("Read(%v) got %d records, want %d", input, got, 0)
			}
		}
	}
}

func TestReadADX(t *testing.T) {
	input := StringReader{Filename: "test.adx", Reader: strings.NewReader(xml.Header +
		`<ADX>
		<HEADER>
		<ADIF_VER>3.1.4</ADIF_VER>
		<CREATED_TIMESTAMP>20220102 153456</CREATED_TIMESTAMP>
		<PROGRAMID>adx_test</PROGRAMID>
		<PROGRAMVERSION>1.2.3</PROGRAMVERSION>
		</HEADER>
		<RECORDS>
		<RECORD>
<QSO_DATE>19901031</QSO_DATE> <TIME_ON>1234</TIME_ON>  <BAND>40M</BAND><CALLSIGN>W1AW</CALLSIGN>
<NAME>Hiram Percey Maxim</NAME> </RECORD>
<RECORD>
	<QSO_DATE TYPE="D">20221224</QSO_DATE>
	<TIME_ON TYPE="T">095846</TIME_ON>
	<BAND TYPE="E">1.25cm</BAND>
	<CALLSIGN TYPE="S">N0P</CALLSIGN>
	<NAME TYPE="S">Santa Claus</NAME>
</RECORD>
<RECORD>
<QSO_DATE>19190219</QSO_DATE>
<RIG TYPE="M">100 watt C.W.
Armstrong regenerative circuit
Inverted L antenna, 70' above ground
</RIG><FREQ TYPE="N">7.654</FREQ>
<CALLSIGN TYPE="S">1AY</CALLSIGN>
<!-- This is a random comment -->
<NAME TYPE="S">"C.G." Tuska</NAME></RECORD>
</RECORDS>
</ADX>
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
	adi := NewADXIO()
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

func TestWriteADX(t *testing.T) {
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
	l.Header.Set(Field{Name: "PROGRAMID", Value: "adx_test"})
	l.Header.Set(Field{Name: "PROGRAMVERSION", Value: "1.2.3"})
	l.Header.Set(Field{Name: "CREATED_TIMESTAMP", Value: "20220102 153456"})
	want := xml.Header + `<ADX>
  <HEADER>
    <ADIF_VER>3.1.4</ADIF_VER>
    <PROGRAMID>adx_test</PROGRAMID>
    <PROGRAMVERSION>1.2.3</PROGRAMVERSION>
    <CREATED_TIMESTAMP>20220102 153456</CREATED_TIMESTAMP>
  </HEADER>
  <RECORDS>
    <RECORD>
      <QSO_DATE TYPE="D">19901031</QSO_DATE>
      <TIME_ON TYPE="T">1234</TIME_ON>
      <BAND>40M</BAND>
      <CALLSIGN>W1AW</CALLSIGN>
      <NAME TYPE="S">Hiram Percey Maxim</NAME>
    </RECORD>
    <RECORD>
      <QSO_DATE>20221224</QSO_DATE>
      <TIME_ON>095846</TIME_ON>
      <BAND TYPE="E">1.25cm</BAND>
      <CALLSIGN TYPE="S">N0P</CALLSIGN>
      <NAME>Santa Claus</NAME>
    </RECORD>
    <RECORD>
      <QSO_DATE TYPE="D">19190219</QSO_DATE>
      <RIG TYPE="M">100 watt C.W.&#xA;Armstrong regenerative circuit&#xA;Inverted L antenna, 70&#39; above ground&#xA;</RIG>
      <FREQ TYPE="N">7.654</FREQ>
      <CALLSIGN TYPE="S">1AY</CALLSIGN>
      <NAME TYPE="S">&#34;C.G.&#34; Tuska</NAME>
    </RECORD>
  </RECORDS>
</ADX>`
	adx := NewADXIO()
	adx.Indent = 2
	out := &strings.Builder{}
	if err := adx.Write(l, out); err != nil {
		t.Errorf("Read(%v) got error %v", l, err)
	} else {
		got := out.String()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Write(%v) had diff with expected:\n%s", l, diff)
		}
	}
}

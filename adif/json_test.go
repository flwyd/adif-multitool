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

func TestEmptyJson(t *testing.T) {
	tests := []string{
		// "",
		"{}",
		`{"RECORDS": []}`,
		`{"HEADER": {}, "RECORDS": []}`,
	}
	for _, tc := range tests {
		json := NewJSONIO()
		if parsed, err := json.Read(strings.NewReader(tc)); err != nil {
			t.Errorf("Read(%q) got error %v", tc, err)
		} else {
			if got := len(parsed.Records); got != 0 {
				t.Errorf("Read(%q) got %d records, want %d", tc, got, 0)
			}
		}
	}
}

func TestReadJSON(t *testing.T) {
	input := `{
    "HEADER": {
      "ADIF_VER": "3.1.4",
      "CREATED_TIMESTAMP": "20220102 153456",
      "PROGRAMID": "adx_test",
      "PROGRAMVERSION": "1.2.3"
    },
		"RECORDS": [
		{
"QSO_DATE": "19901031", "TIME_ON": "1234",  "BAND": "40M","CALLSIGN": "W1AW",
"NAME": "Hiram Percy Maxim" },
{
	"QSO_DATE": 20221224,
	"TIME_ON": "095846",
	"BAND": "1.25cm",
	"CALLSIGN": "N0P",
	"NAME": "Santa Claus",
	"QSO_RANDOM": false
},
{
"QSO_DATE": "19190219",
"RIG": "100 watt C.W.\nArmstrong regenerative circuit\nInverted L antenna, 70' above ground\n",
"FREQ": 7.654,
"CALLSIGN": "1AY",
"NAME": "\"C.G.\" Tuska",
"SILENT_KEY": true}
]
}
`
	wantFields := []*Record{
		NewRecord(Field{Name: "QSO_DATE", Value: "19901031"},
			Field{Name: "TIME_ON", Value: "1234"},
			Field{Name: "BAND", Value: "40M"},
			Field{Name: "CALLSIGN", Value: "W1AW"},
			Field{Name: "NAME", Value: "Hiram Percy Maxim"},
		),
		NewRecord(Field{Name: "QSO_DATE", Value: "20221224", Type: TypeNumber},
			Field{Name: "TIME_ON", Value: "095846"},
			Field{Name: "BAND", Value: "1.25cm"},
			Field{Name: "CALLSIGN", Value: "N0P"},
			Field{Name: "NAME", Value: "Santa Claus"},
			Field{Name: "QSO_RANDOM", Value: "N", Type: TypeBoolean},
		),
		NewRecord(Field{Name: "QSO_DATE", Value: "19190219"},
			Field{Name: "RIG", Value: `100 watt C.W.
Armstrong regenerative circuit
Inverted L antenna, 70' above ground
`},
			Field{Name: "FREQ", Value: "7.654", Type: TypeNumber},
			Field{Name: "CALLSIGN", Value: "1AY"},
			Field{Name: "NAME", Value: `"C.G." Tuska`},
			Field{Name: "SILENT_KEY", Value: "Y", Type: TypeBoolean},
		),
	}
	json := NewJSONIO()
	if parsed, err := json.Read(strings.NewReader(input)); err != nil {
		t.Errorf("Read(%q) got error %v", input, err)
	} else {
		if diff := cmp.Diff(wantFields, parsed.Records); diff != "" {
			t.Errorf("Read(%q) got diff:\n%s", input, diff)
		}
	}
}

func TestWriteJSON(t *testing.T) {
	l := NewLogfile()
	l.Comment = "JSON ignores comments"
	l.AddRecord(NewRecord(
		Field{Name: "QSO_DATE", Value: "19901031", Type: TypeDate},
		Field{Name: "TIME_ON", Value: "1234", Type: TypeTime},
		Field{Name: "BAND", Value: "40M"},
		Field{Name: "CALLSIGN", Value: "W1AW"},
		Field{Name: "NAME", Value: "Hiram Percy Maxim", Type: TypeString},
	)).AddRecord(NewRecord(
		Field{Name: "QSO_DATE", Value: "20221224"},
		Field{Name: "TIME_ON", Value: "095846"},
		Field{Name: "BAND", Value: "1.25cm", Type: TypeEnumeration},
		Field{Name: "QSO_RANDOM", Value: "N", Type: TypeBoolean},
		Field{Name: "CALLSIGN", Value: "N0P", Type: TypeString},
		Field{Name: "NAME", Value: "Santa Claus"},
	)).AddRecord(NewRecord(
		Field{Name: "QsO_dAtE", Value: "19190219", Type: TypeDate},
		Field{Name: "RIG", Value: `100 watt C.W.
Armstrong regenerative circuit
Inverted L antenna, < 70' above ground
`, Type: TypeMultilineString},
		Field{Name: "FREQ", Value: "7.654", Type: TypeNumber},
		Field{Name: "silent_key", Value: "Y", Type: TypeBoolean},
		Field{Name: "callsign", Value: "1AY", Type: TypeString},
		Field{Name: "NAME", Value: `"C.G." Tuska`, Type: TypeString},
	))
	l.Records[1].SetComment("Record comment")
	l.Header.Set(Field{Name: "adif_ver", Value: "3.1.4"})
	l.Header.Set(Field{Name: "PROGRAMID", Value: "adx_test"})
	l.Header.Set(Field{Name: "PROGRAMVERSION", Value: "1.2.3"})
	l.Header.Set(Field{Name: "CREATED_TIMESTAMP", Value: "20220102 153456"})
	want := `{
 "HEADER": {
  "ADIF_VER": "3.1.4",
  "CREATED_TIMESTAMP": "20220102 153456",
  "PROGRAMID": "adx_test",
  "PROGRAMVERSION": "1.2.3"
 },
 "RECORDS": [
  {
   "BAND": "40M",
   "CALLSIGN": "W1AW",
   "NAME": "Hiram Percy Maxim",
   "QSO_DATE": "19901031",
   "TIME_ON": "1234"
  },
  {
   "BAND": "1.25cm",
   "CALLSIGN": "N0P",
   "NAME": "Santa Claus",
   "QSO_DATE": "20221224",
   "QSO_RANDOM": false,
   "TIME_ON": "095846"
  },
  {
   "CALLSIGN": "1AY",
   "FREQ": 7.654,
   "NAME": "\"C.G.\" Tuska",
   "QSO_DATE": "19190219",
   "RIG": "100 watt C.W.\nArmstrong regenerative circuit\nInverted L antenna, < 70' above ground\n",
   "SILENT_KEY": true
  }
 ]
}
`
	json := NewJSONIO()
	json.Indent = 1
	json.HTMLSafe = false
	json.TypedOutput = true
	out := &strings.Builder{}
	if err := json.Write(l, out); err != nil {
		t.Errorf("Write(%v) got error %v", l, err)
	} else {
		got := out.String()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Write(%v) had diff with expected:\n%s", l, diff)
		}
	}
}

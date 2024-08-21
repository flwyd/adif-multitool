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

func TestEmptyCSV(t *testing.T) {
	input := "QSO_DATE,TIME_ON,BAND,CALLSIGN,NAME\n"
	csv := NewCSVIO()
	if parsed, err := csv.Read(strings.NewReader(input)); err != nil {
		t.Errorf("Read(%q) got error %v", input, err)
	} else {
		if got := len(parsed.Records); got != 0 {
			t.Errorf("Read(%q) got %d records, want %d", input, got, 0)
		}
	}
}

func TestReadCSV(t *testing.T) {
	// first record doesn't have trraiiling comma for notes
	input := `QSO_DATE,TIME_ON,BAND,CALLSIGN,NAME,FREQ,NOTES_INTL
19901031,1234,40M,W1AW,Hiram Percy Maxim,7.054
20221224,095846,1.25cm,N0P,Santa Claus,,Þhrough the /❄\ bringing 🎁 \to the child\re\n
`
	wantFields := [][]Field{
		{
			{Name: "QSO_DATE", Value: "19901031"},
			{Name: "TIME_ON", Value: "1234"},
			{Name: "BAND", Value: "40M"},
			{Name: "CALLSIGN", Value: "W1AW"},
			{Name: "NAME", Value: "Hiram Percy Maxim"},
			{Name: "FREQ", Value: "7.054"},
			{Name: "NOTES_INTL", Value: ""},
		},
		{
			{Name: "QSO_DATE", Value: "20221224"},
			{Name: "TIME_ON", Value: "095846"},
			{Name: "BAND", Value: "1.25cm"},
			{Name: "CALLSIGN", Value: "N0P"},
			{Name: "NAME", Value: "Santa Claus"},
			{Name: "FREQ", Value: ""},
			{Name: "NOTES_INTL", Value: `Þhrough the /❄\ bringing 🎁 \to the child\re\n`},
		},
	}
	csv := NewCSVIO()
	if parsed, err := csv.Read(strings.NewReader(input)); err != nil {
		t.Errorf("Read(%q) got error %v", input, err)
	} else {
		for i, r := range parsed.Records {
			fields := r.Fields()
			if diff := cmp.Diff(wantFields[i], fields); diff != "" {
				t.Errorf("Read(%q) record %d did not match expected, diff:\n%s", input, i+1, diff)
			}
		}
		if gotlen := len(parsed.Records); gotlen != len(wantFields) {
			t.Errorf("Read(%q) got %d records:\n%v\nwant %d\n%v", input, gotlen, parsed.Records[len(wantFields):], len(wantFields), wantFields)
		}
	}
}

func TestWriteCSV(t *testing.T) {
	l := NewLogfile()
	l.Comment = "CSV ignores comments"
	l.AddRecord(NewRecord(
		Field{Name: "QSO_DATE", Value: "19901031", Type: TypeDate},
		Field{Name: "TIME_ON", Value: "1234", Type: TypeTime},
		Field{Name: "BAND", Value: "40M"},
		Field{Name: "CALLSIGN", Value: "W1AW"},
		Field{Name: "NAME", Value: "Hiram Percy Maxim", Type: TypeString},
		Field{Name: "FREQ", Value: "7.054"},
	)).AddRecord(NewRecord(
		Field{Name: "QSO_DATE", Value: "20221224"},
		Field{Name: "TIME_ON", Value: "095846"},
		Field{Name: "BAND", Value: "1.25cm", Type: TypeEnumeration},
		Field{Name: "CALLSIGN", Value: "N0P", Type: TypeString},
		Field{Name: "NAME", Value: "Santa Claus"},
		Field{Name: "NOTES_INTL", Value: `Þhrough the /❄\ bringing 🎁 \to the child\re\n`},
	))
	l.Records[1].SetComment("Record comment")
	want := `QSO_DATE,TIME_ON,BAND,CALLSIGN,NAME,FREQ,NOTES_INTL
19901031,1234,40M,W1AW,Hiram Percy Maxim,7.054,
20221224,095846,1.25cm,N0P,Santa Claus,,Þhrough the /❄\ bringing 🎁 \to the child\re\n
`
	csv := NewCSVIO()
	out := &strings.Builder{}
	if err := csv.Write(l, out); err != nil {
		t.Errorf("Write(%v) got error %v", l, err)
	} else {
		got := out.String()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Write(%v) had diff with expected:\n%s", l, diff)
		}
	}
}

func TestCSVQuotes(t *testing.T) {
	input := `QSO_DATE,NAME,VUCC_GRIDS,ADDRESS_INTL,NOTES
19990101,"""C.G."" Tuska","FN31,FN21","225 Main St.
Newington, CT 06111",Unquoted note field
19990202,,AA00,"000 ""The South Pole"", Antarctica","
 leading newline and trailing newline
"
19990303,"Maxim, Hiram Percy,","DN00,DN01,DM90,DM91","Pl. des Nations 1211
1202 Genève
Switzerland",",comma notes,,,"
`
	wantFields := [][]Field{
		{
			{Name: "QSO_DATE", Value: "19990101"},
			{Name: "NAME", Value: `"C.G." Tuska`},
			{Name: "VUCC_GRIDS", Value: "FN31,FN21"},
			{Name: "ADDRESS_INTL", Value: "225 Main St.\nNewington, CT 06111"},
			{Name: "NOTES", Value: "Unquoted note field"},
		},
		{
			{Name: "QSO_DATE", Value: "19990202"},
			{Name: "NAME", Value: ""},
			{Name: "VUCC_GRIDS", Value: "AA00"},
			{Name: "ADDRESS_INTL", Value: `000 "The South Pole", Antarctica`},
			{Name: "NOTES", Value: "\n leading newline and trailing newline\n"},
		},
		{
			{Name: "QSO_DATE", Value: "19990303"},
			{Name: "NAME", Value: "Maxim, Hiram Percy,"},
			{Name: "VUCC_GRIDS", Value: "DN00,DN01,DM90,DM91"},
			{Name: "ADDRESS_INTL", Value: "Pl. des Nations 1211\n1202 Genève\nSwitzerland"},
			{Name: "NOTES", Value: ",comma notes,,,"},
		},
	}
	csv := NewCSVIO()
	if parsed, err := csv.Read(strings.NewReader(input)); err != nil {
		t.Errorf("Read(%q) got error %v", input, err)
	} else {
		for i, r := range parsed.Records {
			fields := r.Fields()
			if diff := cmp.Diff(wantFields[i], fields); diff != "" {
				t.Errorf("Read(%q) record %d did not match expected, diff:\n%s", input, i+1, diff)
			}
		}
		if gotlen := len(parsed.Records); gotlen != len(wantFields) {
			t.Errorf("Read(%q) got %d records:\n%v\nwant %d\n%v", input, gotlen, parsed.Records[len(wantFields):], len(wantFields), wantFields)
		}
		out := &strings.Builder{}
		if err := csv.Write(parsed, out); err != nil {
			t.Errorf("Write after Read(%q) got error %v", input, err)
		} else if diff := cmp.Diff(input, out.String()); diff != "" {
			t.Errorf("CSV round trip got diff:\n%s", diff)
		}
	}
}

func TestCSVOmitHeader(t *testing.T) {
	l := NewLogfile()
	l.Comment = "CSV ignores comments"
	l.AddRecord(NewRecord(
		Field{Name: "QSO_DATE", Value: "19901031", Type: TypeDate},
		Field{Name: "TIME_ON", Value: "1234", Type: TypeTime},
		Field{Name: "BAND", Value: "40M"},
		Field{Name: "CALLSIGN", Value: "W1AW"},
		Field{Name: "NAME", Value: "Hiram Percy Maxim", Type: TypeString},
		Field{Name: "FREQ", Value: "7.054"},
	)).AddRecord(NewRecord(
		Field{Name: "QSO_DATE", Value: "20221224"},
		Field{Name: "TIME_ON", Value: "095846"},
		Field{Name: "BAND", Value: "1.25cm", Type: TypeEnumeration},
		Field{Name: "CALLSIGN", Value: "N0P", Type: TypeString},
		Field{Name: "NAME", Value: "Santa Claus"},
		Field{Name: "NOTES_INTL", Value: `Þhrough the /❄\ bringing 🎁 \to the child\re\n`},
	))
	l.Records[1].SetComment("Record comment")
	want := `19901031,1234,40M,W1AW,Hiram Percy Maxim,7.054,
20221224,095846,1.25cm,N0P,Santa Claus,,Þhrough the /❄\ bringing 🎁 \to the child\re\n
`
	csv := NewCSVIO()
	csv.OmitHeader = true
	out := &strings.Builder{}
	if err := csv.Write(l, out); err != nil {
		t.Errorf("Write(%v) got error %v", l, err)
	} else {
		got := out.String()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Write(%v) had diff with expected:\n%s", l, diff)
		}
	}
}

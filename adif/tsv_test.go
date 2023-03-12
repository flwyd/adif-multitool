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

func TestEmptyTSV(t *testing.T) {
	input := "QSO_DATE\tTIME_ON\tBAND\tCALLSIGN\tNAME\n"
	tsv := NewTSVIO()
	if parsed, err := tsv.Read(strings.NewReader(input)); err != nil {
		t.Errorf("Read(%q) got error %v", input, err)
	} else {
		if got := len(parsed.Records); got != 0 {
			t.Errorf("Read(%q) got %d records, want %d", input, got, 0)
		}
	}
}

func TestReadTSV(t *testing.T) {
	input := "QSO_DATE\tTIME_ON\tBAND\tCALL\tNAME\tFREQ\tNOTES_INTL\n" +
		"19901031\t1234\t40M\tW1AW\tHiram Percy Maxim\t7.054\t\n" +
		"20221224\t095846\t1.25cm\tN0P\tSanta \"St. Nick\" Claus\t\tÞhrough the /❄\\ bringing 🎁 \\to the child\\re\\n\n"
	wantFields := [][]Field{
		{
			{Name: "QSO_DATE", Value: "19901031"},
			{Name: "TIME_ON", Value: "1234"},
			{Name: "BAND", Value: "40M"},
			{Name: "CALL", Value: "W1AW"},
			{Name: "NAME", Value: "Hiram Percy Maxim"},
			{Name: "FREQ", Value: "7.054"},
			{Name: "NOTES_INTL", Value: ""},
		}, {
			{Name: "QSO_DATE", Value: "20221224"},
			{Name: "TIME_ON", Value: "095846"},
			{Name: "BAND", Value: "1.25cm"},
			{Name: "CALL", Value: "N0P"},
			{Name: "NAME", Value: "Santa \"St. Nick\" Claus"},
			{Name: "FREQ", Value: ""},
			{Name: "NOTES_INTL", Value: "Þhrough the /❄\\ bringing 🎁 \\to the child\\re\\n"},
		},
	}
	tsv := NewTSVIO()
	if parsed, err := tsv.Read(strings.NewReader(input)); err != nil {
		t.Errorf("Read(%q) got error %v", input, err)
	} else {
		for i, r := range parsed.Records {
			fields := r.Fields()
			if diff := cmp.Diff(wantFields[i], fields); diff != "" {
				t.Errorf("Read(%q) record %d did not match expected, diff:\n%s", input, i, diff)
			}
		}
		if gotlen := len(parsed.Records); gotlen != len(wantFields) {
			t.Errorf("Read(%q) got %d records:\n%v\nwant %d\n%v", input, gotlen, parsed.Records[len(wantFields):], len(wantFields), wantFields)
		}
	}
}

func TestWriteTSV(t *testing.T) {
	l := NewLogfile()
	l.Comment = "TSV ignores comments"
	l.FieldOrder = []string{"QSO_DATE", "TIME_ON", "BAND", "CALL"}
	l.AddRecord(NewRecord(
		Field{Name: "TIME_ON", Value: "1234", Type: TypeTime},
		Field{Name: "QSO_DATE", Value: "19901031", Type: TypeDate},
		Field{Name: "NAME", Value: "Hiram Percy Maxim", Type: TypeString},
		Field{Name: "BAND", Value: "40M"},
		Field{Name: "FREQ", Value: "7.054"},
		Field{Name: "CALL", Value: "W1AW"},
	)).AddRecord(NewRecord(
		Field{Name: "notes_intl", Value: "Þhrough the /❄\\ bringing 🎁 \\to the child\\re\\n"},
		Field{Name: "qso_date", Value: "20221224"},
		Field{Name: "call", Value: "N0P", Type: TypeString},
		Field{Name: "time_on", Value: "095846"},
		Field{Name: "band", Value: "1.25cm", Type: TypeEnumeration},
		Field{Name: "name", Value: "Santa \"St. Nick\" Claus"},
	))
	l.Records[1].SetComment("Record comment")
	want := "QSO_DATE\tTIME_ON\tBAND\tCALL\tNAME\tFREQ\tNOTES_INTL\n" +
		"19901031\t1234\t40M\tW1AW\tHiram Percy Maxim\t7.054\t\n" +
		"20221224\t095846\t1.25cm\tN0P\tSanta \"St. Nick\" Claus\t\tÞhrough the /❄\\ bringing 🎁 \\to the child\\re\\n\n"
	tsv := NewTSVIO()
	out := &strings.Builder{}
	if err := tsv.Write(l, out); err != nil {
		t.Errorf("Write(%v) got error %v", l, err)
	} else {
		got := out.String()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Write(%v) had diff with expected:\n%s", l, diff)
		}
	}
}

func TestTSVRejectSpecialCharacters(t *testing.T) {
	tests := []*Logfile{
		NewLogfile().AddRecord(NewRecord(Field{Name: "APP_TEST_TrailingNewline", Value: "Foo\n"})),
		NewLogfile().AddRecord(NewRecord(Field{Name: "FREQ", Value: "146.52"}, Field{Name: "TAB\tHEADER", Value: "Foo"})),
		NewLogfile().AddRecord(NewRecord(Field{Name: "CALL", Value: "W1AW"}, Field{Name: "Notes", Value: "\r"})),
		NewLogfile().AddRecord(NewRecord(Field{Name: "CALL", Value: "W1AW"}, Field{Name: "ADDRESS", Value: "225 Main St\r\nNewington, CT 06111"})),
		NewLogfile().AddRecord(NewRecord(Field{Name: "RIG", Value: "\tHeathkit HW-100"}, Field{Name: "FREQ", Value: "14.300"})),
		NewLogfile().AddRecord(NewRecord(Field{Name: "CALL", Value: "W1AW"}, Field{Name: "BAND", Value: "40m"})).
			AddRecord(NewRecord(Field{Name: "CALL", Value: "N0P"}, Field{Name: "NAME", Value: "Santa\nClaus"}, Field{Name: "BAND", Value: "2m"})),
	}
	tabInOrder := NewLogfile().AddRecord(NewRecord(Field{Name: "CALL", Value: "W1AW"}, Field{Name: "NOTTES", Value: "Record isi fine"}))
	tabInOrder.FieldOrder = []string{"CALL", "NOTES", "BAD\tFIELD"}
	for _, l := range tests {
		tsv := NewTSVIO()
		out := &strings.Builder{}
		if err := tsv.Write(l, out); err == nil {
			t.Errorf("Write(%v) got %q, want error", l, out)
		}
	}
}

func TestTSVEscapeSpecialCharacters(t *testing.T) {
	tests := []struct {
		log  *Logfile
		want string
	}{
		{
			log:  NewLogfile().AddRecord(NewRecord(Field{Name: "APP_TEST_TrailingNewline", Value: "Foo\n"})),
			want: "APP_TEST_TRAILINGNEWLINE\nFoo\\n\n",
		},
		{
			log:  NewLogfile().AddRecord(NewRecord(Field{Name: "FREQ", Value: "146.52"}, Field{Name: "TAB\tHEADER", Value: "Foo"})),
			want: "FREQ\tTAB\\tHEADER\n146.52\tFoo\n",
		},
		{
			log:  NewLogfile().AddRecord(NewRecord(Field{Name: "CALL", Value: "W1AW"}, Field{Name: "Notes", Value: "\r"})),
			want: "CALL\tNOTES\nW1AW\t\\r\n",
		},
		{
			log:  NewLogfile().AddRecord(NewRecord(Field{Name: "CALL", Value: "W1AW"}, Field{Name: "ADDRESS", Value: "225 Main St\r\nNewington, CT 06111"})),
			want: "CALL\tADDRESS\nW1AW\t225 Main St\\r\\nNewington, CT 06111\n",
		},
		{
			log:  NewLogfile().AddRecord(NewRecord(Field{Name: "RIG", Value: "\tHeathkit HW-100"}, Field{Name: "FREQ", Value: "14.300"})),
			want: "RIG\tFREQ\n\\tHeathkit HW-100\t14.300\n",
		},
		{
			log:  NewLogfile().AddRecord(NewRecord(Field{Name: "NOTES_INTL", Value: "好话\r\nVery \"big\" signal\r\n"}, Field{Name: "DXCC", Value: "318"})),
			want: "NOTES_INTL\tDXCC\n好话\\r\\nVery \"big\" signal\\r\\n\t318\n",
		},
		{
			log: NewLogfile().AddRecord(NewRecord(Field{Name: "CALL", Value: "W1AW"}, Field{Name: "BAND", Value: "40m"})).
				AddRecord(NewRecord(Field{Name: "CALL", Value: "N0P"}, Field{Name: "NAME", Value: "Santa\nClaus"}, Field{Name: "BAND", Value: "2m"})),
			want: "CALL\tBAND\tNAME\nW1AW\t40m\t\nN0P\t2m\tSanta\\nClaus\n",
		},
	}
	tabInOrder := NewLogfile().AddRecord(NewRecord(Field{Name: "CALL", Value: "W1AW"}, Field{Name: "NOTTES", Value: "Record isi fine"}))
	tabInOrder.FieldOrder = []string{"CALL", "NOTES", "BAD\tFIELD"}
	for _, tc := range tests {
		tsv := NewTSVIO()
		tsv.EscapeSpecial = true
		out := &strings.Builder{}
		if err := tsv.Write(tc.log, out); err != nil {
			t.Errorf("Write(%v) got error %v", tc.log, err)
		} else if diff := cmp.Diff(tc.want, out.String()); diff != "" {
			t.Errorf("Write(%v) got diff:\n%s", tc.log, diff)
		} else {
			in := strings.NewReader(out.String())
			if parsed, err := tsv.Read(in); err != nil {
				t.Errorf("Read(%q) got error %v", out, err)
			} else if diff := cmp.Diff(tc.log.Records, parsed.Records); diff != "" {
				t.Errorf("Read(%q) does not match original logfile, diff:\n%s", out, diff)
			}
		}
	}
}

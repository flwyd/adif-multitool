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

package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/flwyd/adif-multitool/adif"
	"github.com/google/go-cmp/cmp"
)

func TestFixEmpty(t *testing.T) {
	adi := adif.NewADIIO()
	csv := adif.NewCSVIO()
	out := &bytes.Buffer{}
	file1 := "QSO_DATE,TIME_ON,TIME_OFF,CALL,MODE,BAND\n"
	ctx := &Context{
		OutputFormat: adif.FormatADI,
		Readers:      readers(adi, csv),
		Writers:      writers(adi, csv),
		Out:          out,
		Prepare:      testPrepare("My Comment", "3.1.4", "fix test", "1.2.3"),
		fs:           fakeFilesystem{map[string]string{"foo.csv": file1}},
		CommandCtx: &EditContext{
			Add:    FieldAssignments{values: []adif.Field{{Name: "BAZ", Value: "Baz value"}}, validate: ValidateAlphanumName},
			Set:    FieldAssignments{values: []adif.Field{{Name: "FOO", Value: "Foo value"}}, validate: ValidateAlphanumName},
			Remove: []string{"BAR"},
		}}
	if err := Fix.Run(ctx, []string{"foo.csv"}); err != nil {
		t.Errorf("Fix.Run(ctx, foo.csv) got error %v", err)
	} else {
		got := out.String()
		want := "My Comment\n<ADIF_VER:5>3.1.4 <PROGRAMID:8>fix test <PROGRAMVERSION:5>1.2.3 <EOH>\n"
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Fix.Run(ctx, foo.csv) unexpected output, diff:\n%s", diff)
		}
	}
}

func TestFixDate(t *testing.T) {
	adi := adif.NewADIIO()
	csv := adif.NewCSVIO()
	header := "My Comment\n<ADIF_VER:5>3.1.4 <PROGRAMID:8>fix test <PROGRAMVERSION:5>1.2.3 <EOH>\n"
	fields := []string{
		"CLUBLOG_QSO_UPLOAD_DATE", "EQSL_QSLRDATE", "EQSL_QSLSDATE",
		"HAMLOGEU_QSO_UPLOAD_DATE", "HAMQTH_QSO_UPLOAD_DATE", "HRDLOG_QSO_UPLOAD_DATE",
		"LOTW_QSLRDATE", "LOTW_QSLSDATE", "QRZCOM_QSO_UPLOAD_DATE",
		"QSLRDATE", "QSLSDATE", "QSO_DATE", "QSO_DATE_OFF",
	}
	tests := []struct{ source, want string }{
		{source: "19870615", want: "19870615"},
		{source: "2009-08-07", want: "20090807"},
		{source: "1925.11.10", want: "19251110"},
		{source: "2001/2/3", want: "20010203"},
		{source: "1/2/2003", want: "1/2/2003"},       // don't know if m/d or d/m
		{source: "2013:04:05", want: "2013:04:05"},   // unknown delimiter
		{source: "10002-09-08", want: "10002-09-08"}, // too many digits
		{source: "100020908", want: "100020908"},     // too many digits
		{source: "12345", want: "12345"},             // too few digits
		{source: "123", want: "123"},                 // too few digits
		{source: "2012-13-05", want: "2012-13-05"},   // unknown month
		{source: "2001年5月30日", want: "2001年5月30日"},   // Chinese date
	}
	for _, tc := range tests {
		for _, f := range fields {
			out := &bytes.Buffer{}
			file1 := fmt.Sprintf("NOT_DATE,%s\n2022-02-22,%s\n", f, tc.source)
			ctx := &Context{
				OutputFormat: adif.FormatADI,
				Readers:      readers(adi, csv),
				Writers:      writers(adi, csv),
				Out:          out,
				Prepare:      testPrepare("My Comment", "3.1.4", "fix test", "1.2.3"),
				fs:           fakeFilesystem{map[string]string{"foo.csv": file1}}}
			if err := Fix.Run(ctx, []string{"foo.csv"}); err != nil {
				t.Errorf("Fix.Run(ctx, foo.csv) got error %v", err)
			} else {
				got := out.String()
				want := fmt.Sprintf("%s<NOT_DATE:10>2022-02-22 <%s:%d>%s <EOR>\n", header, f, len(tc.want), tc.want)
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("fix %s=%s want %s got diff %s", f, tc.source, tc.want, diff)
				}
			}
		}
	}
}

func TestFixTime(t *testing.T) {
	adi := adif.NewADIIO()
	csv := adif.NewCSVIO()
	header := "My Comment\n<ADIF_VER:5>3.1.4 <PROGRAMID:8>fix test <PROGRAMVERSION:5>1.2.3 <EOH>\n"
	fields := []string{"TIME_ON", "TIME_OFF"}
	tests := []struct{ source, want string }{
		{source: "010203", want: "010203"},
		{source: "0102", want: "0102"},
		{source: "10203", want: "010203"},
		{source: "102", want: "0102"},
		{source: "131415", want: "131415"},
		{source: "1314", want: "1314"},
		{source: "14:15:16", want: "141516"},
		{source: "14:15", want: "1415"},
		{source: "4:56:30 pm", want: "165630"},
		{source: "4:56 pm", want: "1656"},
		{source: "4:56:30 PM", want: "165630"},
		{source: "4:56 PM", want: "1656"},
		{source: "4:56:30 AM", want: "045630"},
		{source: "4:56 AM", want: "0456"},
		{source: "4:56:30pm", want: "165630"},
		{source: "4:56pm", want: "1656"},
		{source: "4:56:30PM", want: "165630"},
		{source: "4:56PM", want: "1656"},
		{source: "4:56:30AM", want: "045630"},
		{source: "4:56AM", want: "0456"},
		{source: "2009-08-07", want: "2009-08-07"}, // not a time
		{source: "89:12:31", want: "89:12:31"},     // invalid hour
		{source: "12/34/56", want: "12/34/56"},     // unknown delimiter
		{source: "12/34", want: "12/34"},           // unknown delimiter
		{source: "1234567", want: "1234567"},       // too many digits
		{source: "11", want: "11"},                 // too few digits
	}
	for _, tc := range tests {
		for _, f := range fields {
			out := &bytes.Buffer{}
			file1 := fmt.Sprintf("NOT_TIME,%s\n07:06:05,%s\n", f, tc.source)
			ctx := &Context{
				OutputFormat: adif.FormatADI,
				Readers:      readers(adi, csv),
				Writers:      writers(adi, csv),
				Out:          out,
				Prepare:      testPrepare("My Comment", "3.1.4", "fix test", "1.2.3"),
				fs:           fakeFilesystem{map[string]string{"foo.csv": file1}}}
			if err := Fix.Run(ctx, []string{"foo.csv"}); err != nil {
				t.Errorf("Fix.Run(ctx, foo.csv) got error %v", err)
			} else {
				got := out.String()
				want := fmt.Sprintf("%s<NOT_TIME:8>07:06:05 <%s:%d>%s <EOR>\n", header, f, len(tc.want), tc.want)
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("fix %s=%s want %s got diff %s", f, tc.source, tc.want, diff)
				}
			}
		}
	}
}

func TestFixLocation(t *testing.T) {
	adi := adif.NewADIIO()
	csv := adif.NewCSVIO()
	header := "My Comment\n<ADIF_VER:5>3.1.4 <PROGRAMID:8>fix test <PROGRAMVERSION:5>1.2.3 <USERDEF1:13:L>LONG_LATITUDE <USERDEF2:14:L>LATE_LONGITUDE <EOH>\n"
	tests := []struct{ name, value, want string }{
		{name: "LAT", value: "N012 34.567", want: "N012 34.567"},
		{name: "LON", value: "W123 45.678", want: "W123 45.678"},
		{name: "MY_LAT", value: "S000 59.987", want: "S000 59.987"},
		{name: "MY_LON", value: "E009 09.000", want: "E009 09.000"},
		{name: "LONG_LATITUDE", value: "N089 59.999", want: "N089 59.999"},
		{name: "LATE_LONGITUDE", value: "E179 00.000", want: "E179 00.000"},
		// to check deg mm.mmm values see https://pgc.umn.edu/apps/convert/
		{name: "LAT", value: "12.50123", want: "N012 30.074"},
		{name: "LON", value: "-123.25987654", want: "W123 15.593"},
		{name: "MY_LAT", value: "-6.75", want: "S006 45.000"},
		{name: "MY_LON", value: "0.003", want: "E000 00.180"},
		{name: "LONG_LATITUDE", value: "9.99999", want: "N009 59.999"},
		{name: "LATE_LONGITUDE", value: "+179.0", want: "E179 00.000"},
		// leave out-of-range values alone
		{name: "LAT", value: "90.0001", want: "90.0001"},
		{name: "LON", value: "181.123", want: "181.123"},
		{name: "MY_LAT", value: "-123.456", want: "-123.456"},
		{name: "MY_LON", value: "-200.1234", want: "-200.1234"},
	}
	for _, tc := range tests {
		out := &bytes.Buffer{}
		file1 := fmt.Sprintf("LONG_WAVE,%s\n83.54294,%s\n", tc.name, tc.value)
		ctx := &Context{
			OutputFormat: adif.FormatADI,
			Readers:      readers(adi, csv),
			Writers:      writers(adi, csv),
			Out:          out,
			UserdefFields: UserdefFieldList{adif.UserdefField{Name: "LONG_LATITUDE", Type: adif.TypeLocation},
				adif.UserdefField{Name: "LATE_LONGITUDE", Type: adif.TypeLocation}},
			Prepare: testPrepare("My Comment", "3.1.4", "fix test", "1.2.3"),
			fs:      fakeFilesystem{map[string]string{"foo.csv": file1}}}
		if err := Fix.Run(ctx, []string{"foo.csv"}); err != nil {
			t.Errorf("Fix.Run(ctx, foo.csv) got error %v", err)
		} else {
			got := out.String()
			want := fmt.Sprintf("%s<LONG_WAVE:8>83.54294 <%s:%d>%s <EOR>\n", header, tc.name, len(tc.want), tc.want)
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("fix %s=%s want %s got diff %s", tc.name, tc.value, tc.want, diff)
			}
		}
	}
}

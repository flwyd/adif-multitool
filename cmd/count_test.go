// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/flwyd/adif-multitool/adif"
	"github.com/google/go-cmp/cmp"
)

func TestCountEmpty(t *testing.T) {
	adi := adif.NewADIIO()
	csv := adif.NewCSVIO()
	out := &bytes.Buffer{}
	file1 := "FOO,BAR\n"
	ctx := &Context{
		OutputFormat: adif.FormatCSV,
		Readers:      readers(adi, csv),
		Writers:      writers(adi, csv),
		Out:          out,
		Prepare:      testPrepare("My Comment", "3.1.5", "count test", "1.2.3"),
		fs:           fakeFilesystem{map[string]string{"foo.csv": file1}},
		CommandCtx:   &CountContext{CountFieldName: "NUM", Fields: FieldList{"BAZ", "FOO"}}}
	if err := Count.Run(ctx, []string{"foo.csv"}); err != nil {
		t.Errorf("Count.Run(ctx, foo.csv) got error %v", err)
	} else {
		got := out.String()
		want := "NUM,BAZ,FOO\n0,,\n"
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Count.Run(ctx, foo.csv) unexpected output, diff:\n%s", diff)
		}
	}
}

func TestCountEmptyFieldList(t *testing.T) {
	adi := adif.NewADIIO()
	csv := adif.NewCSVIO()
	out := &bytes.Buffer{}
	file1 := `QSO_DATE,TIME_ON,TIME_OFF,CALL,BAND,FREQ,STATE,DXCC,TX_PWR,MODE,SUBMODE
20190202,1402,1405,N2B,80m,3.502,NY,291,20,CW,
20190101,1301,,K1A,160m,1.810,ME,291,10,SSB,LSB
20190303,1503,,W3C,40m,7.203,PA,291,30,SSB,LSB
`
	ctx := &Context{
		OutputFormat: adif.FormatCSV,
		Readers:      readers(adi, csv),
		Writers:      writers(adi, csv),
		Out:          out,
		Prepare:      testPrepare("My Comment", "3.1.5", "count test", "1.2.3"),
		fs:           fakeFilesystem{map[string]string{"foo.csv": file1}},
		CommandCtx:   &CountContext{}}
	if err := Count.Run(ctx, []string{"foo.csv"}); err != nil {
		t.Errorf("Count.Run(ctx, foo.csv) got error %v", err)
	} else {
		got := out.String()
		want := "COUNT\n3\n"
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Count.Run(ctx, foo.csv) unexpected output, diff:\n%s", diff)
		}
	}
}

func TestCountFields(t *testing.T) {
	csv := adif.NewCSVIO()
	file := `QSO_DATE,TIME_ON,CALL,BAND,FREQ,MODE,SUBMODE,DXCC,STATE
20010101,0123,AA1A,40m,7.200,SSB,LSB,291,AR
20010101,0234,VE2B,20M,14.012,CW,,1,QC
20020202,0345,K3C,20m,14.321,SSB,USB,291,DE
20020202,0346,K3C,2m,146.520,FM,,291,DE
20030303,0456,N4D,40m,7.2,SSB,LSB,291,DE
20040404,0506,I5E,20m,14.012,CW,,248,AR
20050505,0607,VE6F,40m,7.020,CW,,1,AB
20050505,0708,SG7G,15m,21.420,SSB,USB,284,AB
`
	tests := []struct {
		name   string
		want   []string
		fields FieldList
	}{
		{
			name:   "single string",
			fields: FieldList{"state"},
			want:   []string{"NUM,STATE", "2,AB", "2,AR", "3,DE", "1,QC"},
		},
		{
			name:   "two strings",
			fields: FieldList{"dxcc", "state"},
			want:   []string{"NUM,DXCC,STATE", "1,1,AB", "1,1,QC", "1,248,AR", "1,284,AB", "1,291,AR", "3,291,DE"},
		},
		{
			name:   "two strings, sort by second",
			fields: FieldList{"state", "dxcc"},
			want:   []string{"NUM,STATE,DXCC", "1,AB,1", "1,AB,284", "1,AR,248", "1,AR,291", "3,DE,291", "1,QC,1"},
		},
		{
			name:   "band sort",
			fields: FieldList{"band", "mode"},
			want:   []string{"NUM,BAND,MODE", "1,40m,CW", "2,40m,SSB", "2,20m,CW", "1,20m,SSB", "1,15m,SSB", "1,2m,FM"},
		},
		{
			name:   "number sort",
			fields: FieldList{"freq"},
			want:   []string{"NUM,FREQ", "1,7.020", "2,7.200", "2,14.012", "1,14.321", "1,21.420", "1,146.520"},
		},
		{
			name:   "some missing",
			fields: FieldList{"submode"},
			want:   []string{"NUM,SUBMODE", "4,", "2,LSB", "2,USB"},
		},
		{
			name:   "absent field",
			fields: FieldList{"name"},
			want:   []string{"NUM,NAME", "8,"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			ctx := &Context{
				OutputFormat: adif.FormatCSV,
				Readers:      readers(csv),
				Writers:      writers(csv),
				Out:          out,
				Prepare:      testPrepare("My Comment", "3.1.5", "count test", "1.2.3"),
				fs:           fakeFilesystem{map[string]string{"foo.csv": file}},
				CommandCtx:   &CountContext{Fields: tc.fields, CountFieldName: "num"},
			}
			if err := Count.Run(ctx, []string{"foo.csv"}); err != nil {
				t.Fatalf("Count.Run(ctx, foo.csv) got error %v", err)
			}
			if diff := cmp.Diff(strings.Join(tc.want, "\n")+"\n", out.String()); diff != "" {
				t.Errorf("%s %v got diff\n%s", tc.name, tc.fields, diff)
			}
		})
	}
}

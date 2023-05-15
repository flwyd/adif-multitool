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
	"testing"

	"github.com/flwyd/adif-multitool/adif"
	"github.com/google/go-cmp/cmp"
)

func TestSortEmpty(t *testing.T) {
	adi := adif.NewADIIO()
	csv := adif.NewCSVIO()
	out := &bytes.Buffer{}
	file1 := "FOO,BAR\n"
	ctx := &Context{
		OutputFormat: adif.FormatADI,
		Readers:      readers(adi, csv),
		Writers:      writers(adi, csv),
		Out:          out,
		Prepare:      testPrepare("My Comment", "3.1.4", "sort test", "1.2.3"),
		fs:           fakeFilesystem{map[string]string{"foo.csv": file1}},
		CommandCtx:   &SortContext{Fields: FieldList{"BAZ", "FOO"}}}
	if err := Sort.Run(ctx, []string{"foo.csv"}); err != nil {
		t.Errorf("Sort.Run(ctx, foo.csv) got error %v", err)
	} else {
		got := out.String()
		want := "My Comment\n<ADIF_VER:5>3.1.4 <PROGRAMID:9>sort test <PROGRAMVERSION:5>1.2.3 <EOH>\n"
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Sort.Run(ctx, foo.csv) unexpected output, diff:\n%s", diff)
		}
	}
}

func TestSortEmptyFieldList(t *testing.T) {
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
		Prepare:      testPrepare("My Comment", "3.1.4", "sort test", "1.2.3"),
		fs:           fakeFilesystem{map[string]string{"foo.csv": file1}},
		CommandCtx:   &SortContext{}}
	if err := Sort.Run(ctx, []string{"foo.csv"}); err != nil {
		t.Errorf("Sort.Run(ctx, foo.csv) got error %v", err)
	} else {
		got := out.String()
		if diff := cmp.Diff(file1, got); diff != "" {
			t.Errorf("Sort.Run(ctx, foo.csv) unexpected output, diff:\n%s", diff)
		}
	}
}

func TestSortCondition(t *testing.T) {
	csv := adif.NewCSVIO()
	tests := []struct {
		name, file, want string
		fields           FieldList
	}{
		{
			name:   "single string",
			fields: FieldList{"state"},
			file: `QSO_DATE,TIME_ON,CALL,STATE
20010101,0101,W1W,CT
20020202,0202,KL7K,AK
20030303,0303,A5A,TX
20040404,0404,N0N,MN
`,
			want: `QSO_DATE,TIME_ON,CALL,STATE
20020202,0202,KL7K,AK
20010101,0101,W1W,CT
20040404,0404,N0N,MN
20030303,0303,A5A,TX
`,
		},
		{
			name:   "single string reversed",
			fields: FieldList{"-state"},
			file: `QSO_DATE,TIME_ON,CALL,STATE
20010101,0101,W1W,CT
20020202,0202,KL7K,AK
20030303,0303,A5A,TX
20040404,0404,N0N,MN
`,
			want: `QSO_DATE,TIME_ON,CALL,STATE
20030303,0303,A5A,TX
20040404,0404,N0N,MN
20010101,0101,W1W,CT
20020202,0202,KL7K,AK
`,
		},
		{
			name:   "two strings",
			fields: FieldList{"state", "CALL"},
			file: `QSO_DATE,TIME_ON,CALL,STATE
20010101,0101,W1W,CT
20010101,0101,K1K,CT
20020202,0202,KL7K,AK
20030303,0303,N5N,TX
20030303,0303,A5A,TX
20040404,0404,K5K,NM
20040404,0404,N5M,NM
`,
			want: `QSO_DATE,TIME_ON,CALL,STATE
20020202,0202,KL7K,AK
20010101,0101,K1K,CT
20010101,0101,W1W,CT
20040404,0404,K5K,NM
20040404,0404,N5M,NM
20030303,0303,A5A,TX
20030303,0303,N5N,TX
`,
		},
		{
			name:   "strings ascending descending",
			fields: FieldList{"STATE", "-CALL"},
			file: `QSO_DATE,TIME_ON,CALL,STATE
20010101,0101,W1W,CT
20010101,0101,K1K,CT
20020202,0202,KL7K,AK
20030303,0303,N5N,TX
20030303,0303,A5A,TX
20040404,0404,K5K,NM
20040404,0404,N5M,NM
`,
			want: `QSO_DATE,TIME_ON,CALL,STATE
20020202,0202,KL7K,AK
20010101,0101,W1W,CT
20010101,0101,K1K,CT
20040404,0404,N5M,NM
20040404,0404,K5K,NM
20030303,0303,N5N,TX
20030303,0303,A5A,TX
`,
		},
		{
			name:   "one number, stable sort",
			fields: FieldList{"freq"},
			file: `QSO_DATE,TIME_ON,CALL,FREQ,STATE
20010101,0101,W1W,18.068,CT
20010101,0101,K1K,1.888,CT
20020202,0202,W4W,14.250,SC
20020202,0202,KL7K,7.123,AK
20030303,0303,N5N,28.28,TX
20030303,0303,A5A,14.250,TX
20040404,0404,K5K,14.025,NM
20040404,0404,N5M,7,NM
20040404,0404,W9W,28.28,MI
`,
			want: `QSO_DATE,TIME_ON,CALL,FREQ,STATE
20010101,0101,K1K,1.888,CT
20040404,0404,N5M,7,NM
20020202,0202,KL7K,7.123,AK
20040404,0404,K5K,14.025,NM
20020202,0202,W4W,14.250,SC
20030303,0303,A5A,14.250,TX
20010101,0101,W1W,18.068,CT
20030303,0303,N5N,28.28,TX
20040404,0404,W9W,28.28,MI
`,
		},
		{
			name:   "reverse number, stable sort",
			fields: FieldList{"-freq"},
			file: `QSO_DATE,TIME_ON,CALL,FREQ,STATE
20010101,0101,W1W,18.068,CT
20010101,0101,K1K,1.888,CT
20020202,0202,W4W,14.250,SC
20020202,0202,KL7K,7.123,AK
20030303,0303,N5N,28.28,TX
20030303,0303,A5A,14.250,TX
20040404,0404,K5K,14.025,NM
20040404,0404,N5M,7,NM
20040404,0404,W9W,28.28,MI
`,
			want: `QSO_DATE,TIME_ON,CALL,FREQ,STATE
20030303,0303,N5N,28.28,TX
20040404,0404,W9W,28.28,MI
20010101,0101,W1W,18.068,CT
20020202,0202,W4W,14.250,SC
20030303,0303,A5A,14.250,TX
20040404,0404,K5K,14.025,NM
20020202,0202,KL7K,7.123,AK
20040404,0404,N5M,7,NM
20010101,0101,K1K,1.888,CT
`,
		},
		{
			name:   "missing values",
			fields: FieldList{"time_on", "-state"},
			file: `QSO_DATE,TIME_ON,CALL,STATE
20010101,0101,W1W,CT
20010101,0101,K1K,
20020202,0202,KL7K,AK
20030303,,N5N,TX
20030303,0303,A5A,TX
20040404,0404,K5K,
20040404,0404,N5M,NM
`,
			want: `QSO_DATE,TIME_ON,CALL,STATE
20030303,,N5N,TX
20010101,0101,W1W,CT
20010101,0101,K1K,
20020202,0202,KL7K,AK
20030303,0303,A5A,TX
20040404,0404,N5M,NM
20040404,0404,K5K,
`,
		},
		{
			name:   "absent field",
			fields: FieldList{"SRX", "-qso_date", "state"},
			file: `QSO_DATE,TIME_ON,CALL,FREQ,STATE
20010101,0101,W1W,18.068,CT
20010101,0101,K1K,1.888,CT
20020202,0202,W4W,14.250,SC
20020202,0202,KL7K,7.123,AK
20030303,0303,N5N,28.28,TX
20030303,0303,A5A,14.250,TX
20040404,0404,K5K,14.025,NM
20040404,0404,N5M,7,NM
20040404,0404,W9W,28.28,MI
`,
			want: `QSO_DATE,TIME_ON,CALL,FREQ,STATE
20040404,0404,W9W,28.28,MI
20040404,0404,K5K,14.025,NM
20040404,0404,N5M,7,NM
20030303,0303,N5N,28.28,TX
20030303,0303,A5A,14.250,TX
20020202,0202,KL7K,7.123,AK
20020202,0202,W4W,14.250,SC
20010101,0101,W1W,18.068,CT
20010101,0101,K1K,1.888,CT
`,
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
				Prepare:      testPrepare("My Comment", "3.1.4", "sort test", "1.2.3"),
				fs:           fakeFilesystem{map[string]string{"foo.csv": tc.file}},
				CommandCtx:   &SortContext{Fields: tc.fields},
			}
			if err := Sort.Run(ctx, []string{"foo.csv"}); err != nil {
				t.Fatalf("Sort.Run(ctx, foo.csv) got error %v", err)
			}
			if diff := cmp.Diff(tc.want, out.String()); diff != "" {
				t.Errorf("%s %v got diff\n%s", tc.name, tc.fields, diff)
			}
		})
	}
}

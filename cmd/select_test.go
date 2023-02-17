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

const adiFile = `<CALL:4>W1AW <QSO_DATE:8>19870605 <NAME:17>Hiram Maxim Percy <FOO:16>Talked about Foo <EMPTY_FIELD:0> <EOR>
<CALL:3>N0P <QSO_DATE:8>20221224 <NAME:11>Santa Claus <EMPTY_FIELD:0> <BAND:3>80m <EOR>
`

func TestSelectNoFields(t *testing.T) {
	adi := adif.NewADIIO()
	csv := adif.NewCSVIO()
	out := &bytes.Buffer{}
	cctx := &SelectContext{Fields: make(FieldList, 0)}
	ctx := &Context{
		OutputFormat: adif.FormatCSV,
		Readers:      readers(adi, csv),
		Writers:      writers(adi, csv),
		Out:          out,
		Prepare:      testPrepare("My Comment", "3.1.4", "select test", "1.2.3"),
		fs:           fakeFilesystem{map[string]string{"foo.adi": adiFile}},
		CommandCtx:   cctx}
	if err := Select.Run(ctx, []string{"foo.adi"}); err == nil {
		t.Errorf("Select.Run(ctx) with no fields expected error but got\n%s", out.String())
	}
}

func TestSelectSingleFile(t *testing.T) {
	tests := []struct {
		fields []string
		want   string
	}{
		{fields: []string{"CALL"}, want: "CALL\nW1AW\nN0P\n"},
		{fields: []string{"foo"}, want: "FOO\nTalked about Foo\n"},
		{fields: []string{"band"}, want: "BAND\n80m\n"},
		{fields: []string{"EMPTY_FIELD"}, want: "EMPTY_FIELD\n\n\n"},
		{fields: []string{"NOT_PRESENT"}, want: "NOT_PRESENT\n"},
		{fields: []string{"NAME", "CALL"}, want: "NAME,CALL\nHiram Maxim Percy,W1AW\nSanta Claus,N0P\n"},
		{fields: []string{"band", "qso_date"}, want: "BAND,QSO_DATE\n,19870605\n80m,20221224\n"},
	}
	adi := adif.NewADIIO()
	csv := adif.NewCSVIO()
	for _, tc := range tests {
		out := &bytes.Buffer{}
		cctx := &SelectContext{Fields: tc.fields}
		ctx := &Context{
			OutputFormat: adif.FormatCSV,
			Readers:      readers(adi, csv),
			Writers:      writers(adi, csv),
			Out:          out,
			Prepare:      testPrepare("My Comment", "3.1.4", "select test", "1.2.3"),
			fs:           fakeFilesystem{map[string]string{"foo.adi": adiFile}},
			CommandCtx:   cctx}
		if err := Select.Run(ctx, []string{"foo.adi"}); err != nil {
			t.Errorf("Select.Run(ctx, foo.adi) got error %v", err)
		} else {
			got := out.String()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Select.Run(ctx, foo.adi) with fields %v got diff\n%s", tc.fields, diff)
			}
		}
	}
}

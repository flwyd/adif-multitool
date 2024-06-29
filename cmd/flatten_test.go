// Copyright 2024 Google LLC
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
	"testing"

	"github.com/flwyd/adif-multitool/adif"
	"github.com/google/go-cmp/cmp"
)

func TestFlattenEmpty(t *testing.T) {
	adi := adif.NewADIIO()
	csv := adif.NewCSVIO()
	out := &bytes.Buffer{}
	file1 := "VUCC_GRIDS,AWARD_SUBMITTED\n"
	ctx := &Context{
		OutputFormat: adif.FormatADI,
		Readers:      readers(adi, csv),
		Writers:      writers(adi, csv),
		Out:          out,
		Prepare:      testPrepare("My Comment", "3.1.4", "flatten test", "1.2.3"),
		fs:           fakeFilesystem{map[string]string{"foo.csv": file1}},
		CommandCtx:   &FlattenContext{Fields: FieldList{"AWARD_SUBMITTED", "VUCC_GRIDS"}}}
	if err := Flatten.Run(ctx, []string{"foo.csv"}); err != nil {
		t.Errorf("Flatten.Run(ctx, foo.csv) got error %v", err)
	} else {
		got := out.String()
		want := "My Comment\n<ADIF_VER:5>3.1.4 <PROGRAMID:12>flatten test <PROGRAMVERSION:5>1.2.3 <EOH>\n"
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Flatten.Run(ctx, foo.csv) unexpected output, diff:\n%s", diff)
		}
	}
}

func TestFlattenCartesian(t *testing.T) {
	tsv := adif.NewTSVIO()
	out := &bytes.Buffer{}
	file1 := `CALL	TIME_ON	POTA_REF	MY_USACA_COUNTIES
K1A	0101		AZ,Apache:CO,Montezuma:NM,San Juan:UT,San Juan
K2B	0202	US-1234	AZ,Apache:CO,Montezuma:NM,San Juan:UT,San Juan
K3C	0303	US-0655,US-4567	AZ,Apache:CO,Montezuma:NM,San Juan:UT,San Juan
`
	ctx := &Context{
		OutputFormat: adif.FormatTSV,
		Readers:      readers(tsv),
		Writers:      writers(tsv),
		Out:          out,
		Prepare:      testPrepare("My Comment", "3.1.4", "flatten test", "1.2.3"),
		fs:           fakeFilesystem{map[string]string{"foo.tsv": file1}},
		CommandCtx:   &FlattenContext{Fields: FieldList{"MY_USACA_COUNTIES", "POTA_REF"}}}
	if err := Flatten.Run(ctx, []string{"foo.tsv"}); err != nil {
		t.Errorf("Flatten.Run(ctx, foo.tsv) got error %v", err)
	} else {
		got := out.String()
		want := `CALL	TIME_ON	POTA_REF	MY_USACA_COUNTIES
K1A	0101		AZ,Apache
K1A	0101		CO,Montezuma
K1A	0101		NM,San Juan
K1A	0101		UT,San Juan
K2B	0202	US-1234	AZ,Apache
K2B	0202	US-1234	CO,Montezuma
K2B	0202	US-1234	NM,San Juan
K2B	0202	US-1234	UT,San Juan
K3C	0303	US-0655	AZ,Apache
K3C	0303	US-4567	AZ,Apache
K3C	0303	US-0655	CO,Montezuma
K3C	0303	US-4567	CO,Montezuma
K3C	0303	US-0655	NM,San Juan
K3C	0303	US-4567	NM,San Juan
K3C	0303	US-0655	UT,San Juan
K3C	0303	US-4567	UT,San Juan
`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Flatten.Run(ctx, foo.tsv) unexpected output, diff:\n%s", diff)
		}
	}
}

func TestFlattenCustomDelimiter(t *testing.T) {
	tsv := adif.NewTSVIO()
	out := &bytes.Buffer{}
	file1 := `CALL	TIME_ON	SRX_STRING	STX_STRING
K1A	0101	ABC	ZYX
K2B	0202	DEF/GHI	WVU
K3C	0303	JKL	TSR QPO
`
	ctx := &Context{
		OutputFormat: adif.FormatTSV,
		Readers:      readers(tsv),
		Writers:      writers(tsv),
		Out:          out,
		Prepare:      testPrepare("My Comment", "3.1.4", "flatten test", "1.2.3"),
		fs:           fakeFilesystem{map[string]string{"foo.tsv": file1}},
		CommandCtx: &FlattenContext{
			Fields: FieldList{"SRX_STRING", "STX_STRING"},
			Delimiters: FieldAssignments{values: []adif.Field{
				{Name: "SRX_STRING", Value: "/"},
				{Name: "STX_STRING", Value: " "}}},
		}}
	if err := Flatten.Run(ctx, []string{"foo.tsv"}); err != nil {
		t.Errorf("Flatten.Run(ctx, foo.tsv) got error %v", err)
	} else {
		got := out.String()
		want := `CALL	TIME_ON	SRX_STRING	STX_STRING
K1A	0101	ABC	ZYX
K2B	0202	DEF	WVU
K2B	0202	GHI	WVU
K3C	0303	JKL	TSR
K3C	0303	JKL	QPO
`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Flatten.Run(ctx, foo.tsv) unexpected output, diff:\n%s", diff)
		}
	}
}

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

func TestCatEmpty(t *testing.T) {
	io := adif.NewADIIO()
	out := &bytes.Buffer{}
	ctx := &Context{
		InputFormat:  adif.FormatADI,
		OutputFormat: adif.FormatADI,
		Readers:      readers(io),
		Writers:      writers(io),
		Out:          out,
		Prepare:      testPrepare("My Comment", "3.1.4", "cat test", "1.2.3"),
		fs:           fakeFilesystem{map[string]string{"-": ""}}}
	if err := Cat.Run(ctx, []string{}); err != nil {
		t.Errorf("Cat.Run(ctx) got error %v", err)
	} else {
		got := out.String()
		want := "My Comment\n<ADIF_VER:5>3.1.4 <PROGRAMID:8>cat test <PROGRAMVERSION:5>1.2.3 <EOH>\n"
		if got != want {
			t.Errorf("Cat.Run(ctx) got %s, want %s", got, want)
		}
	}
}

func TestEditADIToCSV(t *testing.T) {
	adi := adif.NewADIIO()
	csv := adif.NewCSVIO()
	out := &bytes.Buffer{}
	file1 := `Generated 2020-06-21
<ADIF_VER:5>3.1.4 <PROGRAMID:8>cat test <PROGRAMVERSION:5>1.2.3 <EOH>
<FIELD_1:4>Alfa <FOO:5>Bravo <FIELD_2:7>Charlie <EOR>
<FIELD_1:5>Delta <FOO:4>Echo <EOR>
`
	file2 := `Generated 1999-12-31
<ADIF_VER:5>2.1.9 <PROGRAMID:8>Fancy Software <PROGRAMVERSION:5>(devel) <EOH>

<BAR:4>Golf <FIELD_2:5>Hotel <FIELD_1:8>Fox Trot <EOR>

<FIELD_2:5>India <BAR:7>Juliett <today:8:D>19870605 <now:4:t>1234 <EOR>
`
	ctx := &Context{
		OutputFormat: adif.FormatCSV,
		Readers:      readers(adi, csv),
		Writers:      writers(adi, csv),
		Out:          out,
		Prepare:      testPrepare("My Comment", "3.1.4", "cat test", "1.2.3"),
		fs:           fakeFilesystem{map[string]string{"foo.adi": file1, "bar.adi": file2}}}
	if err := Cat.Run(ctx, []string{"foo.adi", "bar.adi"}); err != nil {
		t.Errorf("Cat.Run(ctx) got error %v", err)
	} else {
		got := out.String()
		want := `FIELD_1,FOO,FIELD_2,BAR,TODAY,NOW
Alfa,Bravo,Charlie,,,
Delta,Echo,,,,
Fox Trot,,Hotel,Golf,,
,,India,Juliett,19870605,1234
`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Cat.Run(ctx, foo.adi, bar.adi) unexpected output, diff:\n%s", diff)
		}
	}
}

func TestCatCSVToADI(t *testing.T) {
	adi := adif.NewADIIO()
	csv := adif.NewCSVIO()
	out := &bytes.Buffer{}
	file1 := `FIELD_1,FOO,FIELD_2

Alfa,Bravo,Charlie

Delta,Echo,
`
	file2 := `FIELD_1,BAR,FIELD_2,TODAY,NOW
Fox Trot,Golf,Hotel,,
,Juliett,India,19870605,1234
`
	ctx := &Context{
		OutputFormat: adif.FormatADI,
		Readers:      readers(adi, csv),
		Writers:      writers(adi, csv),
		Out:          out,
		Prepare:      testPrepare("My Comment", "3.1.4", "cat test", "1.2.3"),
		fs:           fakeFilesystem{map[string]string{"foo.csv": file1, "bar.csv": file2}}}
	if err := Cat.Run(ctx, []string{"foo.csv", "bar.csv"}); err != nil {
		t.Errorf("Cat.Run(ctx) got error %v", err)
	} else {
		got := out.String()
		want := `My Comment
<ADIF_VER:5>3.1.4 <PROGRAMID:8>cat test <PROGRAMVERSION:5>1.2.3 <EOH>
<FIELD_1:4>Alfa <FOO:5>Bravo <FIELD_2:7>Charlie <EOR>
<FIELD_1:5>Delta <FOO:4>Echo <FIELD_2:0> <EOR>
<FIELD_1:8>Fox Trot <FIELD_2:5>Hotel <BAR:4>Golf <TODAY:0> <NOW:0> <EOR>
<FIELD_1:0> <FIELD_2:5>India <BAR:7>Juliett <TODAY:8>19870605 <NOW:4>1234 <EOR>
`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Cat.Run(ctx, foo.ctx, bar.ctx) unexpected output, diff:\n%s", diff)
		}
	}
}

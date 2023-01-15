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

func TestEditEmpty(t *testing.T) {
	adi := adif.NewADIIO()
	adi.HeaderCommentFn = func(*adif.Logfile) string { return "My Comment" }
	csv := adif.NewCSVIO()
	out := &bytes.Buffer{}
	file1 := "FOO,BAR\n"
	ctx := &Context{
		ProgramName:    "edit test",
		ProgramVersion: "1.2.3",
		ADIFVersion:    "3.1.4",
		OutputFormat:   adif.FormatADI,
		Readers:        map[adif.Format]adif.Reader{adif.FormatADI: adi, adif.FormatCSV: csv},
		Writers:        map[adif.Format]adif.Writer{adif.FormatADI: adi, adif.FormatCSV: csv},
		Out:            out,
		fs:             fakeFilesystem{map[string]string{"foo.csv": file1}},
		CommandCtx: &editContext{
			add:    fieldAssignments{values: []adif.Field{{Name: "BAZ", Value: "Baz value"}}, validate: validateAlphanumName},
			set:    fieldAssignments{values: []adif.Field{{Name: "FOO", Value: "Foo value"}}, validate: validateAlphanumName},
			remove: []string{"BAR"},
		}}
	if err := Edit.Run(ctx, []string{"foo.csv"}); err != nil {
		t.Errorf("Edit.Run(ctx, foo.csv) got error %v", err)
	} else {
		got := out.String()
		want := "My Comment\n<ADIF_VER:5>3.1.4 <PROGRAMID:9>edit test <PROGRAMVERSION:5>1.2.3 <EOH>\n"
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Edit.Run(ctx, foo.csv) unexpected output, diff:\n%s", diff)
		}
	}
}

func TestEditAddSetRemove(t *testing.T) {
	adi := adif.NewADIIO()
	adi.HeaderCommentFn = func(*adif.Logfile) string { return "My Comment" }
	out := &bytes.Buffer{}
	file1 := `<FOO:7>old foo <BAR:7>old bar <CALL:4>W1AW <EOR>
<BAZ:7>old baz <CALL:3>N0P <APP_MONOLOG_FOO:7>monofoo <EOR>
<foo:4>foo2 <bar:4>bar2 <baz:4>baz2 <app_monolog_bar:7>monobar <eor>
`
	ctx := &Context{
		ProgramName:    "edit test",
		ProgramVersion: "1.2.3",
		ADIFVersion:    "3.1.4",
		OutputFormat:   adif.FormatADI,
		Readers:        map[adif.Format]adif.Reader{adif.FormatADI: adi},
		Writers:        map[adif.Format]adif.Writer{adif.FormatADI: adi},
		Out:            out,
		fs:             fakeFilesystem{map[string]string{"foo.adi": file1}},
		CommandCtx: &editContext{
			add:    fieldAssignments{values: []adif.Field{{Name: "BAZ", Value: "Baz value"}}, validate: validateAlphanumName},
			set:    fieldAssignments{values: []adif.Field{{Name: "FOO", Value: "Foo value"}}, validate: validateAlphanumName},
			remove: []string{"BAR"},
		}}
	if err := Edit.Run(ctx, []string{"foo.adi"}); err != nil {
		t.Errorf("Edit.Run(ctx) got error %v", err)
	} else {
		got := out.String()
		want := `My Comment
<ADIF_VER:5>3.1.4 <PROGRAMID:9>edit test <PROGRAMVERSION:5>1.2.3 <EOH>
<FOO:9>Foo value <CALL:4>W1AW <BAZ:9>Baz value <EOR>
<BAZ:7>old baz <CALL:3>N0P <APP_MONOLOG_FOO:7>monofoo <FOO:9>Foo value <EOR>
<FOO:9>Foo value <BAZ:4>baz2 <APP_MONOLOG_BAR:7>monobar <EOR>
`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Edit.Run(ctx, foo.adi) unexpected output, diff:\n%s", diff)
		}
	}
}

func TestEditRemoveEmpty(t *testing.T) {
	adi := adif.NewADIIO()
	adi.HeaderCommentFn = func(*adif.Logfile) string { return "My Comment" }
	csv := adif.NewCSVIO()
	out := &bytes.Buffer{}
	file1 := `FOO,BAR,BAZ
foo1,,baz1
,bar2,baz2
foo3,bar3,
,bar4,
,,
`
	ctx := &Context{
		ProgramName:    "edit test",
		ProgramVersion: "1.2.3",
		ADIFVersion:    "3.1.4",
		OutputFormat:   adif.FormatADI,
		Readers:        map[adif.Format]adif.Reader{adif.FormatADI: adi, adif.FormatCSV: csv},
		Writers:        map[adif.Format]adif.Writer{adif.FormatADI: adi, adif.FormatCSV: csv},
		Out:            out,
		fs:             fakeFilesystem{map[string]string{"foo.csv": file1}},
		CommandCtx:     &editContext{removeBlank: true}}
	if err := Edit.Run(ctx, []string{"foo.csv"}); err != nil {
		t.Errorf("Edit.Run(ctx, foo.csv) got error %v", err)
	} else {
		got := out.String()
		want := `My Comment
<ADIF_VER:5>3.1.4 <PROGRAMID:9>edit test <PROGRAMVERSION:5>1.2.3 <EOH>
<FOO:4>foo1 <BAZ:4>baz1 <EOR>
<BAR:4>bar2 <BAZ:4>baz2 <EOR>
<FOO:4>foo3 <BAR:4>bar3 <EOR>
<BAR:4>bar4 <EOR>
`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Edit.Run(ctx, foo.adi) unexpected output, diff:\n%s", diff)
		}
	}
}

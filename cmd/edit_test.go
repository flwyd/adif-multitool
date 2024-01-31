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
	"time"

	"github.com/flwyd/adif-multitool/adif"
	"github.com/google/go-cmp/cmp"
)

func TestEditEmpty(t *testing.T) {
	adi := adif.NewADIIO()
	csv := adif.NewCSVIO()
	out := &bytes.Buffer{}
	file1 := "FOO,BAR\n"
	ctx := &Context{
		OutputFormat: adif.FormatADI,
		Readers:      readers(adi, csv),
		Writers:      writers(adi, csv),
		Out:          out,
		Prepare:      testPrepare("My Comment", "3.1.4", "edit test", "1.2.3"),
		fs:           fakeFilesystem{map[string]string{"foo.csv": file1}},
		CommandCtx: &EditContext{
			Add:    FieldAssignments{values: []adif.Field{{Name: "BAZ", Value: "Baz value"}}, validate: ValidateAlphanumName},
			Set:    FieldAssignments{values: []adif.Field{{Name: "FOO", Value: "Foo value"}}, validate: ValidateAlphanumName},
			Rename: FieldAssignments{values: []adif.Field{{Name: "OLD", Value: "NEW"}}},
			Remove: []string{"BAR"},
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

func TestEditAddSetRemoveRename(t *testing.T) {
	adi := adif.NewADIIO()
	out := &bytes.Buffer{}
	file1 := `<FOO:7>old foo <BAR:7>old bar <CALL:4>W1AW <OLD:2:N>42 <EOR>
<BAZ:7>old baz <CALL:3>N0P <APP_MONOLOG_FOO:7>monofoo <EOR>
<old:8:d>20131031 <foo:4>foo2 <bar:4>bar2 <baz:4>baz2 <app_monolog_bar:7>monobar <eor>
`
	ctx := &Context{
		OutputFormat: adif.FormatADI,
		Readers:      readers(adi),
		Writers:      writers(adi),
		Out:          out,
		Prepare:      testPrepare("My Comment", "3.1.4", "edit test", "1.2.3"),
		fs:           fakeFilesystem{map[string]string{"foo.adi": file1}},
		CommandCtx: &EditContext{
			Add:    FieldAssignments{values: []adif.Field{{Name: "BAZ", Value: "Baz value"}}, validate: ValidateAlphanumName},
			Set:    FieldAssignments{values: []adif.Field{{Name: "FOO", Value: "Foo value"}}, validate: ValidateAlphanumName},
			Rename: FieldAssignments{values: []adif.Field{{Name: "OLD", Value: "NEW"}}},
			Remove: []string{"BAR"},
		}}
	if err := Edit.Run(ctx, []string{"foo.adi"}); err != nil {
		t.Errorf("Edit.Run(ctx) got error %v", err)
	} else {
		got := out.String()
		want := `My Comment
<ADIF_VER:5>3.1.4 <PROGRAMID:9>edit test <PROGRAMVERSION:5>1.2.3 <EOH>
<FOO:9>Foo value <CALL:4>W1AW <NEW:2:N>42 <BAZ:9>Baz value <EOR>
<BAZ:7>old baz <CALL:3>N0P <APP_MONOLOG_FOO:7>monofoo <FOO:9>Foo value <EOR>
<NEW:8:D>20131031 <FOO:9>Foo value <BAZ:4>baz2 <APP_MONOLOG_BAR:7>monobar <EOR>
`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Edit.Run(ctx, foo.adi) unexpected output, diff:\n%s", diff)
		}
	}
}

func TestEditIf(t *testing.T) {
	adi := adif.NewADIIO()
	out := &bytes.Buffer{}
	file1 := `<FOO:7>old foo <BAR:7>old bar <OLD:2>CA <CALL:4>W1AW <BAND:3>20M <MODE:2>CW <EOR>
<FOO:7>old foo <BAR:7>old bar <OLD:2>US <CALL:4>W1AW <BAND:3>40M <MODE:2>CW <EOR>
<BAZ:7>old baz <OLD:2>MX <CALL:3>N0P <APP_MONOLOG_FOO:7>monofoo <BAND:3>40m <MODE:3>SSB <EOR>
<BAZ:7>old baz <OLD:2>BZ <CALL:3>N0P <APP_MONOLOG_FOO:7>monofoo <BAND:3>40M <MODE:2>cw <EOR>
<foo:4>foo2 <bar:4>bar2 <baz:4>baz2 <old:2>GT <app_monolog_bar:7>monobar <BAND:3>40M <eor>
<foo:4>foo2 <bar:4>bar2 <baz:4>baz2 <old:2>SV <app_monolog_bar:7>monobar <BAND:3>40m <MODE:2>CW <eor>
`
	cond := ConditionValue{}
	cond.IfFlag().Set("MODE=CW")
	cond.IfFlag().Set("band=40m")
	ctx := &Context{
		OutputFormat: adif.FormatADI,
		Readers:      readers(adi),
		Writers:      writers(adi),
		Out:          out,
		Prepare:      testPrepare("My Comment", "3.1.4", "edit test", "1.2.3"),
		fs:           fakeFilesystem{map[string]string{"foo.adi": file1}},
		CommandCtx: &EditContext{
			Cond:   cond,
			Add:    FieldAssignments{values: []adif.Field{{Name: "BAZ", Value: "Baz value"}}, validate: ValidateAlphanumName},
			Set:    FieldAssignments{values: []adif.Field{{Name: "FOO", Value: "Foo value"}}, validate: ValidateAlphanumName},
			Rename: FieldAssignments{values: []adif.Field{{Name: "OLD", Value: "NEW"}}},
			Remove: []string{"BAR"},
		}}
	if err := Edit.Run(ctx, []string{"foo.adi"}); err != nil {
		t.Errorf("Edit.Run(ctx) got error %v", err)
	} else {
		got := out.String()
		want := `My Comment
<ADIF_VER:5>3.1.4 <PROGRAMID:9>edit test <PROGRAMVERSION:5>1.2.3 <EOH>
<FOO:7>old foo <BAR:7>old bar <OLD:2>CA <CALL:4>W1AW <BAND:3>20M <MODE:2>CW <EOR>
<FOO:9>Foo value <NEW:2>US <CALL:4>W1AW <BAND:3>40M <MODE:2>CW <BAZ:9>Baz value <EOR>
<BAZ:7>old baz <OLD:2>MX <CALL:3>N0P <APP_MONOLOG_FOO:7>monofoo <BAND:3>40m <MODE:3>SSB <EOR>
<BAZ:7>old baz <NEW:2>BZ <CALL:3>N0P <APP_MONOLOG_FOO:7>monofoo <BAND:3>40M <MODE:2>cw <FOO:9>Foo value <EOR>
<FOO:4>foo2 <BAR:4>bar2 <BAZ:4>baz2 <OLD:2>GT <APP_MONOLOG_BAR:7>monobar <BAND:3>40M <EOR>
<FOO:9>Foo value <BAZ:4>baz2 <NEW:2>SV <APP_MONOLOG_BAR:7>monobar <BAND:3>40m <MODE:2>CW <EOR>
`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Edit.Run(ctx, foo.adi) unexpected output, diff:\n%s", diff)
		}
	}
}

func TestRenameDoesNotClobber(t *testing.T) {
	adi := adif.NewADIIO()
	// output is empty string if error expected
	tests := []struct{ name, input, output string }{
		{
			name:   "new has value and before old",
			input:  "<STX:1>1 <NEW:9>new value <OLD:9>old value <OTHER:11>other value <EOR>",
			output: "",
		},
		{
			name:   "new has value and after old",
			input:  "<STX:1>2 <OLD:9>old value <NEW:9>new value <OTHER:11>other value <EOR>",
			output: "",
		},
		{
			name:   "new is empty",
			input:  "<STX:1>3 <OLD:9>old value <NEW:0> <OTHER:11>other value <EOR>",
			output: "<STX:1>3 <NEW:9>old value <OTHER:11>other value <EOR>",
		},
		{
			name:   "new is empty and before old",
			input:  "<STX:1>4 <NEW:0> <OLD:9>old value <OTHER:11>other value <EOR>",
			output: "<STX:1>4 <NEW:9>old value <OTHER:11>other value <EOR>",
		},
		{
			name:   "new is not set",
			input:  "<STX:1>5 <OLD:9>old value <OTHER:11>other value <EOR>",
			output: "<STX:1>5 <NEW:9>old value <OTHER:11>other value <EOR>",
		},
		{
			name:   "old is empty",
			input:  "<STX:1>6 <OLD:0> <NEW:9>new value <OTHER:11>other value <EOR>",
			output: "<STX:1>6 <NEW:9>new value <OTHER:11>other value <EOR>",
		},
		{
			name:   "old is empty and after new",
			input:  "<STX:1>7 <NEW:9>new value <OLD:0> <OTHER:11>other value <EOR>",
			output: "<STX:1>7 <NEW:9>new value <OTHER:11>other value <EOR>",
		},
		{
			name:   "old is not set",
			input:  "<STX:1>8 <NEW:9>new value <OTHER:11>other value <EOR>",
			output: "<STX:1>8 <NEW:9>new value <OTHER:11>other value <EOR>",
		},
	}
	wanthead := "<ADIF_VER:5>3.1.4 <PROGRAMID:9>edit test <PROGRAMVERSION:5>1.2.3 <EOH>"
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			ctx := &Context{
				OutputFormat: adif.FormatADI,
				Readers:      readers(adi),
				Writers:      writers(adi),
				Out:          out,
				Prepare:      testPrepare(tc.name, "3.1.4", "edit test", "1.2.3"),
				fs:           fakeFilesystem{map[string]string{"foo.adi": tc.input}},
				CommandCtx: &EditContext{
					Rename: FieldAssignments{values: []adif.Field{{Name: "OLD", Value: "NEW"}}},
				}}
			if err := Edit.Run(ctx, []string{"foo.adi"}); err != nil {
				if tc.output != "" {
					t.Errorf("Edit.Run(ctx) got error %v", err)
				}
			} else if tc.output == "" {
				t.Errorf("Edit.Run(ctx) want error, got %s", out.String())
			} else {
				got := out.String()
				want := fmt.Sprintf("%s\n%s\n%s\n", tc.name, wanthead, tc.output)
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("Edit.Run(ctx, foo.adi) unexpected output, diff:\n%s", diff)
				}
			}
		})
	}
}

func TestCyclicRename(t *testing.T) {
	adi := adif.NewADIIO()
	out := &bytes.Buffer{}
	file1 := `<STX:1>1 <OLD:9>old value <NEW:9>new value <OTHER:11>other value <EOR>
<STX:1>2 <OLD:9>old value <OTHER:11>other value <EOR>
<STX:1>3 <OLD:9>old value <NEW:9>new value <EOR>
<STX:1>4 <NEW:9>new value <OTHER:11>other value <EOR>
<STX:1>5 <NEW:9>new value <OLD:9>old value <OTHER:11>other value <EOR>
<STX:1>6 <OTHER:11>other value <NEW:9>new value <OLD:9>old value <EOR>
`
	ctx := &Context{
		OutputFormat: adif.FormatADI,
		Readers:      readers(adi),
		Writers:      writers(adi),
		Out:          out,
		Prepare:      testPrepare("My Comment", "3.1.4", "edit test", "1.2.3"),
		fs:           fakeFilesystem{map[string]string{"foo.adi": file1}},
		CommandCtx: &EditContext{
			Rename: FieldAssignments{values: []adif.Field{{Name: "OLD", Value: "NEW"}, {Name: "NEW", Value: "OTHER"}, {Name: "OTHER", Value: "OLD"}}},
		}}
	if err := Edit.Run(ctx, []string{"foo.adi"}); err != nil {
		t.Errorf("Edit.Run(ctx) got error %v", err)
	} else {
		got := out.String()
		want := `My Comment
<ADIF_VER:5>3.1.4 <PROGRAMID:9>edit test <PROGRAMVERSION:5>1.2.3 <EOH>
<STX:1>1 <NEW:9>old value <OTHER:9>new value <OLD:11>other value <EOR>
<STX:1>2 <NEW:9>old value <OLD:11>other value <EOR>
<STX:1>3 <NEW:9>old value <OTHER:9>new value <EOR>
<STX:1>4 <OTHER:9>new value <OLD:11>other value <EOR>
<STX:1>5 <OTHER:9>new value <NEW:9>old value <OLD:11>other value <EOR>
<STX:1>6 <OLD:11>other value <OTHER:9>new value <NEW:9>old value <EOR>
`
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Edit.Run(ctx, foo.adi) unexpected output, diff:\n%s", diff)
		}
	}
}

func TestEditRemoveEmpty(t *testing.T) {
	adi := adif.NewADIIO()
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
		OutputFormat: adif.FormatADI,
		Readers:      readers(adi, csv),
		Writers:      writers(adi, csv),
		Out:          out,
		Prepare:      testPrepare("My Comment", "3.1.4", "edit test", "1.2.3"),
		fs:           fakeFilesystem{map[string]string{"foo.csv": file1}},
		CommandCtx:   &EditContext{RemoveBlank: true}}
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

func TestAdjustTimeZone(t *testing.T) {
	type state struct {
		dateOn, dateOff, timeOn, timeOff string
	}
	toCsv := func(s state) string {
		return fmt.Sprintf("%s,%s,%s,%s", s.dateOn, s.timeOn, s.dateOff, s.timeOff)
	}
	zones := make(map[string]*time.Location)
	for _, z := range []string{"UTC", "Asia/Bangkok", "Asia/Kolkata", "Asia/Shanghai", "America/New_York"} {
		l, err := time.LoadLocation(z)
		if err != nil {
			t.Fatalf("could not load time zone %s: %v", z, err)
		}
		zones[z] = l
	}
	tests := []struct {
		start, want state
		from, to    *time.Location
		wantErr     bool
	}{
		{
			start: state{dateOn: "", dateOff: "", timeOn: "2008", timeOff: "0102"},
			want:  state{dateOn: "", dateOff: "", timeOn: "1208", timeOff: "1702"},
			from:  zones["Asia/Shanghai"], to: zones["UTC"],
			wantErr: true,
		},
		{
			start: state{dateOn: "", dateOff: "", timeOn: "2008", timeOff: "0102"},
			want:  state{dateOn: "", dateOff: "", timeOn: "2008", timeOff: "0102"},
			from:  zones["America/New_York"], to: zones["America/New_York"],
			wantErr: false,
		},
		{
			start: state{dateOn: "20080808", dateOff: "20080809", timeOn: "2008", timeOff: "0102"},
			want:  state{dateOn: "20080808", dateOff: "20080808", timeOn: "1208", timeOff: "1702"},
			from:  zones["Asia/Shanghai"], to: zones["UTC"],
		},
		{
			start: state{dateOn: "20201231", dateOff: "20201231", timeOn: "231545", timeOff: "232010"},
			want:  state{dateOn: "20210101", dateOff: "20210101", timeOn: "041545", timeOff: "042010"},
			from:  zones["America/New_York"], to: zones["UTC"],
		},
		{
			start: state{dateOn: "20210519", dateOff: "20210519", timeOn: "123456", timeOff: "1312"},
			want:  state{dateOn: "20210519", dateOff: "20210519", timeOn: "110456", timeOff: "1142"},
			from:  zones["Asia/Bangkok"], to: zones["Asia/Kolkata"],
		},
	}
	header := "QSO_DATE,TIME_ON,QSO_DATE_OFF,TIME_OFF"
	for _, tc := range tests {
		csv := adif.NewCSVIO()
		out := &bytes.Buffer{}
		file1 := fmt.Sprintf("%s\n%s\n", header, toCsv(tc.start))
		ctx := &Context{
			OutputFormat: adif.FormatCSV,
			Readers:      readers(csv),
			Writers:      writers(csv),
			Out:          out,
			fs:           fakeFilesystem{map[string]string{"foo.csv": file1}},
			CommandCtx:   &EditContext{FromZone: TimeZone{tz: tc.from}, ToZone: TimeZone{tz: tc.to}}}
		if tc.wantErr {
			if err := Edit.Run(ctx, []string{"foo.csv"}); err == nil {
				got := out.String()
				t.Errorf("edit -time-zone-from=%s -time-zone-to=%s from\n%swant error, got\n%s", tc.from, tc.to, file1, got)
			}
		} else if err := Edit.Run(ctx, []string{"foo.csv"}); err != nil {
			t.Errorf("edit -time-zone-from=%s -time-zone-to=%s from\n%sgot error %v", tc.from, tc.to, file1, err)
		} else {
			got := out.String()
			want := fmt.Sprintf("%s\n%s\n", header, toCsv(tc.want))
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("edit -time-zone-from=%s -time-zone-to=%s from\n%sgot diff\n%s", tc.from, tc.to, file1, diff)
			}
		}
	}
}

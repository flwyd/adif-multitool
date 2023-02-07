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

func TestEditAddSetRemove(t *testing.T) {
	adi := adif.NewADIIO()
	out := &bytes.Buffer{}
	file1 := `<FOO:7>old foo <BAR:7>old bar <CALL:4>W1AW <EOR>
<BAZ:7>old baz <CALL:3>N0P <APP_MONOLOG_FOO:7>monofoo <EOR>
<foo:4>foo2 <bar:4>bar2 <baz:4>baz2 <app_monolog_bar:7>monobar <eor>
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
			Remove: []string{"BAR"},
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

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

func TestValidateEmpty(t *testing.T) {
	io := adif.NewADIIO()
	out := &bytes.Buffer{}
	ctx := &Context{
		InputFormat:  adif.FormatADI,
		OutputFormat: adif.FormatADI,
		Readers:      readers(io),
		Writers:      writers(io),
		Out:          out,
		Prepare:      testPrepare("My Comment", "3.1.4", "validate test", "1.2.3"),
		fs:           fakeFilesystem{map[string]string{"-": ""}}}
	if err := Validate.Run(ctx, []string{}); err != nil {
		t.Errorf("Validate.Run(ctx) got error %v", err)
	} else {
		got := out.String()
		want := "My Comment\n<ADIF_VER:5>3.1.4 <PROGRAMID:13>validate test <PROGRAMVERSION:5>1.2.3 <EOH>\n"
		if got != want {
			t.Errorf("Validate.Run(ctx) got %s, want %s", got, want)
		}
	}
}

func TestValidateNoErrors(t *testing.T) {
	adi := adif.NewADIIO()
	adi.FieldSep = adif.SeparatorSpace
	adi.RecordSep = adif.SeparatorNewline
	out := &bytes.Buffer{}
	file1 := `My Comment
<ADIF_VER:5>3.1.4 <PROGRAMID:13>validate test <PROGRAMVERSION:5>1.2.3 <EOH>
<QSO_DATE:8>19901031 <TIME_ON:4>1234 <BAND:3>40M <FREQ:5>7.123 <MODE:2>CW <CALLSIGN:4>W1AW <NAME:17>Hiram Percy Maxim <AGE:2>31 <ARRL_SECT:2>CT <CONT:2>NA <GRIDSQUARE:6>FN31pr <SILENT_KEY:1>y <DXCC:3>291 <COUNTRY:24>UNITED STATES OF AMERICA <EOR>
<QSO_DATE:8:D>20221224 <TIME_ON:6:T>095846 <BAND:6:E>1.25cm <FREQ:5>24240 <MODE:3>PSK <SUBMODE:7>QPSK500 <CALLSIGN:3:S>N0P <NAME:11:S>Santa Claus <EMAIL:23>santa AT north DOT pole <QSO_RANDOM:1>N <IOTA_REF:6>EU-019 <K_INDEX:1>9 <GRIDSQUARE:4>LR49 <CONT:2>AS <DXCC:2>54 <COUNTRY:15>EUROPEAN RUSSIA <EOR>
`
	ctx := &Context{
		OutputFormat: adif.FormatADI,
		Readers:      readers(adi),
		Writers:      writers(adi),
		Out:          out,
		Prepare:      testPrepare("My Comment", "3.1.4", "validate test", "1.2.3"),
		fs:           fakeFilesystem{map[string]string{"foo.adi": file1}}}
	if err := Validate.Run(ctx, []string{"foo.adi"}); err != nil {
		t.Errorf("Validate.Run(ctx) got error on file without problems: %v", err)
	} else {
		got := out.String()
		want := file1
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Validate.Run(ctx, foo.adi, bar.adi) unexpected output, diff:\n%s", diff)
		}
	}
}

func TestValidateErrors(t *testing.T) {
	tests := []struct {
		name    string
		record  []adif.Field
		userdef []adif.UserdefField
	}{
		{
			name:   "non-ascii string",
			record: []adif.Field{{Name: "NAME", Value: "Pedro Peña"}},
		},
		{
			name:   "control characters",
			record: []adif.Field{{Name: "Notes", Value: "Vertical\vtab"}},
		},
		{
			name:   "invalid date",
			record: []adif.Field{{Name: "QSO_DATE", Value: "19754269"}},
		},
		{
			name:   "non-numeric format",
			record: []adif.Field{{Name: "ALTITUDE", Value: "forty-two"}},
		},
		{
			name:   "number out of range",
			record: []adif.Field{{Name: "AGE", Value: "345"}},
		},
		{
			name:   "wrong enum value",
			record: []adif.Field{{Name: "BAND", Value: "70m"}},
		},
		{
			name:    "userdef invalid number",
			record:  []adif.Field{{Name: "LUGGAGE_CODE", Value: "IZEA"}},
			userdef: []adif.UserdefField{{Name: "LUGGAGE_CODE", Type: adif.TypeNumber}},
		},
		{
			name:    "userdef number out of range",
			record:  []adif.Field{{Name: "YAGI_AZIMUTH_RADIANS", Value: "6.5"}},
			userdef: []adif.UserdefField{{Name: "YAGI_AZIMUTH_RADIANS", Type: adif.TypeNumber, Min: 0, Max: 6.283}},
		},
		{
			name:    "userdef non-matching enum",
			record:  []adif.Field{{Name: "ADJECTIVE", Value: "PRETTY"}},
			userdef: []adif.UserdefField{{Name: "ADJECTIVE", Type: adif.TypeNumber, EnumValues: []string{"GOOD", "BAD", "UGLY"}}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			inlog := adif.NewLogfile()
			inlog.Records = append(inlog.Records, adif.NewRecord(tc.record...))
			for _, u := range tc.userdef {
				if err := inlog.AddUserdef(u); err != nil {
					t.Fatalf("could not add userdef field %s: %v", u, err)
				}
			}
			adi := adif.NewADIIO()
			infile := &bytes.Buffer{}
			if err := adi.Write(inlog, infile); err != nil {
				t.Fatalf("could not create input logfile: %v", err)
			}
			out := &bytes.Buffer{}
			ctx := &Context{
				OutputFormat: adif.FormatADI,
				Readers:      readers(adi),
				Writers:      writers(adi),
				Out:          out,
				Prepare:      testPrepare("My Comment", "3.1.4", "validate test", "1.2.3"),
				fs:           fakeFilesystem{map[string]string{"foo.adi": infile.String()}}}
			if err := Validate.Run(ctx, []string{"foo.adi"}); err == nil {
				t.Errorf("Validate.Run(ctx) want error, got output:\n%s", out)
			} else if out.String() != "" {
				t.Errorf("Validate.Run(ctx) want empty output, got error %v and output:\n%s", err, out)
			}
		})
	}
}

// TODO test warnings (which are printed to stderr)

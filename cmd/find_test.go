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
	"errors"
	"testing"

	"github.com/flwyd/adif-multitool/adif"
	"github.com/google/go-cmp/cmp"
)

func TestFindEmpty(t *testing.T) {
	adi := adif.NewADIIO()
	csv := adif.NewCSVIO()
	out := &bytes.Buffer{}
	file1 := "FOO,BAR\n"
	cond := ConditionValue{}
	cond.IfFlag().Set("FOO=yes")
	cond.OrIfFlag().Set("BAR>=0")
	ctx := &Context{
		OutputFormat: adif.FormatADI,
		Readers:      readers(adi, csv),
		Writers:      writers(adi, csv),
		Out:          out,
		Prepare:      testPrepare("My Comment", "3.1.4", "find test", "1.2.3"),
		fs:           fakeFilesystem{map[string]string{"foo.csv": file1}},
		CommandCtx:   &FindContext{Cond: cond}}
	if err := Find.Run(ctx, []string{"foo.csv"}); err != nil {
		t.Errorf("Find.Run(ctx, foo.csv) got error %v", err)
	} else {
		got := out.String()
		want := "My Comment\n<ADIF_VER:5>3.1.4 <PROGRAMID:9>find test <PROGRAMVERSION:5>1.2.3 <EOH>\n"
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Find.Run(ctx, foo.csv) unexpected output, diff:\n%s", diff)
		}
	}
}

func TestFindEmptyCondition(t *testing.T) {
	adi := adif.NewADIIO()
	csv := adif.NewCSVIO()
	out := &bytes.Buffer{}
	file1 := `QSO_DATE,TIME_ON,TIME_OFF,CALL,BAND,FREQ,STATE,DXCC,TX_PWR,MODE,SUBMODE
20190101,1301,,K1A,160m,1.810,ME,291,10,SSB,LSB
20190202,1402,1405,N2B,80m,3.502,NY,291,20,CW,
20190303,1503,,W3C,40m,7.203,PA,291,30,SSB,LSB
`
	cond := ConditionValue{}
	ctx := &Context{
		OutputFormat: adif.FormatCSV,
		Readers:      readers(adi, csv),
		Writers:      writers(adi, csv),
		Out:          out,
		Prepare:      testPrepare("My Comment", "3.1.4", "find test", "1.2.3"),
		fs:           fakeFilesystem{map[string]string{"foo.csv": file1}},
		CommandCtx:   &FindContext{Cond: cond}}
	if err := Find.Run(ctx, []string{"foo.csv"}); err != nil {
		t.Errorf("Find.Run(ctx, foo.csv) got error %v", err)
	} else {
		got := out.String()
		if diff := cmp.Diff(file1, got); diff != "" {
			t.Errorf("Find.Run(ctx, foo.csv) unexpected output, diff:\n%s", diff)
		}
	}
}

func TestFindCondition(t *testing.T) {
	adi := adif.NewADIIO()
	csv := adif.NewCSVIO()
	file1 := `QSO_DATE,TIME_ON,TIME_OFF,CALL,BAND,FREQ,STATE,DXCC,TX_PWR,MODE,SUBMODE
20190101,1301,,K1A,160m,1.810,ME,291,10,SSB,LSB
20190202,1402,1405,N2B,80m,3.502,NY,291,20,CW,
20190303,1503,,W3C,40m,7.203,PA,291,30,SSB,LSB
20190404,1604,1604,W4D,30m,10.104,FL,291,40,RTTY,
20190505,1705,1707,N5E,20m,14.077,TX,291,50,PSK,PSK31
20190606,1806,,KH6F,17m,18.160,HI,110,60,SSTV,
20190707,1907,1910,AL7G,15m,21.270,AK,6,70,SSB,USB
20190808,2008,,W8H,12m,24.908,MI,291,80,CW,
20190909,2109,2109,N9I,10m,28.909,IL,291,90,FM,
20191010,2210,2210,K0J,6m,50.110,CO,291,100,SSB,USB
`
	tests := []struct {
		name  string
		setup func(*ConditionValue) error
		want  []string // callsigns
	}{
		{
			name:  "single equal",
			setup: func(c *ConditionValue) error { return c.IfFlag().Set("MODE=SSB") },
			want:  []string{"K1A", "W3C", "AL7G", "K0J"},
		},
		{
			name:  "single equal single result",
			setup: func(c *ConditionValue) error { return c.IfFlag().Set("CALL=W3C") },
			want:  []string{"W3C"},
		},
		{
			name:  "multi equal",
			setup: func(c *ConditionValue) error { return c.IfFlag().Set("MODE=cw|rtty|psk") },
			want:  []string{"N2B", "W4D", "N5E", "W8H"},
		},
		{
			name:  "equal empty",
			setup: func(c *ConditionValue) error { return c.IfFlag().Set("TIME_OFF=") },
			want:  []string{"K1A", "W3C", "KH6F", "W8H"},
		},
		{
			name:  "equal field",
			setup: func(c *ConditionValue) error { return c.IfFlag().Set("TIME_OFF={TIME_ON}") },
			want:  []string{"W4D", "N9I", "K0J"},
		},
		{
			name:  "equal field",
			setup: func(c *ConditionValue) error { return c.IfFlag().Set("TIME_ON={TIME_OFF}") },
			want:  []string{"W4D", "N9I", "K0J"},
		},
		{
			name:  "less than number",
			setup: func(c *ConditionValue) error { return c.IfFlag().Set("freq<10.104") },
			want:  []string{"K1A", "N2B", "W3C"},
		},
		{
			name:  "less than string",
			setup: func(c *ConditionValue) error { return c.IfFlag().Set("MODE<RTTY") },
			want:  []string{"N2B", "N5E", "W8H", "N9I"},
		},
		{
			name:  "less than equal date",
			setup: func(c *ConditionValue) error { return c.IfFlag().Set("qso_date<=20190404") },
			want:  []string{"K1A", "N2B", "W3C", "W4D"},
		},
		{
			name:  "less than equal band",
			setup: func(c *ConditionValue) error { return c.IfFlag().Set("BAND<=17M") },
			want:  []string{"K1A", "N2B", "W3C", "W4D", "N5E", "KH6F"},
		},
		{
			name:  "greater than number",
			setup: func(c *ConditionValue) error { return c.IfFlag().Set("TX_PWR>50") },
			want:  []string{"KH6F", "AL7G", "W8H", "N9I", "K0J"},
		},
		{
			name:  "greater than empty",
			setup: func(c *ConditionValue) error { return c.IfFlag().Set("SUBMODE>") },
			want:  []string{"K1A", "W3C", "N5E", "AL7G", "K0J"},
		},
		{
			name:  "greater than field",
			setup: func(c *ConditionValue) error { return c.IfFlag().Set("time_off>{time_on}") },
			want:  []string{"N2B", "N5E", "AL7G"},
		},
		{
			name:  "greater than multi",
			setup: func(c *ConditionValue) error { return c.IfFlag().Set("freq>27|{dxcc}") },
			want:  []string{"AL7G", "N9I", "K0J"},
		},
		{
			name:  "greater than equal number",
			setup: func(c *ConditionValue) error { return c.IfFlag().Set("FREQ>=18.160") },
			want:  []string{"KH6F", "AL7G", "W8H", "N9I", "K0J"},
		},
		{
			name:  "greater than equal string",
			setup: func(c *ConditionValue) error { return c.IfFlag().Set("state>=HI") },
			want:  []string{"K1A", "N2B", "W3C", "N5E", "KH6F", "W8H", "N9I"},
		},
		{
			name:  "greater than equal no match",
			setup: func(c *ConditionValue) error { return c.IfFlag().Set("band>=2m") },
			want:  []string{},
		},
		{
			name:  "not equal enum",
			setup: func(c *ConditionValue) error { return c.IfNotFlag().Set("dxcc=291") },
			want:  []string{"KH6F", "AL7G"},
		},
		{
			name:  "not equal empty",
			setup: func(c *ConditionValue) error { return c.IfNotFlag().Set("SUBMODE=") },
			want:  []string{"K1A", "W3C", "N5E", "AL7G", "K0J"},
		},
		{
			name:  "not equal multi",
			setup: func(c *ConditionValue) error { return c.IfNotFlag().Set("state=ak|hi|tx") },
			want:  []string{"K1A", "N2B", "W3C", "W4D", "W8H", "N9I", "K0J"},
		},
		{
			name:  "not less than",
			setup: func(c *ConditionValue) error { return c.IfNotFlag().Set("qso_date<20190808") },
			want:  []string{"W8H", "N9I", "K0J"},
		},
		{
			name:  "not less than equal",
			setup: func(c *ConditionValue) error { return c.IfNotFlag().Set("TIME_OFF<=1707") },
			want:  []string{"AL7G", "N9I", "K0J"},
		},
		{
			name:  "not greater than",
			setup: func(c *ConditionValue) error { return c.IfNotFlag().Set("submode>lsb") },
			want:  []string{"K1A", "N2B", "W3C", "W4D", "KH6F", "W8H", "N9I"},
		},
		{
			name:  "not greater than equal",
			setup: func(c *ConditionValue) error { return c.IfNotFlag().Set("band>=40m") },
			want:  []string{"K1A", "N2B"},
		},
		{
			name: "between numbers",
			setup: func(c *ConditionValue) error {
				return errors.Join(c.IfFlag().Set("FREQ>7"), c.IfFlag().Set("FREQ<20"))
			},
			want: []string{"W3C", "W4D", "N5E", "KH6F"},
		},
		{
			name:  "or alone",
			setup: func(c *ConditionValue) error { return c.OrIfFlag().Set("time_on>1800") },
			want:  []string{"KH6F", "AL7G", "W8H", "N9I", "K0J"},
		},
		{
			name: "equal or equal",
			setup: func(c *ConditionValue) error {
				return errors.Join(c.IfFlag().Set("CALL=W3C"), c.OrIfFlag().Set("CALL=N5E"))
			},
			want: []string{"W3C", "N5E"},
		},
		{
			name: "equal or not equal",
			setup: func(c *ConditionValue) error {
				return errors.Join(c.IfFlag().Set("CALL=K0J"), c.OrIfNotFlag().Set("MODE=SSB"))
			},
			want: []string{"N2B", "W4D", "N5E", "KH6F", "W8H", "N9I", "K0J"},
		},
		{
			name: "or and",
			setup: func(c *ConditionValue) error {
				return errors.Join(c.IfFlag().Set("MODE=CW"), c.OrIfFlag().Set("mode=SSB"), c.IfFlag().Set("submode=LSB"))
			},
			want: []string{"K1A", "N2B", "W3C", "W8H"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cond := ConditionValue{}
			if err := tc.setup(&cond); err != nil {
				t.Fatalf("Error parsing conditional flags: %v", err)
			}
			ctx := &Context{
				OutputFormat: adif.FormatADI,
				Readers:      readers(adi, csv),
				Writers:      writers(adi, csv),
				Out:          out,
				Prepare:      testPrepare("My Comment", "3.1.4", "find test", "1.2.3"),
				fs:           fakeFilesystem{map[string]string{"foo.csv": file1}},
				CommandCtx:   &FindContext{Cond: cond},
			}
			if err := Find.Run(ctx, []string{"foo.csv"}); err != nil {
				t.Fatalf("Find.Run(ctx, foo.csv) got error %v", err)
			}
			gotlog, err := adi.Read(out)
			if err != nil {
				t.Fatalf("Find output could not be parsed: %v", err)
			}
			got := make([]string, len(gotlog.Records))
			for i, r := range gotlog.Records {
				if call, ok := r.Get("CALL"); !ok {
					t.Errorf("Record %d had no CALL: %v", i, r)
				} else {
					got[i] = call.Value
				}
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("%s %v got diff\n%s", tc.name, cond.Get(), diff)
			}
		})
	}
}

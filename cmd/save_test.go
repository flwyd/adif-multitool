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
	"encoding/xml"
	"os"
	"testing"

	"github.com/flwyd/adif-multitool/adif"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/maps"
)

func TestSaveInferFormat(t *testing.T) {
	adiio := adif.NewADIIO()
	adiio.FieldSep = adif.SeparatorSpace
	adiio.RecordSep = adif.SeparatorNewline
	adxio := adif.NewADXIO()
	csvio := adif.NewCSVIO()
	jsonio := adif.NewJSONIO()
	input := `Test ADI input file
<ADIF_VER:5>3.2.1
<PROGRAMID:7>monolog
<PROGRAMVERSION:5>9.8.7
<EOH>
<CALL:4>W1AW
<QSO_DATE:8>19870605
<BAND:3>40m
<EOR>
Comment about Santa on the Air
<CALL:3>N0P
<QSO_DATE:8>20221224
<BAND:2>2m
<EOR>
`
	tests := []struct {
		name     string
		filename string
		want     string
		wantErr  bool
	}{
		{
			name:     "unknown format",
			filename: "out.txt",
			wantErr:  true,
		},
		{
			name:     "infer ADI",
			filename: "out.adi",
			want: `Test comment
<ADIF_VER:5>3.1.4 <PROGRAMID:9>test save <PROGRAMVERSION:5>5.6.7 <EOH>
<CALL:4>W1AW <QSO_DATE:8>19870605 <BAND:3>40m <EOR>
Comment about Santa on the Air <CALL:3>N0P <QSO_DATE:8>20221224 <BAND:2>2m <EOR>
`,
		},
		{
			name:     "infer ADX",
			filename: "out.adx",
			want:     xml.Header + `<ADX><HEADER><!--Test comment--><ADIF_VER>3.1.4</ADIF_VER><PROGRAMID>test save</PROGRAMID><PROGRAMVERSION>5.6.7</PROGRAMVERSION></HEADER><RECORDS><RECORD><CALL>W1AW</CALL><QSO_DATE>19870605</QSO_DATE><BAND>40m</BAND></RECORD><RECORD><!--Comment about Santa on the Air--><CALL>N0P</CALL><QSO_DATE>20221224</QSO_DATE><BAND>2m</BAND></RECORD></RECORDS></ADX>`,
		},
		{
			name:     "infer CSV",
			filename: "out.csv",
			want: `CALL,QSO_DATE,BAND
W1AW,19870605,40m
N0P,20221224,2m
`,
		},
		{
			name:     "infer JSON",
			filename: "out.json",
			want: `{"HEADER":{"ADIF_VER":"3.1.4","PROGRAMID":"test save","PROGRAMVERSION":"5.6.7"},"RECORDS":[{"BAND":"40m","CALL":"W1AW","QSO_DATE":"19870605"},{"BAND":"2m","CALL":"N0P","QSO_DATE":"20221224"}]}
`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fs := fakeFilesystem{files: map[string]string{os.Stdin.Name(): input}}
			ctx := &Context{
				Readers:    map[adif.Format]adif.Reader{adif.FormatADI: adiio, adif.FormatADX: adxio, adif.FormatCSV: csvio, adif.FormatJSON: jsonio},
				Writers:    map[adif.Format]adif.Writer{adif.FormatADI: adiio, adif.FormatADX: adxio, adif.FormatCSV: csvio, adif.FormatJSON: jsonio},
				Out:        os.Stdout,
				CommandCtx: &SaveContext{},
				Prepare:    testPrepare("Test comment", "3.1.4", "test save", "5.6.7"),
				fs:         fs,
			}
			if err := runSave(ctx, []string{tc.filename}); err != nil {
				if !tc.wantErr {
					t.Errorf("runSave got error: %v", err)
				}
			} else if got, ok := fs.files[tc.filename]; !ok {
				t.Errorf("runSave didn't write to file %s", tc.filename)
			} else if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("runSave(%q) got diff\n%s", tc.filename, diff)
			}
		})
	}
}

func TestSaveFileTemplate(t *testing.T) {
	adiio := adif.NewADIIO()
	adiio.FieldSep = adif.SeparatorSpace
	adiio.RecordSep = adif.SeparatorNewline
	adxio := adif.NewADXIO()
	csvio := adif.NewCSVIO()
	jsonio := adif.NewJSONIO()
	header := "QSO_DATE,TIME_ON,BAND,MODE,COUNTRY,OPERATOR,NOTES\n"
	w1aw := "19870605,0830,40m,CW,UNITED STATES OF AMERICA,1AY,\n"
	rs0iss := "20200317,1951,2m,FM,,KH9ELF,Hams:in*space?\n"
	n0p := "20221224,1234,2m,FM,United States of America,KH9ELF/P,Asked for presents!\n"
	ve0b := "19450701,1345,1.25m,CW,Canada,1AY,\"tab\tnewline\nslash/symbol‽greek(γράμμα)emoji📻Ge'ez<ደብዳቤ>\"\n"
	input := header + w1aw + rs0iss + n0p + ve0b
	tests := []struct {
		template string
		files    map[string]string
	}{
		{template: "{MODE}on{BAND}.csv",
			files: map[string]string{
				"FMon2M.csv":    header + rs0iss + n0p,
				"CWon40M.csv":   header + w1aw,
				"CWon1.25M.csv": header + ve0b,
			},
		},
		{template: "{operator}:{country}.csv",
			files: map[string]string{
				"1AY:UNITED STATES OF AMERICA.csv":      header + w1aw,
				"KH9ELF-P:UNITED STATES OF AMERICA.csv": header + n0p,
				"KH9ELF:COUNTRY-EMPTY.csv":              header + rs0iss,
				"1AY:CANADA.csv":                        header + ve0b,
			},
		},
		{template: "dir/{NoTeS}.csv",
			files: map[string]string{
				"dir/NOTES-EMPTY.csv":         header + w1aw,
				"dir/HAMS-IN-SPACE-.csv":      header + rs0iss,
				"dir/ASKED FOR PRESENTS!.csv": header + n0p,
				"dir/TAB_NEWLINE_SLASH-SYMBOL‽GREEK-ΓΡΆΜΜΑ-EMOJI📻GE-EZ-ደብዳቤ-.csv": header + ve0b,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.template, func(t *testing.T) {
			fs := fakeFilesystem{files: map[string]string{os.Stdin.Name(): input}}
			ctx := &Context{
				Readers:    map[adif.Format]adif.Reader{adif.FormatADI: adiio, adif.FormatADX: adxio, adif.FormatCSV: csvio, adif.FormatJSON: jsonio},
				Writers:    map[adif.Format]adif.Writer{adif.FormatADI: adiio, adif.FormatADX: adxio, adif.FormatCSV: csvio, adif.FormatJSON: jsonio},
				Out:        os.Stdout,
				CommandCtx: &SaveContext{},
				fs:         fs,
			}
			if err := runSave(ctx, []string{tc.template}); err != nil {
				t.Fatalf("runSave(%q) got error: %v", tc.template, err)
			}
			for name, want := range tc.files {
				if got, ok := fs.files[name]; !ok {
					t.Errorf("runSave(%q) did not create %q, files are %v", tc.template, name, maps.Keys(fs.files))
				} else if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("runSave(%q) got diff on %q:\n%s", tc.template, name, diff)
				}
			}
		})
	}
}

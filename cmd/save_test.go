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
)

func TestSaveInferFormat(t *testing.T) {
	adiio := adif.NewADIIO()
	adiio.FieldSep = adif.SeparatorSpace
	adiio.RecordSep = adif.SeparatorNewline
	adiio.HeaderCommentFn = func(l *adif.Logfile) string { return "Test comment" }
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
<CALL:3>N0P <QSO_DATE:8>20221224 <BAND:2>2m <EOR>
`,
		},
		{
			name:     "infer ADX",
			filename: "out.adx",
			want:     xml.Header + `<ADX><HEADER><ADIF_VER>3.1.4</ADIF_VER><PROGRAMID>test save</PROGRAMID><PROGRAMVERSION>5.6.7</PROGRAMVERSION></HEADER><RECORDS><RECORD><CALL>W1AW</CALL><QSO_DATE>19870605</QSO_DATE><BAND>40m</BAND></RECORD><RECORD><CALL>N0P</CALL><QSO_DATE>20221224</QSO_DATE><BAND>2m</BAND></RECORD></RECORDS></ADX>`,
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
		fs := fakeFilesystem{files: map[string]string{os.Stdin.Name(): input}}
		t.Run(tc.name, func(t *testing.T) {
			ctx := &Context{
				ProgramName:    "test save",
				ProgramVersion: "5.6.7",
				ADIFVersion:    "3.1.4",
				Readers:        map[adif.Format]adif.Reader{adif.FormatADI: adiio, adif.FormatADX: adxio, adif.FormatCSV: csvio, adif.FormatJSON: jsonio},
				Writers:        map[adif.Format]adif.Writer{adif.FormatADI: adiio, adif.FormatADX: adxio, adif.FormatCSV: csvio, adif.FormatJSON: jsonio},
				Out:            os.Stdout,
				CommandCtx:     &SaveContext{},
				fs:             fs,
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

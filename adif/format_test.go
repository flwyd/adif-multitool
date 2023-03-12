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

package adif

import (
	"bufio"
	"encoding/xml"
	"strings"
	"testing"
)

func TestGuessFormatFromName(t *testing.T) {
	tests := []struct {
		name    string
		want    Format
		wantErr bool
	}{
		{name: "foo.adi", want: FormatADI},
		{name: "foo.adx", want: FormatADX},
		{name: "foo.csv", want: FormatCSV},
		{name: "foo.json", want: FormatJSON},
		{name: "foo.tsv", want: FormatTSV},
		{name: "bar.ADI", want: FormatADI},
		{name: "bar.ADX", want: FormatADX},
		{name: "bar.CSV", want: FormatCSV},
		{name: "bar.JSON", want: FormatJSON},
		{name: "bar.TSV", want: FormatTSV},
		{name: "BAZ.tmp.adx", want: FormatADX},
		{name: "/path/to/file.csv", want: FormatCSV},
		{name: "nodotcsv", wantErr: true},
		{name: "log.txt", wantErr: true},
		{name: "file.tsv.data", wantErr: true},
		{name: "compressed.adx.gz", wantErr: true},
		{name: "/path/to/files.adi/noext", wantErr: true},
	}
	for _, tc := range tests {
		if got, err := GuessFormatFromName(tc.name); err != nil {
			if !tc.wantErr {
				t.Errorf("GuessFormatFromName(%q) got error %v, want %s", tc.name, err, tc.want)
			}
		} else if got != tc.want {
			t.Errorf("GuessFormatFromName(%q) got %s, want %s", tc.name, got, tc.want)
		}
	}
}

func TestGuessFormatFromContent(t *testing.T) {
	shortSpace := "\r\n\t  "
	longSpace := strings.Repeat(" ", contentPeekSize)
	tests := []struct {
		name    string
		want    Format
		wantErr bool
		records int
		text    string
	}{
		{
			name:    "ADI no header",
			want:    FormatADI,
			records: 1,
			text:    "<CALL:4>W1AW <MODE:2>CW <EOR>",
		},
		{
			name:    "ADI field type no header",
			want:    FormatADI,
			records: 1,
			text:    "<CALL:4:S>W1AW <MODE:2>CW <EOR>",
		},
		{
			name:    "ADI no header lower case",
			want:    FormatADI,
			records: 1,
			text:    "<call:4>W1AW <mode:2>CW <eor>\n",
		},
		{
			name:    "ADI space header",
			want:    FormatADI,
			records: 1,
			text:    shortSpace + "<PROGRAMID:11>format test <EOH>\n<CALL:4>W1AW <MODE:2>CW <EOR>\n",
		},
		{
			name:    "ADI comment header",
			want:    FormatADI,
			records: 1,
			text: `Generated on 2011-11-22 at 02:15:23Z for WN4AZY

<adif_ver:5>3.0.5
<programid:7>MonoLog
<USERDEF1:8:N>QRP_ARCI
<EOH>

<cAlL:4>W1AW
<MoDe:2>CW
<EOR>`,
		},
		{
			name:    "ADX basic",
			want:    FormatADX,
			records: 1,
			text:    xml.Header + "\n<ADX>\n<RECORDS>\n<RECORD><CALL>W1AW</CALL><MODE>CW</MODE></RECORD></RECORDS>\n</ADX>\n",
		},
		{
			name:    "ADX leading spaec",
			want:    FormatADX,
			records: 1,
			text:    shortSpace + "<ADX>\n<RECORDS>\n<RECORD><CALL>W1AW</CALL><MODE>CW</MODE></RECORD></RECORDS>\n</ADX>\n",
		},
		{
			name: "ADX no records",
			want: FormatADX,
			text: "<ADX>\n<HEADER>\n</HEADER>\n<RECORDS>\n</RECORDS>\n</ADX>\n",
		},
		{
			name: "ADX no records no newlines",
			want: FormatADX,
			text: "<ADX><HEADER></HEADER><RECORDS></RECORDS></ADX>",
		},
		{
			name: "ADX XML comment",
			want: FormatADX,
			text: xml.Header + "\n<!-- This is my log file -->\n<ADX><HEADER></HEADER><RECORDS></RECORDS></ADX>\n",
		},
		{
			name:    "CSV basic",
			want:    FormatCSV,
			records: 1,
			text:    "CALL,MODE\nW1AW,CW\n",
		},
		{
			name:    "CSV Windows",
			want:    FormatCSV,
			records: 1,
			text:    "CALL,MODE\r\nW1AW,CW\r\n",
		},
		{
			name: "CSV header only",
			want: FormatCSV,
			text: "CALL,MODE\n",
		},
		{
			name:    "CSV leading space",
			want:    FormatCSV,
			records: 1,
			text:    shortSpace + "CALL,MODE\nW1AW,CW\n",
		},
		{
			name:    "JSON basic",
			want:    FormatJSON,
			records: 1,
			text: `{
				"HEADER": {},
				"RECORDS": [{"CALL": "W1AW", "MODE": "CW"}]
			}
			`,
		},
		{
			name:    "JSON no header",
			want:    FormatJSON,
			records: 1,
			text:    `{"RECORDS": [{"CALL": "W1AW", "MODE": "CW"}]}`,
		},
		{
			name:    "JSON leading spacce",
			want:    FormatJSON,
			records: 1,
			text:    shortSpace + `{"RECORDS": [{"CALL": "W1AW", "MODE": "CW"}]}`,
		},
		{
			name: "JSON no records", // could later decide this is invalid
			want: FormatJSON,
			text: `{}`,
		},
		{
			name:    "TSV basic",
			want:    FormatTSV,
			records: 1,
			text:    "CALL\tMODE\nW1AW\tCW\n",
		},
		{
			name:    "TSV Windows",
			want:    FormatTSV,
			records: 1,
			text:    "CALL\tMODE\r\nW1AW\tCW\r\n",
		},
		{
			name: "TSV header only",
			want: FormatTSV,
			text: "CALL\tMODE\n",
		},
		{
			name:    "TSV leading space",
			want:    FormatTSV,
			records: 1,
			text:    "  CALL\tMODE\nW1AW\tCW\n",
		},
		{
			name:    "ADI looks like CSV",
			want:    FormatADI,
			records: 1,
			text:    "Comment,Text\n<PROGRAMID:11>format test <EOH>\n<CALL:4>W1AW <MODE:2>CW <EOR>\n",
		},
		{
			name:    "ADI looks like JSON",
			want:    FormatADI,
			records: 1,
			text:    `{"RECORDS": []} Comment <PROGRAMID:11>format test <EOH> <CALL:4>W1AW <MODE:2>CW <EOR>`,
		},
		{
			name:    "ADI looks like TSV",
			want:    FormatADI,
			records: 1,
			text:    "Comment\tText\n<PROGRAMID:11>format test <EOH>\n<CALL:4>W1AW <MODE:2>CW <EOR>\n",
		},
		{
			name:    "HTML",
			wantErr: true,
			text:    "<HTML><BODY><RECORDS><RECORD><CALL>W1AW</CALL><MODE>CW</MODE></RECORD></RECORDS></BODY></HTML>",
		},
		{
			name:    "Empty",
			wantErr: true,
			text:    "",
		},
		{
			name:    "SpaceOnly",
			wantErr: true,
			text:    shortSpace,
		},
		{
			name:    "ADX with doctype",
			wantErr: true,
			text:    "<!DOCTYPE ADX><ADX><HEADER></HEADER><RECORDS></RECORDS></ADX>\n",
		},
		{
			name:    "Long leading space ADI",
			wantErr: true,
			text:    longSpace + "<CALL:4>W1AW <MODE:2>CW <EOR>\n",
		},
		{
			name:    "Long leading space ADX",
			wantErr: true,
			text:    longSpace + "<ADX><HEADER></HEADER><RECORDS></RECORDS></ADX>\n",
		},
		{
			name:    "Long leading space CSV",
			wantErr: true,
			text:    longSpace + "CALL,MODE\nW1AW,CW\n",
		},
		{
			name:    "Long leading space JSON",
			wantErr: true,
			text:    longSpace + `{"HEADER": {}, "RECORDS": []}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := bufio.NewReader(strings.NewReader(tc.text))
			if got, err := GuessFormatFromContent(r); err != nil {
				if !tc.wantErr {
					t.Errorf("GuessFormatFromContent got error %v, want %s, text was\n%s", err, tc.want, tc.text)
				}
			} else if got != tc.want {
				t.Errorf("GuessFormatFromContent got %s, want %s, text was\n%s", got, tc.want, tc.text)
			} else {
				var fr Reader
				switch got {
				case FormatADI:
					fr = NewADIIO()
				case FormatADX:
					fr = NewADXIO()
				case FormatCSV:
					fr = NewCSVIO()
				case FormatJSON:
					fr = NewJSONIO()
				case FormatTSV:
					fr = NewTSVIO()
				}
				if l, err := fr.Read(r); err != nil {
					t.Errorf("GuessFormatFromContent left Reader in an unsuitable state, %s.Read() got error %v, text:\n%s", fr, err, tc.text)
				} else if gotrec := len(l.Records); gotrec != tc.records {
					t.Errorf("GuessFormatFromContent got %d records, want %d\n%v", gotrec, tc.records, l)
				}
			}
		})
	}
}

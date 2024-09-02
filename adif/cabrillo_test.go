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
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestEmptyCabrillo(t *testing.T) {
	input := `START-OF-LOG: 3.0
EMAIL: ham@example.com
NAME: Ham Operator
CALLSIGN: W1AW
LOCATION: CO
GRID-LOCATOR: DN70
CONTEST: TEST-CONTEST-ID
CATEGORY-OPERATOR: SINGLE-OP
CATEGORY-TRANSMITTER: ONE
CATEGORY-ASSISTED: ASSISTED
CATEGORY-MODE: SSB
CATEGORY-POWER: LOW
CATEGORY-STATION: FIXED
CREATED-BY: HAND
CLAIMED-SCORE: 0
END-OF-LOG:
`
	cab := NewCabrilloIO()
	if parsed, err := cab.Read(strings.NewReader(input)); err != nil {
		t.Errorf("Read(%q) got error %v", input, err)
	} else {
		if got := len(parsed.Records); got != 0 {
			t.Errorf("Read(%q) got %d records, want %d", input, got, 0)
		}
	}
}

func TestReadCabrillo(t *testing.T) {
	input := `START-OF-LOG: 3.0
EMAIL: ham@example.com
NAME: Ham Operator
CALLSIGN: W1AW
LOCATION: CT
GRID-LOCATOR: FN31
CONTEST: TEST-CONTEST-ID
CATEGORY-OPERATOR: MULTI-OP
CATEGORY-TRANSMITTER: TWO
CATEGORY-ASSISTED: ASSISTED
CATEGORY-MODE: MIXED
CATEGORY-POWER: LOW
CATEGORY-STATION: ROVER
CREATED-BY: HAND
CLAIMED-SCORE: 16
QSO: 14234 PH 2023-10-31 1234 W1AW   57  CT  AA1A   48  PAC 0
QSO:  7012 CW 2023-11-01 0123 W1AW   599 CT  WX0YZ  432 MN  1
QSO:  1.2G DG 2023-11-01 1415 W1AW   23  CT  N7N    45  WWA 1
QSO: 222 FM 2023-10-31 2345 W1AW/M 59 CT KJ4LMN 59 NC 0
X-QSO: 21123 RY 2023-10-31 1920 W1AW/M 46 RI EA1OU 53 DX 0
END-OF-LOG:
`
	wantHeaders := []Field{
		{Name: "APP_CABRILLO_EMAIL", Value: "ham@example.com", Type: TypeString},
		{Name: "APP_CABRILLO_NAME", Value: "Ham Operator", Type: TypeString},
		{Name: "APP_CABRILLO_CALLSIGN", Value: "W1AW", Type: TypeString},
		{Name: "APP_CABRILLO_LOCATION", Value: "CT", Type: TypeString},
		{Name: "APP_CABRILLO_GRID_LOCATOR", Value: "FN31", Type: TypeString},
		{Name: "APP_CABRILLO_CONTEST", Value: "TEST-CONTEST-ID", Type: TypeString},
		{Name: "APP_CABRILLO_CATEGORY_OPERATOR", Value: "MULTI-OP", Type: TypeString},
		{Name: "APP_CABRILLO_CATEGORY_TRANSMITTER", Value: "TWO", Type: TypeString},
		{Name: "APP_CABRILLO_CATEGORY_ASSISTED", Value: "ASSISTED", Type: TypeString},
		{Name: "APP_CABRILLO_CATEGORY_MODE", Value: "MIXED", Type: TypeString},
		{Name: "APP_CABRILLO_CATEGORY_POWER", Value: "LOW", Type: TypeString},
		{Name: "APP_CABRILLO_CATEGORY_STATION", Value: "ROVER", Type: TypeString},
		{Name: "APP_CABRILLO_CREATED_BY", Value: "HAND", Type: TypeString},
		{Name: "APP_CABRILLO_CLAIMED_SCORE", Value: "16", Type: TypeString},
	}
	wantFields := [][]Field{
		{
			{Name: "FREQ", Value: "14.234", Type: TypeNumber},
			{Name: "BAND", Value: "20m", Type: TypeString},
			{Name: "MODE", Value: "SSB", Type: TypeString},
			{Name: "QSO_DATE", Value: "20231031", Type: TypeDate},
			{Name: "TIME_ON", Value: "1234", Type: TypeTime},
			{Name: "STATION_CALLSIGN", Value: "W1AW", Type: TypeString},
			{Name: "RST_SENT", Value: "57", Type: TypeString},
			{Name: "MY_ARRL_SECT", Value: "CT", Type: TypeString},
			{Name: "CALL", Value: "AA1A", Type: TypeString},
			{Name: "RST_RCVD", Value: "48", Type: TypeString},
			{Name: "ARRL_SECT", Value: "PAC", Type: TypeString},
			{Name: "APP_CABRILLO_TRANSMITTER_ID", Value: "0", Type: TypeNumber},
			{Name: "CONTEST_ID", Value: "TEST-CONTEST-ID", Type: TypeString},
			{Name: "GRIDSQUARE", Value: "FN31", Type: TypeString},
		}, {
			{Name: "FREQ", Value: "7.012", Type: TypeNumber},
			{Name: "BAND", Value: "40m", Type: TypeString},
			{Name: "MODE", Value: "CW", Type: TypeString},
			{Name: "QSO_DATE", Value: "20231101", Type: TypeDate},
			{Name: "TIME_ON", Value: "0123", Type: TypeTime},
			{Name: "STATION_CALLSIGN", Value: "W1AW", Type: TypeString},
			{Name: "RST_SENT", Value: "599", Type: TypeString},
			{Name: "MY_ARRL_SECT", Value: "CT", Type: TypeString},
			{Name: "CALL", Value: "WX0YZ", Type: TypeString},
			{Name: "RST_RCVD", Value: "432", Type: TypeString},
			{Name: "ARRL_SECT", Value: "MN", Type: TypeString},
			{Name: "APP_CABRILLO_TRANSMITTER_ID", Value: "1", Type: TypeNumber},
			{Name: "CONTEST_ID", Value: "TEST-CONTEST-ID", Type: TypeString},
			{Name: "GRIDSQUARE", Value: "FN31", Type: TypeString},
		}, {
			{Name: "BAND", Value: "23cm", Type: TypeString},
			{Name: "MODE", Value: "DIGITAL", Type: TypeString},
			{Name: "QSO_DATE", Value: "20231101", Type: TypeDate},
			{Name: "TIME_ON", Value: "1415", Type: TypeTime},
			{Name: "STATION_CALLSIGN", Value: "W1AW", Type: TypeString},
			{Name: "RST_SENT", Value: "23", Type: TypeString},
			{Name: "MY_ARRL_SECT", Value: "CT", Type: TypeString},
			{Name: "CALL", Value: "N7N", Type: TypeString},
			{Name: "RST_RCVD", Value: "45", Type: TypeString},
			{Name: "ARRL_SECT", Value: "WWA", Type: TypeString},
			{Name: "APP_CABRILLO_TRANSMITTER_ID", Value: "1", Type: TypeNumber},
			{Name: "CONTEST_ID", Value: "TEST-CONTEST-ID", Type: TypeString},
			{Name: "GRIDSQUARE", Value: "FN31", Type: TypeString},
		}, {
			{Name: "BAND", Value: "1.25m", Type: TypeString},
			{Name: "MODE", Value: "FM", Type: TypeString},
			{Name: "QSO_DATE", Value: "20231031", Type: TypeDate},
			{Name: "TIME_ON", Value: "2345", Type: TypeTime},
			{Name: "STATION_CALLSIGN", Value: "W1AW/M", Type: TypeString},
			{Name: "RST_SENT", Value: "59", Type: TypeString},
			{Name: "MY_ARRL_SECT", Value: "CT", Type: TypeString},
			{Name: "CALL", Value: "KJ4LMN", Type: TypeString},
			{Name: "RST_RCVD", Value: "59", Type: TypeString},
			{Name: "ARRL_SECT", Value: "NC", Type: TypeString},
			{Name: "APP_CABRILLO_TRANSMITTER_ID", Value: "0", Type: TypeNumber},
			{Name: "CONTEST_ID", Value: "TEST-CONTEST-ID", Type: TypeString},
			{Name: "GRIDSQUARE", Value: "FN31", Type: TypeString},
		}, {
			{Name: "FREQ", Value: "21.123", Type: TypeNumber},
			{Name: "BAND", Value: "15m", Type: TypeString},
			{Name: "MODE", Value: "RTTY", Type: TypeString},
			{Name: "QSO_DATE", Value: "20231031", Type: TypeDate},
			{Name: "TIME_ON", Value: "1920", Type: TypeTime},
			{Name: "STATION_CALLSIGN", Value: "W1AW/M", Type: TypeString},
			{Name: "RST_SENT", Value: "46", Type: TypeString},
			{Name: "MY_ARRL_SECT", Value: "RI", Type: TypeString},
			{Name: "CALL", Value: "EA1OU", Type: TypeString},
			{Name: "RST_RCVD", Value: "53", Type: TypeString},
			{Name: "ARRL_SECT", Value: "DX", Type: TypeString},
			{Name: "APP_CABRILLO_TRANSMITTER_ID", Value: "0", Type: TypeNumber},
			{Name: "APP_CABRILLO_XQSO", Value: "Y", Type: TypeBoolean},
			{Name: "CONTEST_ID", Value: "TEST-CONTEST-ID", Type: TypeString},
			{Name: "GRIDSQUARE", Value: "FN31", Type: TypeString},
		},
	}
	cab := &CabrilloIO{
		MyExchangeField:    "MY_ARRL_SECT",
		TheirExchangeField: "ARRL_SECT",
	}
	parsed, err := cab.Read(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Read(%q) got error %v", input, err)
	}
	less := func(a, b Field) bool {
		if a.Name == b.Name {
			return a.Value < b.Value
		}
		return a.Name < b.Name
	}
	if diff := cmp.Diff(wantHeaders, parsed.Header.Fields(), cmpopts.SortSlices(less)); diff != "" {
		t.Errorf("Read(%q) headers did not match expected, diff:\n%s", input, diff)
	}
	for i, r := range parsed.Records {
		fields := r.Fields()
		if diff := cmp.Diff(wantFields[i], fields); diff != "" {
			t.Errorf("Read(%q) record %d did not match expected, diff:\n%s", input, i+1, diff)
		}
	}
	if gotlen := len(parsed.Records); gotlen != len(wantFields) {
		t.Errorf("Read(%q) got %d records:\n%v\nwant %d\n%v", input, gotlen, parsed.Records[len(wantFields):], len(wantFields), wantFields)
	}
}

func TestWriteCabrillo(t *testing.T) {
	l := NewLogfile()
	l.AddRecord(NewRecord(
		Field{Name: "QSO_DATE", Value: "20231031", Type: TypeDate},
		Field{Name: "TIME_ON", Value: "1234", Type: TypeTime},
		Field{Name: "TIME_OFF", Value: "1239", Type: TypeTime},
		Field{Name: "FREQ", Value: "7.234", Type: TypeNumber},
		Field{Name: "BAND", Value: "20m"},
		Field{Name: "MODE", Value: "SSB"},
		Field{Name: "SUBMODE", Value: "USB"},
		Field{Name: "CALL", Value: "AA1A"},
		Field{Name: "STATION_CALLSIGN", Value: "W1AW"},
		Field{Name: "RST_SENT", Value: "57"},
		Field{Name: "RST_RCVD", Value: "48"},
		Field{Name: "ARRL_SECT", Value: "PAC"},
		Field{Name: "MY_ARRL_SECT", Value: "CT"},
		Field{Name: "APP_CABRILLO_TRANSMITTER_ID", Value: "0"},
	)).AddRecord(NewRecord(
		Field{Name: "FREQ", Value: "14.0461", Type: TypeNumber},
		Field{Name: "MODE", Value: "CW", Type: TypeString},
		Field{Name: "QSO_DATE", Value: "20231101", Type: TypeDate},
		Field{Name: "TIME_ON", Value: "0123", Type: TypeTime},
		Field{Name: "STATION_CALLSIGN", Value: "W1AW", Type: TypeString},
		Field{Name: "RST_SENT", Value: "599", Type: TypeString},
		Field{Name: "MY_ARRL_SECT", Value: "CT", Type: TypeString},
		Field{Name: "CALL", Value: "WX0YZ", Type: TypeString},
		Field{Name: "RST_RCVD", Value: "432", Type: TypeString},
		Field{Name: "ARRL_SECT", Value: "MN", Type: TypeString},
		Field{Name: "APP_CABRILLO_TRANSMITTER_ID", Value: "1", Type: TypeNumber},
	)).AddRecord(NewRecord(
		Field{Name: "BAND", Value: "23cm", Type: TypeString},
		Field{Name: "MODE", Value: "DIGITAL", Type: TypeString},
		Field{Name: "QSO_DATE", Value: "20231101", Type: TypeDate},
		Field{Name: "TIME_ON", Value: "1415", Type: TypeTime},
		Field{Name: "STATION_CALLSIGN", Value: "W1AW", Type: TypeString},
		Field{Name: "RST_SENT", Value: "23", Type: TypeString},
		Field{Name: "MY_ARRL_SECT", Value: "CT", Type: TypeString},
		Field{Name: "CALL", Value: "N7N", Type: TypeString},
		Field{Name: "RST_RCVD", Value: "45", Type: TypeString},
		Field{Name: "ARRL_SECT", Value: "WWA", Type: TypeString},
		Field{Name: "APP_CABRILLO_TRANSMITTER_ID", Value: "1", Type: TypeNumber},
	)).AddRecord(NewRecord(
		Field{Name: "BAND", Value: "1.25m", Type: TypeString},
		Field{Name: "MODE", Value: "FM", Type: TypeString},
		Field{Name: "QSO_DATE", Value: "20231031", Type: TypeDate},
		Field{Name: "TIME_ON", Value: "2345", Type: TypeTime},
		Field{Name: "STATION_CALLSIGN", Value: "W1AW/M", Type: TypeString},
		Field{Name: "RST_SENT", Value: "59", Type: TypeString},
		Field{Name: "MY_ARRL_SECT", Value: "CT", Type: TypeString},
		Field{Name: "CALL", Value: "KJ4LMN", Type: TypeString},
		Field{Name: "RST_RCVD", Value: "59", Type: TypeString},
		Field{Name: "ARRL_SECT", Value: "NC", Type: TypeString},
		Field{Name: "APP_CABRILLO_TRANSMITTER_ID", Value: "0", Type: TypeNumber},
	))
	l.Header.Set(Field{Name: "PROGRAMID", Value: "My Logger"})
	l.Header.Set(Field{Name: "PROGRAMVERSION", Value: "1.2.3"})
	l.Header.Set(Field{Name: "APP_CABRILLO_CLAIMED_SCORE", Value: "42"})
	l.Header.Set(Field{Name: "APP_CABRILLO_CLUB", Value: "Amateur Radio Relay League"})
	l.Header.Set(Field{Name: "APP_CABRILLO_ADDRESS", Value: "225 Main Street\nNewington, CT 06111"})
	l.Header.Set(Field{Name: "APP_CABRILLO_CALLSIGN", Value: "N9N"})
	l.Header.Set(Field{Name: "APP_CABRILLO_CATEGORY_BAND", Value: "14000"})
	l.Header.Set(Field{Name: "APP_CABRILLO_CATEGORY_OVERLAY", Value: "CLASSIC"})
	want := `START-OF-LOG: 3.0
X-INSTRUCTIONS: Fill out headers following contest instructions
X-INSTRUCTIONS: Delete any unnecessary headers
X-INSTRUCTIONS: Double-check QSO lines, keeping columns in order
X-INSTRUCTIONS: Report bugs at https://github.com/flwyd/adif-multitool
CREATED-BY: cabrillo_test
SOAPBOX: 
CONTEST: TEST-CONTEST-ID
CALLSIGN: W1AW (3 records) W1AW/M (1 records)
CLUB: Amateur Radio Relay League
OPERATORS: 
NAME: Ham Operator
EMAIL: ham@example.com
ADDRESS: 225 Main Street
ADDRESS: Newington, CT 06111
GRID-LOCATOR: FN31
LOCATION: CT
CLAIMED-SCORE: 42
OFFTIME: 
CATEGORY-ASSISTED: 
CATEGORY-BAND: ALL
CATEGORY-MODE: MIXED
CATEGORY-OPERATOR: 
CATEGORY-OVERLAY: CLASSIC
CATEGORY-POWER: 
CATEGORY-STATION: 
CATEGORY-TIME: 
CATEGORY-TRANSMITTER: 
X-INSTRUCTIONS: See contest rules for expected category values
X-Q:                            --info sent---- --info rcvd----
X-Q: freq    mo date       time call   rst exch call   rst exch t
QSO: 7234    PH 2023-10-31 1234 W1AW   57  CT   AA1A   48  PAC  0
QSO: 14046.1 CW 2023-11-01 0123 W1AW   599 CT   WX0YZ  432 MN   1
QSO: 1.2G    DG 2023-11-01 1415 W1AW   23  CT   N7N    45  WWA  1
QSO: 222     FM 2023-10-31 2345 W1AW/M 59  CT   KJ4LMN 59  NC   0
END-OF-LOG: 
`
	cab := &CabrilloIO{
		CreatedBy:          "cabrillo_test",
		Contest:            "TEST-CONTEST-ID",
		Email:              "ham@example.com",
		GridLocator:        "FN31",
		Name:               "Ham Operator",
		MyExchange:         "CT",
		MyExchangeField:    "MY_ARRL_SECT",
		TheirExchangeField: "ARRL_SECT"}
	out := &strings.Builder{}
	if err := cab.Write(l, out); err != nil {
		t.Errorf("Write(%v) got error %v", l, err)
	} else {
		got := out.String()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Write(%v) had diff with expected:\n%s", l, diff)
		}
	}
}

func TestInferrCabrilloCategories(t *testing.T) {
	tests := []struct {
		name                          string
		modes, bands, powers          []string
		wantMode, wantBand, wantPower string
	}{
		{
			name:      "just cw",
			modes:     []string{"CW", "CW", "CW", "CW"},
			bands:     []string{"40m", "20m", "15m", "10m"},
			powers:    []string{"5", "25", "50", "100"},
			wantMode:  "CW",
			wantBand:  "ALL",
			wantPower: "HIGH",
		},
		{
			name:      "just ssb",
			modes:     []string{"SSB", "SSB", "SSB", "SSB"},
			bands:     []string{"40m", "20m", "15m", "10m"},
			powers:    []string{"5", "25", "50", "100"},
			wantMode:  "SSB",
			wantBand:  "ALL",
			wantPower: "HIGH",
		},
		{
			name:      "just fm",
			modes:     []string{"FM", "FM", "FM", "FM"},
			bands:     []string{"10m", "6m", "2m", "70cm"},
			powers:    []string{"5", "25", "50", "100"},
			wantMode:  "FM",
			wantBand:  "ALL",
			wantPower: "HIGH",
		},
		{
			name:      "just rtty",
			modes:     []string{"RTTY", "RTTY", "RTTY", "RTTY"},
			bands:     []string{"40m", "20m", "15m", "10m"},
			powers:    []string{"5", "25", "50", "100"},
			wantMode:  "RTTY",
			wantBand:  "ALL",
			wantPower: "HIGH",
		},
		{
			name:      "multi digi",
			modes:     []string{"FT8", "OLIVIA", "PSK", "MFSK"},
			bands:     []string{"40m", "20m", "15m", "10m"},
			powers:    []string{"5", "25", "50", "100"},
			wantMode:  "DIGI",
			wantBand:  "ALL",
			wantPower: "HIGH",
		},
		{
			name:      "rtty and digi",
			modes:     []string{"RTTY", "RTTY", "RTTY", "HELL"},
			bands:     []string{"40m", "20m", "15m", "10m"},
			powers:    []string{"5", "25", "50", "100"},
			wantMode:  "DIGI",
			wantBand:  "ALL",
			wantPower: "HIGH",
		},
		{
			name:      "multi phone",
			modes:     []string{"DIGITALVOICE", "FM", "AM", "SSB"},
			bands:     []string{"70cm", "2m", "6m", "10m"},
			powers:    []string{"1", "2", "3", "4"},
			wantMode:  "SSB",
			wantBand:  "ALL",
			wantPower: "QRP",
		},
		{
			name:      "mixed mode",
			modes:     []string{"FT8", "SSB"},
			bands:     []string{"40m", "20m"},
			powers:    []string{"50", "100"},
			wantMode:  "MIXED",
			wantBand:  "ALL",
			wantPower: "HIGH",
		},
		{
			name:      "band 80m",
			modes:     []string{"FT8", "SSB", "CW"},
			bands:     []string{"80m", "80m", "80m"},
			powers:    []string{"5", "50", "100"},
			wantMode:  "MIXED",
			wantBand:  "80M",
			wantPower: "HIGH",
		},
		{
			name:      "band 10G",
			modes:     []string{"FT8", "FM", "ATV"},
			bands:     []string{"3cm", "3cm", "3cm"},
			powers:    []string{"5", "50", "100"},
			wantMode:  "MIXED",
			wantBand:  "10G",
			wantPower: "HIGH",
		},
		{
			name:      "power QRP",
			modes:     []string{"FT8", "SSB", "CW"},
			bands:     []string{"6m", "2m", "70cm"},
			powers:    []string{"5", "4", "1"},
			wantMode:  "MIXED",
			wantBand:  "ALL",
			wantPower: "QRP",
		},
		{
			name:      "power low",
			modes:     []string{"FT8", "SSB", "CW"},
			bands:     []string{"6m", "2m", "70cm"},
			powers:    []string{"5", "42", "10"},
			wantMode:  "MIXED",
			wantBand:  "ALL",
			wantPower: "LOW",
		},
		{
			name:      "power all high",
			modes:     []string{"FT8", "SSB", "CW"},
			bands:     []string{"6m", "2m", "70cm"},
			powers:    []string{"500", "43", "100"},
			wantMode:  "MIXED",
			wantBand:  "ALL",
			wantPower: "HIGH",
		},
	}

	cab := &CabrilloIO{
		CreatedBy:          "cabrillo_test",
		Contest:            "TEST-CONTEST-ID",
		Email:              "ham@example.com",
		GridLocator:        "FN31",
		Name:               "Ham Operator",
		LowPowerMax:        42,
		QRPPowerMax:        5,
		MyExchangeField:    "STX_STRING",
		TheirExchangeField: "SRX_STRING"}
	calls := []string{"A1AA", "B2BB", "C3CC", "D4DD"}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.modes) != len(tc.bands) || len(tc.bands) != len(tc.powers) {
				t.Fatalf("%s has mismatched lengths", tc.name)
			}
			l := NewLogfile()
			for i := 0; i < len(tc.modes); i++ {
				l.AddRecord(NewRecord(
					Field{Name: "QSO_DATE", Value: "20201031"},
					Field{Name: "TIME_ON", Value: fmt.Sprintf("12%02d", i)},
					Field{Name: "STATION_CALLSIGN", Value: "W1AW"},
					Field{Name: "CALL", Value: calls[i]},
					Field{Name: "BAND", Value: tc.bands[i]},
					Field{Name: "MODE", Value: tc.modes[i]},
					Field{Name: "TX_PWR", Value: tc.powers[i]},
					Field{Name: "RST_SENT", Value: "59"},
					Field{Name: "RST_RCVD", Value: "59"},
					Field{Name: "SRX_STRING", Value: fmt.Sprintf("R%d", i)},
					Field{Name: "STX_STRING", Value: fmt.Sprintf("T%d", i)},
				))
			}
			out := &strings.Builder{}
			if err := cab.Write(l, out); err != nil {
				t.Fatalf("Write(%v) got error %v", l, err)
			}
			lines := strings.Split(out.String(), "\n")
			headers := map[string]string{"CATEGORY-BAND": tc.wantBand, "CATEGORY-MODE": tc.wantMode, "CATEGORY-POWER": tc.wantPower}
			for h, want := range headers {
				found := false
				for _, line := range lines {
					if strings.HasPrefix(line, h+":") {
						found = true
						got := strings.TrimPrefix(line, h+": ")
						if got != want {
							t.Errorf("%s got %q, want %q", h, got, want)
						}
					}
				}
				if !found {
					t.Errorf("%s not in Cabrillo output, want %s", h, want)
				}
			}
		})
	}
}

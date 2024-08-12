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
	"github.com/flwyd/adif-multitool/adif/spec"
)

func TestInfer(t *testing.T) {
	withOtaFields := func(fs ...adif.Field) []adif.Field {
		res := []adif.Field{
			{Name: "IOTA", Value: "NA-001"}, {Name: "IOTA_ISLAND_ID", Value: "99999999"},
			{Name: "MY_IOTA", Value: "OC-047"}, {Name: "MY_IOTA_ISLAND_ID", Value: "8"},
			{Name: "POTA_REF", Value: "8P-0013"}, {Name: "MY_POTA_REF", Value: "VE-3157,K-0028"},
			{Name: "SOTA_REF", Value: "CE3/CO-001"}, {Name: "MY_SOTA_REF", Value: "S7/SC-002"},
			{Name: "WWFF_REF", Value: "3DFF-002"}, {Name: "MY_WWFF_REF", Value: "KFF-1234"},
		}
		for _, f := range fs {
			res = append(res, f)
		}
		return res
	}

	adi := adif.NewADIIO()
	tests := []struct {
		name        string
		infer       FieldList
		start, want []adif.Field
	}{
		{
			name:  "band 2m",
			infer: FieldList{"BAND"},
			start: []adif.Field{{Name: "FREQ", Value: "146.52"}},
			want:  []adif.Field{{Name: "FREQ", Value: "146.52"}, {Name: "BAND", Value: "2m"}},
		},
		{
			name:  "band 20m",
			infer: FieldList{"band"},
			start: []adif.Field{{Name: "freq", Value: "14.25600"}},
			want:  []adif.Field{{Name: "FREQ", Value: "14.25600"}, {Name: "BAND", Value: "20m"}},
		},
		{
			name:  "band no overwrite",
			infer: FieldList{"BAND"},
			start: []adif.Field{{Name: "FREQ", Value: "146.52"}, {Name: "BAND", Value: "rubber"}},
			want:  []adif.Field{{Name: "FREQ", Value: "146.52"}, {Name: "BAND", Value: "rubber"}},
		},
		{
			name:  "band overwrite empty",
			infer: FieldList{"BAND"},
			start: []adif.Field{{Name: "FREQ", Value: "222.222"}, {Name: "BAND", Value: ""}},
			want:  []adif.Field{{Name: "FREQ", Value: "222.222"}, {Name: "BAND", Value: "1.25m"}},
		},
		{
			name:  "lowest band_rx",
			infer: FieldList{"BAND_RX"},
			start: []adif.Field{{Name: "FREQ_RX", Value: ".1378"}},
			want:  []adif.Field{{Name: "FREQ_RX", Value: ".1378"}, {Name: "BAND_RX", Value: "2190m"}},
		},
		{
			name:  "highest band_rx",
			infer: FieldList{"BaNd_rX"},
			start: []adif.Field{{Name: "FREQ_RX", Value: "300000"}},
			want:  []adif.Field{{Name: "FREQ_RX", Value: "300000"}, {Name: "BAND_RX", Value: "submm"}},
		},

		{
			name:  "country not inferred from DXCC 0",
			infer: FieldList{"COUNTRY"},
			start: []adif.Field{{Name: "DXCC", Value: "0"}},
			want:  []adif.Field{{Name: "DXCC", Value: "0"}},
		},
		{
			name:  "country Canada",
			infer: FieldList{"country"},
			start: []adif.Field{{Name: "DXcc", Value: "1"}},
			want:  []adif.Field{{Name: "DXCC", Value: "1"}, {Name: "COUNTRY", Value: "CANADA"}},
		},
		{
			name:  "country Asiatic Russia",
			infer: FieldList{"Country"},
			start: []adif.Field{{Name: "DXCC", Value: "15"}},
			want:  []adif.Field{{Name: "DXCC", Value: "15"}, {Name: "COUNTRY", Value: "ASIATIC RUSSIA"}},
		},
		{
			name:  "dxcc unknown",
			infer: FieldList{"DXCC"},
			start: []adif.Field{{Name: "COUNTRY", Value: "AQUIESTAN"}},
			want:  []adif.Field{{Name: "COUNTRY", Value: "AQUIESTAN"}},
		},
		{
			name:  "dxcc Palestine (not the deleted one)",
			infer: FieldList{"DXCC"},
			start: []adif.Field{{Name: "COUNTRY", Value: "Palestine"}},
			want:  []adif.Field{{Name: "COUNTRY", Value: "Palestine"}, {Name: "DXCC", Value: "510"}},
		},
		{
			name:  "my_dxcc Kosovo",
			infer: FieldList{"MY_DXCC"},
			start: []adif.Field{{Name: "COUNTRY", Value: "Canada"}, {Name: "MY_COUNTRY", Value: "Republic of Kosovo"}},
			want:  []adif.Field{{Name: "COUNTRY", Value: "Canada"}, {Name: "MY_COUNTRY", Value: "Republic of Kosovo"}, {Name: "MY_DXCC", Value: "522"}},
		},

		{
			name:  "mode MFSK",
			infer: FieldList{"MODE"},
			start: []adif.Field{{Name: "SUBMODE", Value: "JS8"}},
			want:  []adif.Field{{Name: "SUBMODE", Value: "JS8"}, {Name: "MODE", Value: "MFSK"}},
		},
		{
			name:  "mode VARA",
			infer: FieldList{"mode"},
			start: []adif.Field{{Name: "submode", Value: "vara fm 1200"}},
			want:  []adif.Field{{Name: "SUBMODE", Value: "vara fm 1200"}, {Name: "MODE", Value: "DYNAMIC"}},
		},
		{
			name:  "mode does not overwrite",
			infer: FieldList{"MODE"},
			start: []adif.Field{{Name: "SUBMODE", Value: "USB"}, {Name: "MODE", Value: "FM"}},
			want:  []adif.Field{{Name: "SUBMODE", Value: "USB"}, {Name: "MODE", Value: "FM"}},
		},

		{
			name:  "gridsquare null island",
			infer: FieldList{"gridsquare"},
			start: []adif.Field{{Name: "LAT", Value: "N000 00.000"}, {Name: "LON", Value: "E000 00.000"}},
			want:  []adif.Field{{Name: "LAT", Value: "N000 00.000"}, {Name: "LON", Value: "E000 00.000"}, {Name: "GRIDSQUARE", Value: "JJ00aa00"}},
		},
		{
			name:  "gridsquare and ext",
			infer: FieldList{"gridsquare", "gridsquare_ext"},
			start: []adif.Field{{Name: "LAT", Value: "S000 00.001"}, {Name: "LON", Value: "W000 00.001"}},
			want:  []adif.Field{{Name: "LAT", Value: "S000 00.001"}, {Name: "LON", Value: "W000 00.001"}, {Name: "GRIDSQUARE", Value: "II99xx99"}, {Name: "GRIDSQUARE_EXT", Value: "xx99"}},
		},
		{
			name:  "gridsquare_ext alone",
			infer: FieldList{"gridsquare_ext"},
			start: []adif.Field{{Name: "LAT", Value: "N042 59.995"}, {Name: "LON", Value: "W179 59.970"}},
			want:  []adif.Field{{Name: "LAT", Value: "N042 59.995"}, {Name: "LON", Value: "W179 59.970"}, {Name: "GRIDSQUARE_EXT", Value: "bx45"}},
		},
		{
			name:  "my_gridsquare South Atlantic",
			infer: FieldList{"MY_GRIDSQUARE"},
			start: []adif.Field{{Name: "MY_LAT", Value: "S045 30.001"}, {Name: "MY_LON", Value: "W045 30.001"}},
			want:  []adif.Field{{Name: "MY_LAT", Value: "S045 30.001"}, {Name: "MY_LON", Value: "W045 30.001"}, {Name: "MY_GRIDSQUARE", Value: "GE74fl99"}},
		},
		{
			name:  "my_gridsquare and ext north pole",
			infer: FieldList{"MY_GRIDSQUARE", "MY_GRIDSQUARE_EXT"},
			start: []adif.Field{{Name: "MY_LAT", Value: "N089 59.999"}, {Name: "MY_LON", Value: "E179 59.999"}},
			want:  []adif.Field{{Name: "MY_LAT", Value: "N089 59.999"}, {Name: "MY_LON", Value: "E179 59.999"}, {Name: "MY_GRIDSQUARE", Value: "RR99xx99"}, {Name: "MY_GRIDSQUARE_EXT", Value: "xx99"}},
		},
		{
			name:  "lat lon null island",
			infer: FieldList{"lat", "lon"},
			start: []adif.Field{{Name: "GRIDSQUARE", Value: "JJ00aa00"}, {Name: "GRIDSQUARE_EXT", Value: "aa00"}},
			want:  []adif.Field{{Name: "GRIDSQUARE", Value: "JJ00aa00"}, {Name: "GRIDSQUARE_EXT", Value: "aa00"}, {Name: "LAT", Value: "N000 00.001"}, {Name: "LON", Value: "E000 00.001"}},
		},
		{
			name:  "my_lat my_lon no ext",
			infer: FieldList{"my_lat", "my_lon"},
			start: []adif.Field{{Name: "MY_GRIDSQUARE", Value: "AB23cd45"}},
			want:  []adif.Field{{Name: "MY_GRIDSQUARE", Value: "AB23cd45"}, {Name: "MY_LAT", Value: "S076 51.125"}, {Name: "MY_LON", Value: "W175 47.750"}},
		},

		{
			name:  "operator from guest_op",
			infer: FieldList{"OPERATOR"},
			start: []adif.Field{{Name: "GUEST_OP", Value: "W1AW"}},
			want:  []adif.Field{{Name: "GUEST_OP", Value: "W1AW"}, {Name: "OPERATOR", Value: "W1AW"}},
		},
		{
			name:  "station_callsign from operator",
			infer: FieldList{"OWNER_CALLSIGN", "STATION_CALLSIGN"},
			start: []adif.Field{{Name: "OPERATOR", Value: "N0P"}, {Name: "OWNER_CALLSIGN", Value: "W1AW"}},
			want:  []adif.Field{{Name: "OPERATOR", Value: "N0P"}, {Name: "OWNER_CALLSIGN", Value: "W1AW"}, {Name: "STATION_CALLSIGN", Value: "N0P"}},
		},
		{
			name:  "owner_callsign and station_callsign from operator",
			infer: FieldList{"OWNER_CALLSIGN", "STATION_CALLSIGN"},
			start: []adif.Field{{Name: "OPERATOR", Value: "W1AW"}},
			want:  []adif.Field{{Name: "OPERATOR", Value: "W1AW"}, {Name: "OWNER_CALLSIGN", Value: "W1AW"}, {Name: "STATION_CALLSIGN", Value: "W1AW"}},
		},

		{
			name:  "usaca_counties from US cnty",
			infer: FieldList{"USACA_COUNTIES", "MY_USACA_COUNTIES"},
			start: []adif.Field{
				{Name: "CNTY", Value: "MD,St. Mary's"}, {Name: "DXCC", Value: spec.CountryUnitedStatesOfAmerica.EntityCode},
				{Name: "MY_CNTY", Value: "HI,Maui"}, {Name: "MY_DXCC", Value: spec.CountryHawaii.EntityCode}},
			want: []adif.Field{
				{Name: "USACA_COUNTIES", Value: "MD,St. Mary's"}, {Name: "CNTY", Value: "MD,St. Mary's"}, {Name: "DXCC", Value: spec.CountryUnitedStatesOfAmerica.EntityCode},
				{Name: "MY_USACA_COUNTIES", Value: "HI,Maui"}, {Name: "MY_CNTY", Value: "HI,Maui"}, {Name: "MY_DXCC", Value: spec.CountryHawaii.EntityCode}},
		},
		{
			name:  "cnty from usaca_counties",
			infer: FieldList{"CNTY", "MY_CNTY"},
			start: []adif.Field{
				{Name: "USACA_COUNTIES", Value: "AK,Prince of Wales-Outer Ketchikan"}, {Name: "DXCC", Value: spec.CountryAlaska.EntityCode},
				{Name: "MY_USACA_COUNTIES", Value: "FL,Miami-Dade"}, {Name: "MY_DXCC", Value: spec.CountryUnitedStatesOfAmerica.EntityCode}},
			want: []adif.Field{
				{Name: "CNTY", Value: "AK,Prince of Wales-Outer Ketchikan"}, {Name: "USACA_COUNTIES", Value: "AK,Prince of Wales-Outer Ketchikan"}, {Name: "DXCC", Value: spec.CountryAlaska.EntityCode},
				{Name: "MY_CNTY", Value: "FL,Miami-Dade"}, {Name: "MY_USACA_COUNTIES", Value: "FL,Miami-Dade"}, {Name: "MY_DXCC", Value: spec.CountryUnitedStatesOfAmerica.EntityCode}},
		},
		{
			name:  "multiple usaca_counties not copied to cnty",
			infer: FieldList{"CNTY", "MY_CNTY"},
			start: []adif.Field{
				{Name: "USACA_COUNTIES", Value: "MA,Franklin:MA,Hampshire"}, {Name: "DXCC", Value: spec.CountryUnitedStatesOfAmerica.EntityCode},
				{Name: "MY_USACA_COUNTIES", Value: "AK,Northwest Arctic:AK,Nome"}, {Name: "MY_DXCC", Value: spec.CountryAlaska.EntityCode}},
			want: []adif.Field{
				{Name: "USACA_COUNTIES", Value: "MA,Franklin:MA,Hampshire"}, {Name: "DXCC", Value: spec.CountryUnitedStatesOfAmerica.EntityCode},
				{Name: "MY_USACA_COUNTIES", Value: "AK,Northwest Arctic:AK,Nome"}, {Name: "MY_DXCC", Value: spec.CountryAlaska.EntityCode}},
		},
		{
			name:  "JA subdivisions not copied to usaca_counties",
			infer: FieldList{"USACA_COUNTIES", "MY_USACA_COUNTIES"},
			start: []adif.Field{
				{Name: "CNTY", Value: "01001"}, {Name: "DXCC", Value: spec.CountryJapan.EntityCode},
				{Name: "MY_CNTY", Value: "100100"}, {Name: "MY_DXCC", Value: spec.CountryJapan.EntityCode}},
			want: []adif.Field{
				{Name: "CNTY", Value: "01001"}, {Name: "DXCC", Value: spec.CountryJapan.EntityCode},
				{Name: "MY_CNTY", Value: "100100"}, {Name: "MY_DXCC", Value: spec.CountryJapan.EntityCode}},
		},
		{
			name:  "non-US subdivisions matching format not copied to usaca_counties",
			infer: FieldList{"USACA_COUNTIES", "MY_USACA_COUNTIES"},
			start: []adif.Field{
				{Name: "CNTY", Value: "NL,Monterrey"}, {Name: "DXCC", Value: spec.CountryMexico.EntityCode},
				{Name: "MY_CNTY", Value: "GD,Guangzhou"}, {Name: "MY_DXCC", Value: spec.CountryChina.EntityCode}},
			want: []adif.Field{
				{Name: "CNTY", Value: "NL,Monterrey"}, {Name: "DXCC", Value: spec.CountryMexico.EntityCode},
				{Name: "MY_CNTY", Value: "GD,Guangzhou"}, {Name: "MY_DXCC", Value: spec.CountryChina.EntityCode}},
		},
		{
			name:  "usaca_counties does not overwrite",
			infer: FieldList{"CNTY", "MY_CNTY", "USACA_COUNTIES", "MY_USACA_COUNTIES"},
			start: []adif.Field{
				{Name: "CNTY", Value: "AL,Autauga"}, {Name: "USACA_COUNTIES", Value: "AL,Winston"}, {Name: "DXCC", Value: spec.CountryUnitedStatesOfAmerica.EntityCode},
				{Name: "MY_CNTY", Value: "WY,Albany"}, {Name: "MY_USACA_COUNTIES", Value: "WY,Weston"}, {Name: "MY_DXCC", Value: spec.CountryUnitedStatesOfAmerica.EntityCode},
			},
			want: []adif.Field{
				{Name: "CNTY", Value: "AL,Autauga"}, {Name: "USACA_COUNTIES", Value: "AL,Winston"}, {Name: "DXCC", Value: spec.CountryUnitedStatesOfAmerica.EntityCode},
				{Name: "MY_CNTY", Value: "WY,Albany"}, {Name: "MY_USACA_COUNTIES", Value: "WY,Weston"}, {Name: "MY_DXCC", Value: spec.CountryUnitedStatesOfAmerica.EntityCode},
			},
		},

		{
			name:  "sig_info nothing when ambiguous",
			infer: FieldList{"SIG_INFO", "MY_SIG_INFO"},
			start: withOtaFields(),
			want:  withOtaFields(),
		},
		{
			name:  "sig_info with sig=IOTA",
			infer: FieldList{"SIG_INFO"},
			start: withOtaFields(adif.Field{Name: "SIG", Value: "IOTA"}, adif.Field{Name: "MY_SIG", Value: "POTA"}),
			want:  withOtaFields(adif.Field{Name: "SIG", Value: "IOTA"}, adif.Field{Name: "SIG_INFO", Value: "NA-001"}, adif.Field{Name: "MY_SIG", Value: "POTA"}),
		},
		{
			name:  "my_sig_info with my_sig=IOTA",
			infer: FieldList{"my_SIG_INFO"},
			start: withOtaFields(adif.Field{Name: "my_siG", Value: "IOTA"}, adif.Field{Name: "SIG", Value: "WWFF"}),
			want:  withOtaFields(adif.Field{Name: "MY_SIG", Value: "IOTA"}, adif.Field{Name: "MY_SIG_INFO", Value: "OC-047"}, adif.Field{Name: "SIG", Value: "WWFF"}),
		},
		{
			name:  "sig_info with sig=POTA",
			infer: FieldList{"sig_info"},
			start: withOtaFields(adif.Field{Name: "SIG", Value: "POTA"}, adif.Field{Name: "MY_SIG", Value: "IOTA"}),
			want:  withOtaFields(adif.Field{Name: "SIG", Value: "POTA"}, adif.Field{Name: "SIG_INFO", Value: "8P-0013"}, adif.Field{Name: "MY_SIG", Value: "IOTA"}),
		},
		{
			name:  "my_sig_info with my_sig=POTA",
			infer: FieldList{"MY_SIG_INFO"},
			start: withOtaFields(adif.Field{Name: "MY_SIG", Value: "POTA"}, adif.Field{Name: "SIG", Value: "SOTA"}),
			want:  withOtaFields(adif.Field{Name: "MY_SIG", Value: "POTA"}, adif.Field{Name: "MY_SIG_INFO", Value: "VE-3157,K-0028"}, adif.Field{Name: "SIG", Value: "SOTA"}),
		},
		{
			name:  "sig_info with sig=SOTA",
			infer: FieldList{"SIG_INFO"},
			start: withOtaFields(adif.Field{Name: "SIG", Value: "SOTA"}, adif.Field{Name: "MY_SIG", Value: "WWFF"}),
			want:  withOtaFields(adif.Field{Name: "SIG", Value: "SOTA"}, adif.Field{Name: "SIG_INFO", Value: "CE3/CO-001"}, adif.Field{Name: "MY_SIG", Value: "WWFF"}),
		},
		{
			name:  "my_sig_info with my_sig=SOTA",
			infer: FieldList{"my_SIG_info"},
			start: withOtaFields(adif.Field{Name: "my_SIG", Value: "SOTA"}, adif.Field{Name: "SIG", Value: "POTA"}),
			want:  withOtaFields(adif.Field{Name: "MY_SIG", Value: "SOTA"}, adif.Field{Name: "MY_SIG_INFO", Value: "S7/SC-002"}, adif.Field{Name: "SIG", Value: "POTA"}),
		},
		{
			name:  "sig_info with sig=WWFF",
			infer: FieldList{"SIG_INFO"},
			start: withOtaFields(adif.Field{Name: "SIG", Value: "WWFF"}, adif.Field{Name: "MY_SIG", Value: "SOTA"}),
			want:  withOtaFields(adif.Field{Name: "SIG", Value: "WWFF"}, adif.Field{Name: "SIG_INFO", Value: "3DFF-002"}, adif.Field{Name: "MY_SIG", Value: "SOTA"}),
		},
		{
			name:  "my_sig_info with my_sig=WWFF",
			infer: FieldList{"MY_SIG_INFO"},
			start: withOtaFields(adif.Field{Name: "MY_SIG", Value: "WWFF"}, adif.Field{Name: "SIG", Value: "IOTA"}),
			want:  withOtaFields(adif.Field{Name: "MY_SIG", Value: "WWFF"}, adif.Field{Name: "MY_SIG_INFO", Value: "KFF-1234"}, adif.Field{Name: "SIG", Value: "IOTA"}),
		},
		{
			name:  "sig_info and sig from IOTA",
			infer: FieldList{"SIG_INFO"},
			start: []adif.Field{{Name: "IOTA", Value: "AF-123"}, {Name: "IOTA_ISLAND_ID", Value: "321"}},
			want:  []adif.Field{{Name: "IOTA", Value: "AF-123"}, {Name: "IOTA_ISLAND_ID", Value: "321"}, {Name: "SIG_INFO", Value: "AF-123"}, {Name: "SIG", Value: "IOTA"}},
		},
		{
			name:  "my_sig_info and my_sig from IOTA",
			infer: FieldList{"MY_SIG_INFO"},
			start: []adif.Field{{Name: "MY_IOTA", Value: "AF-123"}, {Name: "MY_IOTA_ISLAND_ID", Value: "321"}},
			want:  []adif.Field{{Name: "MY_IOTA", Value: "AF-123"}, {Name: "MY_IOTA_ISLAND_ID", Value: "321"}, {Name: "MY_SIG_INFO", Value: "AF-123"}, {Name: "MY_SIG", Value: "IOTA"}},
		},
		{
			name:  "sig_info and sig from POTA",
			infer: FieldList{"SIG_INFO"},
			start: []adif.Field{{Name: "POTA_REF", Value: "C9-0007"}},
			want:  []adif.Field{{Name: "POTA_REF", Value: "C9-0007"}, {Name: "SIG_INFO", Value: "C9-0007"}, {Name: "SIG", Value: "POTA"}},
		},
		{
			name:  "my_sig_info and my_sig from POTA",
			infer: FieldList{"MY_SIG_INFO"},
			start: []adif.Field{{Name: "MY_POTA_REF", Value: "C9-0007"}},
			want:  []adif.Field{{Name: "MY_POTA_REF", Value: "C9-0007"}, {Name: "MY_SIG_INFO", Value: "C9-0007"}, {Name: "MY_SIG", Value: "POTA"}},
		},
		{
			name:  "sig_info and sig from SOTA",
			infer: FieldList{"SIG_INFO"},
			start: []adif.Field{{Name: "SOTA_REF", Value: "UT/CA-109"}},
			want:  []adif.Field{{Name: "SOTA_REF", Value: "UT/CA-109"}, {Name: "SIG_INFO", Value: "UT/CA-109"}, {Name: "SIG", Value: "SOTA"}},
		},
		{
			name:  "my_sig_info and my_sig from SOTA",
			infer: FieldList{"MY_SIG_INFO"},
			start: []adif.Field{{Name: "MY_SOTA_REF", Value: "UT/CA-109"}},
			want:  []adif.Field{{Name: "MY_SOTA_REF", Value: "UT/CA-109"}, {Name: "MY_SIG_INFO", Value: "UT/CA-109"}, {Name: "MY_SIG", Value: "SOTA"}},
		},
		{
			name:  "sig_info and sig from WWFF",
			infer: FieldList{"SIG_INFO"},
			start: []adif.Field{{Name: "WWFF_REF", Value: "P29FF-123"}},
			want:  []adif.Field{{Name: "WWFF_REF", Value: "P29FF-123"}, {Name: "SIG_INFO", Value: "P29FF-123"}, {Name: "SIG", Value: "WWFF"}},
		},
		{
			name:  "my_sig_info and my_sig from WWFF",
			infer: FieldList{"MY_SIG_INFO"},
			start: []adif.Field{{Name: "MY_WWFF_REF", Value: "P29FF-123"}},
			want:  []adif.Field{{Name: "MY_WWFF_REF", Value: "P29FF-123"}, {Name: "MY_SIG_INFO", Value: "P29FF-123"}, {Name: "MY_SIG", Value: "WWFF"}},
		},

		{
			name: "no inference if missing",
			infer: FieldList{
				"BAND", "BAND_RX", "MODE",
				"DXCC", "MY_DXCC", "COUNTRY", "MY_COUNTRY",
				"GRIDSQUARE", "GRIDSQUARE_EXT", "MY_GRIDSQUARE", "MY_GRIDSQUARE_EXT",
				"LAT", "LON", "MY_LAT", "MY_LON",
				"OPERATOR", "STATION_CALLSIGN", "OWNER_CALLSIGN",
				"SIG_INFO", "IOTA", "POTA_REF", "SOTA_REF", "WWFF_REF",
				"MY_SIG_INFO", "MY_IOTA", "MY_POTA_REF", "MY_SOTA_REF", "MY_WWFF_REF",
			},
			start: []adif.Field{{Name: "CALL", Value: "W1AW"}, {Name: "RST_RCVD", Value: "59"}},
			want:  []adif.Field{{Name: "CALL", Value: "W1AW"}, {Name: "RST_RCVD", Value: "59"}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			in := &bytes.Buffer{}
			lin := adif.NewLogfile()
			lin.AddRecord(adif.NewRecord(tc.start...))
			if err := adi.Write(lin, in); err != nil {
				t.Fatalf("Error writing fields %v: %v", tc.start, err)
			}
			out := &bytes.Buffer{}
			ctx := &Context{
				InputFormat:  adif.FormatADI,
				OutputFormat: adif.FormatADI,
				Readers:      readers(adi),
				Writers:      writers(adi),
				Out:          out,
				fs:           fakeFilesystem{map[string]string{"foo.adi": in.String()}},
				CommandCtx:   &InferContext{Fields: tc.infer}}
			if err := Infer.Run(ctx, []string{"foo.adi"}); err != nil {
				t.Fatalf("Infer(%s) got error %v", in.String(), err)
			}
			l, err := adi.Read(out)
			if err != nil {
				t.Fatalf("Read(%s) got error: %v", out.String(), err)
			}
			if len(l.Records) != 1 {
				t.Errorf("Read(%s) got %d records, want 1", out.String(), len(l.Records))
			} else {
				got := l.Records[0]
				want := adif.NewRecord(tc.want...)
				if !want.Equal(l.Records[0]) {
					t.Errorf("infer %v from %v got %v, want %v", tc.infer, tc.start, got, want)
				}
			}
		})
	}
}

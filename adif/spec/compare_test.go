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

package spec

import (
	"testing"

	"golang.org/x/text/language"
)

type comparisonCheck func(t *testing.T, c FieldComparator, values ...string)

func shouldCompareEqual(t *testing.T, c FieldComparator, values ...string) {
	t.Helper()
	for i := 0; i < len(values); i++ {
		a := values[i]
		for j := 0; j < len(values); j++ {
			b := values[j]
			if got, err := c(a, b); err != nil {
				t.Errorf("compare(%q, %q) got error %v", a, b, err)
			} else if got != 0 {
				t.Errorf("compare(%q, %q) got %d, want 0", a, b, got)
			}
		}
	}
}

func shouldCompareLess(t *testing.T, c FieldComparator, values ...string) {
	t.Helper()
	for i := 0; i < len(values); i++ {
		a := values[i]
		if gotself, err := c(a, a); err != nil {
			t.Errorf("compare(%q, %q) got error %v", a, a, err)
		} else if gotself != 0 {
			t.Errorf("compare(%q, %q) got %d, want 0", a, a, gotself)
		}
		for j := i + 1; j < len(values); j++ {
			b := values[j]
			if got, err := c(a, b); err != nil {
				t.Errorf("compare(%q, %q) got error %v", a, b, err)
			} else if got >= 0 {
				t.Errorf("compare(%q, %q) got %d, want < 0", a, b, got)
			}
			if got, err := c(b, a); err != nil {
				t.Errorf("compare(%q, %q) got error %v", b, a, err)
			} else if got <= 0 {
				t.Errorf("compare(%q, %q) got %d, want > 0", b, a, got)
			}
		}
	}
}

func TestCompareStringsBasic(t *testing.T) {
	tests := []struct {
		name  string
		field Field
		want  comparisonCheck
		vals  []string
	}{
		{
			name: "empty string equal", field: ProgramidField, want: shouldCompareEqual,
			vals: []string{"", ""},
		},
		{
			name: "empty string less", field: CheckField, want: shouldCompareLess,
			vals: []string{"", "Apple", "Banana"},
		},
		{
			name: "identical strings", field: CallField, want: shouldCompareEqual,
			vals: []string{"W1AW", "W1AW"},
		},
		{
			name: "different case", field: MyAntennaField, want: shouldCompareEqual,
			vals: []string{"DIPOLE", "dipole", "dIpOlE", "dipolE", "DIPole"},
		},
		{
			name: "diacritics are equal", field: NameField, want: shouldCompareEqual,
			vals: []string{"Jean Michel", "Jéàn Mîçhēl", "JÉÀÑ MÎÇHĒŁ"},
		},
		{
			name: "gridsquare mixed case", field: GridsquareField, want: shouldCompareEqual,
			vals: []string{"AA00XX99", "AA00xx99", "aa00xx99", "aa00XX99", "aA00Xx99"},
		},
		{
			name: "ASCII number order", field: SrxStringField, want: shouldCompareLess,
			vals: []string{"1", "123", "1234", "2", "234", "32", "9"},
		},
		{
			name: "preserves leading zero and space", field: RigField, want: shouldCompareLess,
			vals: []string{"MyHf 00123", "MyHF 0123", "MyHF 123", "MyHF0123", "MyHF123"},
		},
		{
			name: "respects punctuation", field: NameField, want: shouldCompareLess,
			vals: []string{"Peggy Sue", "peggy-sue", "Peggy, Sue", "Peggy/Sue"},
		},
		{
			name: "articles don't change order", field: NotesField, want: shouldCompareLess,
			vals: []string{"A good chat", "Chatted for a bit", "The best chat"},
		},
		{
			name: "gridsquare sort", field: GridsquareField, want: shouldCompareLess,
			vals: []string{"AA00aa00", "AA00aa01", "AA00ab00", "AA09aa00", "AA10aa00", "AA11jk00", "AB00CD00", "dn00ab", "dn07wx", "dn70jk", "RR99xx", "RR99xx99"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := ComparatorForField(tc.field, language.Greek)
			tc.want(t, c, tc.vals...)
		})
	}
}

func TestCompareStringsLocale(t *testing.T) {
	tests := []struct {
		name  string
		field Field
		lang  language.Tag
		want  comparisonCheck
		vals  []string
	}{
		{
			name: "empty string equal", field: CountryIntlField, want: shouldCompareEqual,
			vals: []string{"", ""},
		},
		{
			name: "empty string less", field: CountryIntlField, lang: language.English, want: shouldCompareLess,
			vals: []string{"", "Costa Rica", "Côte d'Ivoire", "Ελλάδα", "中国"},
		},
		{
			name: "identical strings", field: NameIntlField, lang: language.French, want: shouldCompareEqual,
			vals: []string{"Jean Michel", "Jean Michel"},
		},
		{
			name: "French sort", field: NameIntlField, lang: language.French, want: shouldCompareLess,
			vals: []string{"Jean", "Jean Michel", "Jean Michél", "Jean Mîchel", "Jèan Michel", "Jèan Michel ⚜"},
		},
		{
			name: "German sort", field: MyStreetIntlField, lang: language.German, want: shouldCompareLess,
			vals: []string{"Äcker", "Apfel", "Öden", "Ost", "Strasse", "Straße", "Über", "Ulm"},
		},
		{
			// Go bug: Norwegian doesn't sort like Danish https://github.com/golang/go/issues/59908
			name: "Norwegian sort", field: QthIntlField, lang: language.Danish, want: shouldCompareLess,
			vals: []string{"Arendal", "Bergen", "Oslo", "Trondheim", "Ænes", "Østfold", "Ålgård"},
		},
		{
			name: "English order of Norwegian names", field: QthIntlField, lang: language.English, want: shouldCompareLess,
			vals: []string{"Ænes", "Ålgård", "Arendal", "Bergen", "Oslo", "Østfold", "Trondheim"},
		},
		{
			name: "full width equal", field: CommentIntlField, lang: language.Japanese, want: shouldCompareEqual,
			vals: []string{"FOO, bar? (ｯ)", "ＦＯＯ，　ｂａｒ？　（ツ）", "ｆｏｏ，　ＢＡＲ？　（ツ）"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := ComparatorForField(tc.field, tc.lang)
			tc.want(t, c, tc.vals...)
		})
	}
}

func TestCompareNumbers(t *testing.T) {
	tests := []struct {
		name  string
		field Field
		want  comparisonCheck
		vals  []string
	}{
		{
			name: "empty string equal", field: RxPwrField, want: shouldCompareEqual,
			vals: []string{"", ""},
		},
		{
			name: "empty string less", field: MyAltitudeField, want: shouldCompareLess,
			vals: []string{"", "-50", "0", "3"},
		},
		{
			name: "identical strings", field: FreqField, want: shouldCompareEqual,
			vals: []string{"14.300", "14.300"},
		},
		{
			name: "trailing zeroes", field: FreqField, want: shouldCompareEqual,
			vals: []string{"14.300", "14.30", "14.3", "14.30000"},
		},
		{
			name: "leading zeroes", field: SrxField, want: shouldCompareEqual,
			vals: []string{"123", "0123", "00123", "0123.00"},
		},
		{
			name: "ascending integers", field: TxPwrField, want: shouldCompareLess,
			vals: []string{"1", "2", "9", "32", "098", "123", "234", "1234"},
		},
		{
			name: "ascending decimals", field: AltitudeField, want: shouldCompareLess,
			vals: []string{"-2.6", "-2", "-1.0", "-0.5", "0.3", "1.0", "1.00001", "1.5", "2", "3.1", "100.9"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := ComparatorForField(tc.field, language.Thai)
			tc.want(t, c, tc.vals...)
		})
	}
}

func TestCompareDates(t *testing.T) {
	tests := []struct {
		name  string
		field Field
		want  comparisonCheck
		vals  []string
	}{
		{
			name: "empty string equal", field: QsoDateField, want: shouldCompareEqual,
			vals: []string{"", ""},
		},
		{
			name: "empty string less", field: QsoDateField, want: shouldCompareLess,
			vals: []string{"", "19650403", "19991231", "20130719"},
		},
		{
			name: "identical strings", field: QsoDateOffField, want: shouldCompareEqual,
			vals: []string{"20210521", "20210521"},
		},
		{
			name: "year month day", field: QsoDateOffField, want: shouldCompareLess,
			vals: []string{"20090807", "20120101"},
		},
		{
			name: "same year", field: QsoDateField, want: shouldCompareLess,
			vals: []string{"19870101", "19870102", "19870228", "19870321", "19870401", "19870916", "19871031", "19871206"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := ComparatorForField(tc.field, language.Arabic)
			tc.want(t, c, tc.vals...)
		})
	}
}

func TestCompareTimes(t *testing.T) {
	tests := []struct {
		name  string
		field Field
		want  comparisonCheck
		vals  []string
	}{
		{
			name: "empty string equal", field: TimeOnField, want: shouldCompareEqual,
			vals: []string{"", ""},
		},
		{
			name: "empty string less", field: TimeOnField, want: shouldCompareLess,
			vals: []string{"", "0000", "0123", "123456"},
		},
		{
			name: "identical strings", field: TimeOffField, want: shouldCompareEqual,
			vals: []string{"123456", "123456"},
		},
		{
			name: "zero seconds", field: TimeOffField, want: shouldCompareEqual,
			vals: []string{"1234", "123400"},
		},
		{
			name: "leading zeroes", field: TimeOnField, want: shouldCompareLess,
			vals: []string{"001234", "0123", "012304"},
		},
		{
			name: "mixed hhmm and hhmmss", field: TimeOffField, want: shouldCompareLess,
			vals: []string{"0123", "012345", "024530", "0315", "1234", "123456", "1235", "141530", "1416"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := ComparatorForField(tc.field, language.Arabic)
			tc.want(t, c, tc.vals...)
		})
	}
}

func TestCompareBooleans(t *testing.T) {
	tests := []struct {
		name  string
		field Field
		want  comparisonCheck
		vals  []string
	}{
		{
			name: "empty string equal", field: ForceInitField, want: shouldCompareEqual,
			vals: []string{"", ""},
		},
		{
			name: "empty string less", field: ForceInitField, want: shouldCompareLess,
			vals: []string{"", "N", "Y"},
		},
		{
			name: "identical Y strings", field: SilentKeyField, want: shouldCompareEqual,
			vals: []string{"Y", "Y"},
		},
		{
			name: "identical N strings", field: SilentKeyField, want: shouldCompareEqual,
			vals: []string{"N", "N"},
		},
		{
			name: "case-insensitive Y", field: QsoRandomField, want: shouldCompareEqual,
			vals: []string{"y", "Y"},
		},
		{
			name: "case-insensitive N", field: QsoRandomField, want: shouldCompareEqual,
			vals: []string{"n", "N"},
		},
		{
			name: "false then true", field: QsoRandomField, want: shouldCompareLess,
			vals: []string{"n", "Y"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := ComparatorForField(tc.field, language.Tamil)
			tc.want(t, c, tc.vals...)
		})
	}
}

func TestCompareEnums(t *testing.T) {
	tests := []struct {
		name  string
		field Field
		want  comparisonCheck
		vals  []string
	}{
		{
			name: "empty string equal", field: ArrlSectField, want: shouldCompareEqual,
			vals: []string{"", ""},
		},
		{
			name: "empty string less", field: ContField, want: shouldCompareLess,
			vals: []string{"", "EU", "SA"},
		},
		{
			name: "identical strings", field: ModeField, want: shouldCompareEqual,
			vals: []string{"CW", "CW"},
		},
		{
			name: "different case", field: ArrlSectField, want: shouldCompareEqual,
			vals: []string{"PEI", "pei", "PeI", "pEi", "peI", "Pei", "PEi"},
		},
		{
			name: "alphabetic order", field: ArrlSectField, want: shouldCompareLess,
			vals: []string{"AK", "AL", "EPA", "ETX", "STX", "WPA", "WTX", "WY"},
		},
		// special-cased enums
		{
			name: "DXCC in numeric order", field: DxccField, want: shouldCompareLess,
			vals: []string{"1", "7", "13", "39", "117", "129", "160", "202", "222", "262", "309", "344", "382", "414", "499", "517", "522"},
		},
		{
			name: "band in frequency order", field: BandField, want: shouldCompareLess,
			vals: []string{"2190m", "630m", "560m", "160m", "80m", "40m", "30m", "20m", "17m", "15m", "12m", "10m", "8m", "6m", "5m", "4m", "2m", "1.25m", "70cm", "33cm", "23cm", "13cm", "9cm", "6cm", "3cm", "1.25cm", "6mm", "4mm", "2.5mm", "2mm", "1mm", "submm"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := ComparatorForField(tc.field, language.Greek)
			tc.want(t, c, tc.vals...)
		})
	}
}

func TestCompareLocations(t *testing.T) {
	tests := []struct {
		name  string
		field Field
		want  comparisonCheck
		vals  []string
	}{
		{
			name: "empty string equal", field: LonField, want: shouldCompareEqual,
			vals: []string{"", ""},
		},
		{
			name: "empty string less", field: LatField, want: shouldCompareLess,
			vals: []string{"", "N012 34.567"},
		},
		{
			name: "identical strings", field: MyLonField, want: shouldCompareEqual,
			vals: []string{"E123 45.678", "E123 45.678"},
		},
		{
			name: "different case", field: MyLatField, want: shouldCompareEqual,
			vals: []string{"S001 23.456", "s001 23.456"},
		},
		{
			name: "south ascending", field: LatField, want: shouldCompareLess,
			vals: []string{"S090 00.000", "S089 59.999", "s089 00.000", "S045 30.500", "S045 30.432", "S007 53.197", "S000 00.000"},
		},
		{
			name: "north ascending", field: LatField, want: shouldCompareLess,
			vals: []string{"N000 00.000", "N007 53.197", "n045 30.432", "N045 30.500", "N089 00.000", "N089 59.999", "N090 00.000"},
		},
		{
			name: "west ascending", field: LatField, want: shouldCompareLess,
			vals: []string{"W180 00.000", "W123 45.678", "W090 99.999", "W090 00.000", "W089 59.999", "w089 00.000", "W045 30.500", "W045 30.432", "W007 53.197", "W000 00.000"},
		},
		{
			name: "east ascending", field: LatField, want: shouldCompareLess,
			vals: []string{"E000 00.000", "E007 53.197", "E045 30.432", "E045 30.500", "e089 00.000", "E089 59.999", "E090 00.000", "E090 99.999", "E123 45.678", "E180 00.000"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := ComparatorForField(tc.field, language.Greek)
			tc.want(t, c, tc.vals...)
		})
	}
}

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
	"strings"
	"testing"
	"time"
)

type validateTest struct {
	field  Field
	value  string
	want   Validity
	others []validateTest
}

var emptyCtx = ValidationContext{FieldValue: func(name string) string { return "" }}

func validateTestCtx(v validateTest) ValidationContext {
	return ValidationContext{FieldValue: func(name string) string {
		for _, o := range v.others {
			if strings.EqualFold(name, o.field.Name) {
				return o.value
			}
		}
		return ""
	}}
}

func testValidator(t *testing.T, tc validateTest, ctx ValidationContext, funcname string) {
	t.Helper()
	v := TypeValidators[tc.field.Type.Name]
	if got := v(tc.value, tc.field, ctx); got.Validity != tc.want {
		if got.Validity == Valid {
			t.Errorf("%s(%q, %q, ctx) got Valid, want %s", funcname, tc.value, tc.field.Name, tc.want)
		} else {
			t.Errorf("%s(%q, %q, ctx) want %s got %s %s", funcname, tc.value, tc.field.Name, tc.want, got.Validity, got.Message)
		}
	}
}

func TestValidateBoolean(t *testing.T) {
	tests := []validateTest{
		{field: QsoRandomField, value: "Y", want: Valid},
		{field: SilentKeyField, value: "y", want: Valid},
		{field: ForceInitField, value: "N", want: Valid},
		{field: SwlField, value: "n", want: Valid},
		{field: QsoRandomField, value: "YES", want: InvalidError},
		{field: SilentKeyField, value: "true", want: InvalidError},
		{field: ForceInitField, value: "F", want: InvalidError},
		{field: SwlField, value: "false", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateBoolean")
	}
}

func TestValidateNumber(t *testing.T) {
	tests := []validateTest{
		{field: AgeField, value: "120", want: Valid},
		{field: AgeField, value: "-0", want: Valid},
		{field: AntElField, value: "90", want: Valid},
		{field: AntElField, value: "-90", want: Valid},
		{field: AIndexField, value: "123.0", want: Valid},
		{field: DistanceField, value: "9876.", want: Valid},
		{field: DistanceField, value: "1234567890", want: Valid},
		{field: FreqField, value: "146.520000001", want: Valid},
		{field: MaxBurstsField, value: "0", want: Valid},
		{field: MaxBurstsField, value: "00", want: Valid},
		{field: MyAltitudeField, value: "-1234.56789", want: Valid},
		{field: RxPwrField, value: ".7", want: Valid},
		{field: TxPwrField, value: "1499.999", want: Valid},
		{field: AgeField, value: "121", want: InvalidError},
		{field: AgeField, value: "-1", want: InvalidError},
		{field: AntElField, value: "--30", want: InvalidError},
		{field: AntElField, value: "99", want: InvalidError},
		{field: AntElField, value: "-91", want: InvalidError},
		{field: AntElField, value: "2π", want: InvalidError},
		{field: AIndexField, value: "-0.1", want: InvalidError},
		{field: AIndexField, value: "100-1", want: InvalidError},
		{field: AIndexField, value: "420", want: InvalidError},
		{field: DistanceField, value: "١٢٣", want: InvalidError},
		{field: DistanceField, value: "+9876", want: InvalidError},
		{field: DistanceField, value: "1 234", want: InvalidError},
		{field: DistanceField, value: "1,234", want: InvalidError},
		{field: DistanceField, value: "-0.00000001", want: InvalidError},
		{field: FreqField, value: "1.4652e2", want: InvalidError},
		{field: FreqField, value: "7.074-7", want: InvalidError},
		{field: MaxBurstsField, value: "1.2.3", want: InvalidError},
		{field: MaxBurstsField, value: "NaN", want: InvalidError},
		{field: MyAltitudeField, value: "〸", want: InvalidError},
		{field: MyAltitudeField, value: "⁷", want: InvalidError},
		{field: RxPwrField, value: "", want: InvalidError},
		{field: TxPwrField, value: ".", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateNumber")
	}
}

func TestValidateInteger(t *testing.T) {
	// currently all IntegerDataType fields have a minimum >= 0
	tests := []validateTest{
		{field: StxField, value: "0", want: Valid},
		{field: StxField, value: "1234567890", want: Valid},
		{field: NrBurstsField, value: "98765432123456789", want: Valid},
		{field: SfiField, value: "123", want: Valid},
		{field: KIndexField, value: "0", want: Valid},
		{field: KIndexField, value: "5", want: Valid},
		{field: KIndexField, value: "9", want: Valid},
		{field: StxField, value: "1,234", want: InvalidError},
		{field: StxField, value: "-1", want: InvalidError},
		{field: StxField, value: "7thirty", want: InvalidError},
		{field: StxField, value: "III", want: InvalidError},
		{field: NrBurstsField, value: "-1234", want: InvalidError},
		{field: NrBurstsField, value: "", want: InvalidError},
		{field: NrBurstsField, value: "twenty", want: InvalidError},
		{field: NrBurstsField, value: "௮", want: InvalidError},
		{field: NrBurstsField, value: "Ⅺ", want: InvalidError},
		{field: SfiField, value: "301", want: InvalidError},
		{field: SfiField, value: "7F", want: InvalidError},
		{field: SfiField, value: "0x20", want: InvalidError},
		{field: KIndexField, value: "10", want: InvalidError},
		{field: KIndexField, value: "-5", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateInteger")
	}
}

func TestValidatePositiveInteger(t *testing.T) {
	tests := []validateTest{
		{field: CqzField, value: "1", want: Valid},
		{field: CqzField, value: "40", want: Valid},
		{field: TenTenField, value: "1010101010", want: Valid},
		{field: FistsField, value: "1", want: Valid},
		{field: FistsField, value: "0987654321", want: Valid},
		{field: MyIotaIslandIdField, value: "1", want: Valid},
		{field: MyIotaIslandIdField, value: "666", want: Valid},
		{field: MyIotaIslandIdField, value: "99999999", want: Valid},
		{field: ItuzField, value: "1", want: Valid},
		{field: ItuzField, value: "90", want: Valid},
		{field: CqzField, value: "0", want: InvalidError},
		{field: CqzField, value: "42", want: InvalidError},
		{field: TenTenField, value: "0", want: InvalidError},
		{field: TenTenField, value: "-123", want: InvalidError},
		{field: FistsField, value: "five", want: InvalidError},
		{field: FistsField, value: "", want: InvalidError},
		{field: MyIotaIslandIdField, value: "0", want: InvalidError},
		{field: MyIotaIslandIdField, value: "100000000", want: InvalidError},
		{field: ItuzField, value: "111", want: InvalidError},
		{field: ItuzField, value: "99", want: InvalidError},
		{field: ItuzField, value: "８", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidatePositiveInteger")
	}
}

func TestValidateDate(t *testing.T) {
	tests := []validateTest{
		{field: QsoDateField, value: "19300101", want: Valid},
		{field: QsoDateOffField, value: "20200317", want: Valid},
		{field: QslrdateField, value: "19991231", want: Valid},
		{field: QslsdateField, value: "20000229", want: Valid},
		{field: QrzcomQsoUploadDateField, value: "22000101", want: Valid}, // future date, ctx.Now is zero
		{field: LotwQslrdateField, value: "23450607", want: Valid},        // future date, ctx.Now is zero
		{field: QsoDateField, value: "19000101", want: InvalidError},
		{field: QsoDateField, value: "19800100", want: InvalidError},
		{field: QsoDateOffField, value: "202012", want: InvalidError},
		{field: QsoDateOffField, value: "21000229", want: InvalidError},
		{field: QslrdateField, value: "1031", want: InvalidError},
		{field: QslsdateField, value: "2001-02-03", want: InvalidError},
		{field: QrzcomQsoUploadDateField, value: "01/02/2003", want: InvalidError},
		{field: LotwQslrdateField, value: "01022003", want: InvalidError},
		{field: LotwQslsdateField, value: "20220431", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateDate")
	}
}

func TestValidateTime(t *testing.T) {
	tests := []validateTest{
		{field: TimeOnField, value: "0000", want: Valid},
		{field: TimeOffField, value: "000000", want: Valid},
		{field: TimeOnField, value: "1234", want: Valid},
		{field: TimeOffField, value: "123456", want: Valid},
		{field: TimeOnField, value: "235959", want: Valid},
		{field: TimeOffField, value: "2000", want: Valid},
		{field: TimeOnField, value: "012345", want: Valid},
		{field: TimeOffField, value: "0159", want: Valid},
		{field: TimeOnField, value: "00:00", want: InvalidError},
		{field: TimeOffField, value: "00:00:00", want: InvalidError},
		{field: TimeOnField, value: "1234pm", want: InvalidError},
		{field: TimeOffField, value: "12.34.56", want: InvalidError},
		{field: TimeOnField, value: "735", want: InvalidError},
		{field: TimeOffField, value: "12345", want: InvalidError},
		{field: TimeOnField, value: "0630 AM", want: InvalidError},
		{field: TimeOffField, value: "noon", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateTime")
	}
}

func TestValidateDateRelative(t *testing.T) {
	ctx := ValidationContext{
		Now:        time.Date(2023, time.October, 31, 12, 34, 56, 0, time.UTC),
		FieldValue: emptyCtx.FieldValue}
	tests := []validateTest{
		{field: QsoDateField, value: "19991231", want: Valid},
		{field: QsoDateOffField, value: "20000101", want: Valid},
		{field: QslrdateField, value: "20231030", want: Valid},
		{field: QslsdateField, value: "20231031", want: Valid},
		{field: QrzcomQsoUploadDateField, value: "20231101", want: InvalidWarning}, // future date
		{field: LotwQslrdateField, value: "23450607", want: InvalidWarning},        // future date
		{field: QsoDateField, value: "19000101", want: InvalidError},
		{field: QsoDateField, value: "19800100", want: InvalidError},
		{field: QsoDateOffField, value: "202012", want: InvalidError},
		{field: QsoDateOffField, value: "21000229", want: InvalidError},
		{field: QslrdateField, value: "1031", want: InvalidError},
		{field: QslsdateField, value: "2001-02-03", want: InvalidError},
		{field: QrzcomQsoUploadDateField, value: "01/02/2003", want: InvalidError},
		{field: LotwQslrdateField, value: "01022003", want: InvalidError},
		{field: LotwQslsdateField, value: "20220431", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, ctx, "ValidateDate")
	}
}

func TestValidateTimeRelative(t *testing.T) {
	now := time.Date(2023, time.October, 31, 12, 34, 56, 0, time.UTC)
	yesterday := "20231030"
	today := "20231031"
	tomorrow := "20231101" // no warning, since the date field will warn
	tests := []struct {
		timeValue, dateValue string
		want                 Validity
	}{
		{timeValue: "1516", dateValue: yesterday, want: Valid},
		{timeValue: "235959", dateValue: yesterday, want: Valid},
		{timeValue: "1234", dateValue: today, want: Valid},
		{timeValue: "123456", dateValue: today, want: Valid},
		{timeValue: "1314", dateValue: today, want: InvalidWarning},
		{timeValue: "123457", dateValue: today, want: InvalidWarning},
		{timeValue: "1235", dateValue: today, want: InvalidWarning},
		{timeValue: "1516", dateValue: today, want: InvalidWarning},
		{timeValue: "0123", dateValue: tomorrow, want: Valid},
		{timeValue: "000000", dateValue: tomorrow, want: Valid},
	}
	tdFields := []struct{ tf, df Field }{
		{tf: TimeOnField, df: QsoDateField}, {tf: TimeOffField, df: QsoDateOffField},
	}
	for _, tc := range tests {
		for _, td := range tdFields {
			ctx := ValidationContext{
				Now: now,
				FieldValue: func(name string) string {
					if name == td.df.Name {
						return tc.dateValue
					}
					if name == td.tf.Name {
						return tc.timeValue
					}
					return ""
				},
			}
			testValidator(t, validateTest{field: td.tf, value: tc.timeValue, want: tc.want}, ctx, "ValidateTime")
		}
	}
}

func TestValidateString(t *testing.T) {
	tests := []validateTest{
		{field: ProgramidField, value: "Log & Operate", want: Valid},
		{field: CallField, value: "P5K2JI", want: Valid},
		{field: CommentField, value: `~!@#$%^&*()_+=-{}[]|\:;"'<>,.?/`, want: Valid},
		{field: EmailField, value: "flann.o-brien@example.com (Flann O`Brien)", want: Valid},
		{field: MyCityField, value: "", want: Valid},
		{field: ProgramidField, value: "🌲file", want: InvalidError},
		{field: CallField, value: "Я7ПД", want: InvalidError},
		{field: CommentField, value: `full width ．`, want: InvalidError},
		{field: EmailField, value: "thor@planet.earþ", want: InvalidError},
		{field: MyCityField, value: "\n", want: InvalidError},
		{field: CommentField, value: "Good\r\nchat", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateString")
	}
}

func TestValidateMultilineString(t *testing.T) {
	tests := []validateTest{
		{field: AddressField, value: "1600 Pennsylvania Ave\r\nWashington, DC 25000\r\n", want: Valid},
		{field: NotesField, value: "They were\non a boat", want: Valid},
		{field: QslmsgField, value: "\r", want: Valid},
		{field: RigField, value: "5 watts and a wire", want: Valid},
		{field: AddressField, value: "1600 Pennsylvania Ave\r\nWashington, ⏦ 25000\r\n", want: InvalidError},
		{field: NotesField, value: "Vertical\vtab", want: InvalidError},
		{field: QslmsgField, value: "🧶", want: InvalidError},
		{field: RigField, value: "Non‑breaking‑hyphen", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateMultilineString")
	}
}

func TestValidateIntlString(t *testing.T) {
	tests := []validateTest{
		{field: CommentIntlField, value: "१٢௩Ⅳ໖⁷၈🄊〸", want: Valid},
		{field: CountryIntlField, value: "🇭🇰", want: Valid},
		{field: MyAntennaIntlField, value: "null\0000character", want: Valid},
		{field: MyCityIntlField, value: "مكة المكرمة, Makkah the Honored", want: Valid},
		{field: MyCountryIntlField, value: "", want: Valid},
		{field: MyNameIntlField, value: "zero​width	space", want: Valid},
		{field: CommentIntlField, value: "new\nline", want: InvalidError},
		{field: CountryIntlField, value: "carriage\rreturn", want: InvalidError},
		{field: MyAntennaIntlField, value: "blank\r\n\r\nline", want: InvalidError},
		{field: MyCityIntlField, value: "line end\n", want: InvalidError},
		{field: MyCountryIntlField, value: "\r\n", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateIntlString")
	}
}

func TestValidateIntlMultilineString(t *testing.T) {
	tests := []validateTest{
		{field: AddressIntlField, value: "१٢௩Ⅳ໖⁷၈🄊〸", want: Valid},
		{field: NotesIntlField, value: "🇭🇰", want: Valid},
		{field: QslmsgIntlField, value: "null\0000character", want: Valid},
		{field: RigIntlField, value: "مكة المكرمة, Makkah the Honored", want: Valid},
		{field: AddressIntlField, value: "", want: Valid},
		{field: NotesIntlField, value: "zero​width	space", want: Valid},
		{field: QslmsgIntlField, value: "new\nline", want: Valid},
		{field: RigIntlField, value: "carriage\rreturn", want: Valid},
		{field: AddressIntlField, value: "blank\r\n\r\nline", want: Valid},
		{field: NotesIntlField, value: "line end\n", want: Valid},
		{field: QslmsgIntlField, value: "\r\n", want: Valid},
		{field: RigIntlField, value: "hello\tworld\r\n", want: Valid},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateIntlMultilineString")
	}
}

func TestValidateEnumeration(t *testing.T) {
	tests := []validateTest{
		{field: AntPathField, value: "G", want: Valid},
		{field: ArrlSectField, value: "wy", want: Valid},
		{field: BandField, value: "1.25m", want: Valid},
		{field: BandField, value: "13CM", want: Valid},
		{field: DxccField, value: "42", want: Valid},
		{field: QsoCompleteField, value: "?", want: Valid},
		{field: ContField, value: "na", want: Valid},
		{field: AntPathField, value: "X", want: InvalidError},
		{field: ArrlSectField, value: "TX", want: InvalidError},
		{field: ArrlSectField, value: "42", want: InvalidError},
		{field: BandField, value: "18m", want: InvalidError},
		{field: BandField, value: "130mm", want: InvalidError},
		{field: BandField, value: "G", want: InvalidError},
		{field: DxccField, value: "US", want: InvalidError},
		{field: QsoCompleteField, value: "!", want: InvalidError},
		{field: ContField, value: "Europe", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateEnumeration")
	}
}

func TestValidateEnumScope(t *testing.T) {
	tests := []struct {
		validateTest
		values map[string]string
	}{
		// empty string is valid
		{
			validateTest: validateTest{field: StateField, value: "", want: Valid},
			values:       map[string]string{},
		},
		{
			validateTest: validateTest{field: StateField, value: "", want: Valid},
			values:       map[string]string{DxccField.Name: ""},
		},
		// Wyoming, USA
		{
			validateTest: validateTest{field: MyStateField, value: "WY", want: Valid},
			values:       map[string]string{MyDxccField.Name: CountryUnitedStatesOfAmerica.EntityCode},
		},
		// Yukon Territory, Canada
		{
			validateTest: validateTest{field: StateField, value: "YT", want: Valid},
			values:       map[string]string{DxccField.Name: "1"},
		},
		// Sichuan, China
		{
			validateTest: validateTest{field: MyStateField, value: "SC", want: Valid},
			values:       map[string]string{MyDxccField.Name: "318"},
		},
		// Nagano, Japan
		{
			validateTest: validateTest{field: StateField, value: "09", want: Valid},
			values:       map[string]string{DxccField.Name: "339"},
		},
		// Vranov nad Toplou, Slovak Republic
		{
			validateTest: validateTest{field: MyStateField, value: "vrt", want: Valid},
			values:       map[string]string{MyDxccField.Name: "504"},
		},
		// Permskaya oblast (new) or Permskaya Kraj (old) in Russia
		{
			validateTest: validateTest{field: StateField, value: "PM", want: Valid},
			values:       map[string]string{DxccField.Name: "15"},
		},

		// St. John, Barbados, but ADIF doesn't define Barbadan subdivisions
		{
			validateTest: validateTest{field: MyStateField, value: "JN", want: InvalidWarning},
			values:       map[string]string{MyDxccField.Name: CountryBarbados.EntityCode, "MY_STATE": "JN"},
		},
		// not an empty string
		{
			validateTest: validateTest{field: StateField, value: "  ", want: InvalidError},
			values:       map[string]string{DxccField.Name: "1"},
		},
		// not an empty string, but DXCC not set
		{
			validateTest: validateTest{field: StateField, value: "  ", want: InvalidWarning},
			values:       map[string]string{},
		},
		// Not a state abbreviation in any country
		{
			validateTest: validateTest{field: MyStateField, value: "XYZ", want: InvalidError},
			values:       map[string]string{MyDxccField.Name: "291"},
		},
		// Not a state abbreviation in any country, DXCC not set
		{
			validateTest: validateTest{field: MyStateField, value: "XYZ", want: InvalidWarning},
			values:       map[string]string{},
		},
		// Not a valid DXCC code
		{
			validateTest: validateTest{field: MyStateField, value: "CA", want: InvalidWarning},
			values:       map[string]string{MyDxccField.Name: "9876"},
		},
		// Wyoming, presumably but not certain
		{
			validateTest: validateTest{field: MyStateField, value: "WY", want: InvalidWarning},
			values:       map[string]string{},
		},
		// Yukon Territory, but wrong country
		{
			validateTest: validateTest{field: StateField, value: "YT", want: InvalidError},
			values:       map[string]string{DxccField.Name: "291"},
		},
		// Sichuan, but not abbreviated
		{
			validateTest: validateTest{field: MyStateField, value: "Sichuan", want: InvalidError},
			values:       map[string]string{MyDxccField.Name: "318"},
		},
		// Sichuan, not abbreviated, but DXCC is invalid; warn on subdivision, DXCC validation will be error
		{
			validateTest: validateTest{field: MyStateField, value: "Sichuan", want: InvalidWarning},
			values:       map[string]string{MyDxccField.Name: "China"},
		},
		// Nagano, but missing leading 0
		{
			validateTest: validateTest{field: StateField, value: "9", want: InvalidError},
			values:       map[string]string{DxccField.Name: "339"},
		},
		// Vranov nad Toplou, but Slovak Republic is spelled out
		{
			validateTest: validateTest{field: MyStateField, value: "VRT", want: InvalidWarning},
			values:       map[string]string{MyDxccField.Name: "Slovak Republic"},
		},
		// PM in Russia, but Cyrilic characters
		{
			validateTest: validateTest{field: StateField, value: "ПМ", want: InvalidError},
			values:       map[string]string{DxccField.Name: "15"},
		},
	}
	for _, tc := range tests {
		ctx := ValidationContext{FieldValue: func(name string) string { return tc.values[name] }}
		testValidator(t, tc.validateTest, ctx, "ValidateEnumScope")
	}
}

func TestValidateStringEnumScope(t *testing.T) {
	// Submode is a String field but has an associated enumeration
	tests := []struct {
		validateTest
		values map[string]string
	}{
		{
			validateTest: validateTest{field: SubmodeField, value: "", want: Valid},
			values:       map[string]string{},
		},
		{
			validateTest: validateTest{field: SubmodeField, value: "", want: Valid},
			values:       map[string]string{ModeField.Name: "CW", SubmodeField.Name: ""},
		},
		{
			validateTest: validateTest{field: SubmodeField, value: "THOR-M", want: Valid},
			values:       map[string]string{ModeField.Name: "THOR", SubmodeField.Name: "THOR-M"},
		},
		{
			validateTest: validateTest{field: SubmodeField, value: "LSB", want: Valid},
			values:       map[string]string{ModeField.Name: "SSB", SubmodeField.Name: "LSB"},
		},
		{
			validateTest: validateTest{field: SubmodeField, value: "lsb", want: Valid},
			values:       map[string]string{ModeField.Name: "sSb", SubmodeField.Name: "lsb"},
		},
		{
			validateTest: validateTest{field: SubmodeField, value: "PSK31", want: Valid},
			values:       map[string]string{ModeField.Name: "PSK", SubmodeField.Name: "PSK31"},
		},
		{
			validateTest: validateTest{field: SubmodeField, value: "Dstar", want: Valid},
			values:       map[string]string{ModeField.Name: "DIGITALVOICE", SubmodeField.Name: "Dstar"},
		},
		{
			validateTest: validateTest{field: SubmodeField, value: "OLIVIA 16/500", want: Valid},
			values:       map[string]string{ModeField.Name: "OLIVIA", SubmodeField.Name: "OLIVIA 16/500"},
		},

		// valid submode, but mode not set
		{
			validateTest: validateTest{field: SubmodeField, value: "THOR-M", want: InvalidWarning},
			values:       map[string]string{ModeField.Name: "", SubmodeField.Name: "THOR-M"},
		},
		{
			validateTest: validateTest{field: SubmodeField, value: "PSK31", want: InvalidWarning},
			values:       map[string]string{SubmodeField.Name: "PSK31"},
		},
		// Unkknoown submodes
		{
			validateTest: validateTest{field: SubmodeField, value: "NOTAMODE", want: InvalidWarning},
			values:       map[string]string{},
		},
		{
			validateTest: validateTest{field: SubmodeField, value: "lower", want: InvalidWarning},
			values:       map[string]string{ModeField.Name: "SSB", SubmodeField.Name: "lower"},
		},
		// mode/submode mismatch
		{
			validateTest: validateTest{field: SubmodeField, value: "LSB", want: InvalidWarning},
			values:       map[string]string{ModeField.Name: "FM", SubmodeField.Name: "LSB"},
		},
		{
			validateTest: validateTest{field: SubmodeField, value: "PSK31", want: InvalidWarning},
			values:       map[string]string{ModeField.Name: "CW", SubmodeField.Name: "PSK31"},
		},
		{
			validateTest: validateTest{field: SubmodeField, value: "VOCODER", want: InvalidWarning},
			values:       map[string]string{ModeField.Name: "DIGITALVOICE", SubmodeField.Name: "VOCODER"},
		},

		{
			validateTest: validateTest{field: SubmodeField, value: "USB", want: InvalidWarning},
			values:       map[string]string{ModeField.Name: "CW", SubmodeField.Name: "USB"},
		},
		// submode entirely unknown, submode is string so only warning
		{
			validateTest: validateTest{field: SubmodeField, value: "TotalJunk", want: InvalidWarning},
			values:       map[string]string{ModeField.Name: "RTTY", SubmodeField.Name: "TotalJunk"},
		},
		// mode doesn't have any submodes
		{
			validateTest: validateTest{field: SubmodeField, value: "TotalJunk", want: InvalidWarning},
			values:       map[string]string{ModeField.Name: "FAX", SubmodeField.Name: "TotalJunk"},
		},
		{
			validateTest: validateTest{field: SubmodeField, value: "ASCI", want: InvalidWarning},
			values:       map[string]string{ModeField.Name: "FAX", SubmodeField.Name: "ASCI"},
		},
		// mode/submode swapped
		{
			validateTest: validateTest{field: SubmodeField, value: "SSB", want: InvalidWarning},
			values:       map[string]string{ModeField.Name: "USB", SubmodeField.Name: "SSB"},
		},
	}
	for _, tc := range tests {
		ctx := ValidationContext{FieldValue: func(name string) string { return tc.values[name] }}
		testValidator(t, tc.validateTest, ctx, "ValidateStringScope")
	}
}

func TestValidateContestID(t *testing.T) {
	// CONTEST_ID is associated with the Contest_Id enum but it's a String field,
	// allowing non-enumerated contests (but encouraging enum values for
	// interoperability).
	tests := []validateTest{
		{field: ContestIdField, value: "", want: Valid},
		{field: ContestIdField, value: "CO-QSO-PARTY", want: Valid},
		{field: ContestIdField, value: "CO_QSO_PARTY", want: InvalidWarning},
		{field: ContestIdField, value: "CO QSO PARTY", want: InvalidWarning},
		// Enum value for IL has spaces rather than dashes for some reason
		{field: ContestIdField, value: "IL QSO PARTY", want: Valid},
		{field: ContestIdField, value: "IL-QSO-PARTY", want: InvalidWarning},
		{field: ContestIdField, value: "IL_QSO_PARTY", want: InvalidWarning},
		{field: ContestIdField, value: "070-31-FLAVORS", want: Valid},
		{field: ContestIdField, value: "070-31-flavors", want: Valid},
		{field: ContestIdField, value: "70-31-FLAVORS", want: InvalidWarning},
		{field: ContestIdField, value: "RAC", want: Valid},
		{field: ContestIdField, value: "R.A.C.", want: InvalidWarning},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateContestId")
	}
}

func TestValidateGridsquare(t *testing.T) {
	tests := []validateTest{
		{field: GridsquareField, value: "", want: Valid},
		// First letter pair is only valid A-R, thuogh we don't check this
		{field: GridsquareField, value: "AA", want: Valid},
		{field: MyGridsquareField, value: "rr", want: Valid},
		{field: GridsquareField, value: "AA00", want: Valid},
		{field: MyGridsquareField, value: "CD12", want: Valid},
		{field: GridsquareField, value: "jk28", want: Valid},
		{field: MyGridsquareField, value: "XX99", want: Valid},
		// Second letter pair is only valid A-X, thuogh we don't check this
		{field: GridsquareField, value: "AB34ef", want: Valid},
		{field: MyGridsquareField, value: "gh56IJ", want: Valid},
		{field: GridsquareField, value: "KL78mn", want: Valid},
		{field: MyGridsquareField, value: "QR09XX", want: Valid},
		{field: GridsquareField, value: "AA00xx99", want: Valid},
		{field: MyGridsquareField, value: "rh63NG50", want: Valid},
		{field: GridsquareField, value: "rr99aa00", want: Valid},
		{field: MyGridsquareField, value: ",", want: InvalidError},
		{field: MyGridsquareField, value: "F,", want: InvalidError},
		{field: MyGridsquareField, value: "JK3,", want: InvalidError},
		{field: GridsquareField, value: "12", want: InvalidError},
		{field: MyGridsquareField, value: "12CD", want: InvalidError},
		{field: GridsquareField, value: "12CD34", want: InvalidError},
		{field: MyGridsquareField, value: "12CD34GH", want: InvalidError},
		{field: GridsquareField, value: "AB12C3", want: InvalidError},
		{field: MyGridsquareField, value: "AB123C", want: InvalidError},
		{field: GridsquareField, value: "AB12CDEF", want: InvalidError},
		{field: MyGridsquareField, value: "AB12CD3F", want: InvalidError},
		{field: GridsquareField, value: "AB12CDE3", want: InvalidError},
		{field: MyGridsquareField, value: "AB12CD34EF", want: InvalidError},
		{field: GridsquareField, value: "AB12CD34EF56", want: InvalidError},
		{field: MyGridsquareField, value: "AB 12", want: InvalidError},
		{field: GridsquareField, value: "AB 12 cd", want: InvalidError},
		{field: MyGridsquareField, value: "AB 12 cd 34", want: InvalidError},
		{field: GridsquareField, value: "ÐÞ12", want: InvalidError},
		{field: GridsquareField, value: "HI๕๙", want: InvalidError},
		{field: MyGridsquareField, value: "JK99dé", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateGridsquare")
	}
}

func TestGridsquareExt(t *testing.T) {
	// Gridsquare extension has max length 4 (2 letters, 2 numbers)
	tests := []validateTest{
		{field: GridsquareExtField, value: "", want: Valid},
		// First letter pair is only valid A-R, thuogh we don't check this
		{field: GridsquareExtField, value: "AA", want: Valid},
		{field: MyGridsquareExtField, value: "rr", want: Valid},
		{field: GridsquareExtField, value: "AA00", want: Valid},
		{field: MyGridsquareExtField, value: "CD12", want: Valid},
		{field: GridsquareExtField, value: "jk28", want: Valid},
		{field: MyGridsquareExtField, value: "XX99", want: Valid},
		// Second letter pair is only valid A-X, thuogh we don't check this
		{field: GridsquareExtField, value: "AB34ef", want: InvalidError},
		{field: MyGridsquareExtField, value: "gh56IJ", want: InvalidError},
		{field: GridsquareExtField, value: "KL78mn", want: InvalidError},
		{field: MyGridsquareExtField, value: "QR09XX", want: InvalidError},
		{field: GridsquareExtField, value: "AA00xx99", want: InvalidError},
		{field: MyGridsquareExtField, value: "rh63NG50", want: InvalidError},
		{field: GridsquareExtField, value: "rr99aa00", want: InvalidError},
		{field: MyGridsquareExtField, value: ",", want: InvalidError},
		{field: MyGridsquareExtField, value: "F,", want: InvalidError},
		{field: MyGridsquareExtField, value: "JK8,", want: InvalidError},
		{field: MyGridsquareExtField, value: "G73,", want: InvalidError},
		{field: GridsquareExtField, value: "12", want: InvalidError},
		{field: MyGridsquareExtField, value: "12CD", want: InvalidError},
		{field: GridsquareExtField, value: "12CD34", want: InvalidError},
		{field: MyGridsquareExtField, value: "AB 12", want: InvalidError},
		{field: GridsquareExtField, value: "ÐÞ12", want: InvalidError},
		{field: GridsquareExtField, value: "HI๕๙", want: InvalidError},
		{field: MyGridsquareExtField, value: "DÉ41", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateGridsquareExt")
	}
}

func TestGridsquareList(t *testing.T) {
	// Validator doesn't currently check that locators are adjacent
	tests := []validateTest{
		{field: VuccGridsField, value: "", want: Valid},
		{field: MyVuccGridsField, value: ",", want: Valid},
		{field: MyVuccGridsField, value: ",HI59", want: Valid},
		{field: VuccGridsField, value: "DN70,DM79", want: Valid},
		{field: MyVuccGridsField, value: "DN70,DM79,DN80,DM89", want: Valid},
		{field: VuccGridsField, value: "AA00xx,RR99aa", want: Valid},
		{field: MyVuccGridsField, value: "AA00xx,RR99aa,GO83NU", want: Valid},
		{field: VuccGridsField, value: "Rq98Po87,nm65LK43,JH21,ig", want: Valid},
		{field: VuccGridsField, value: " ", want: InvalidError},
		{field: MyVuccGridsField, value: " , ", want: InvalidError},
		{field: VuccGridsField, value: "DN70;DM79", want: InvalidError},
		{field: VuccGridsField, value: "AB12cd,ÃB13gf", want: InvalidError},
		{field: MyVuccGridsField, value: "AB12cd,DE34fg,56hijk", want: InvalidError},
		{field: MyVuccGridsField, value: "RR99,AA0", want: InvalidError},
		{field: VuccGridsField, value: "AB12cd,EF34gh56jk", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateGridsquareList")
	}
}

func TestValidateLocation(t *testing.T) {
	tests := []validateTest{
		{field: LatField, value: "", want: Valid},
		{field: LonField, value: "", want: Valid},
		{field: LatField, value: "N000 00.000", want: Valid, others: []validateTest{{field: LonField, value: "E000 00.000"}}},
		{field: LatField, value: "S000 00.000", want: Valid, others: []validateTest{{field: LonField, value: "W000 00.000"}}},
		{field: LonField, value: "E000 00.000", want: Valid, others: []validateTest{{field: LatField, value: "N000 00.000"}}},
		{field: LonField, value: "W000 00.000", want: Valid, others: []validateTest{{field: LatField, value: "S000 00.000"}}},
		{field: MyLatField, value: "N001 00.000", want: Valid, others: []validateTest{{field: MyLonField, value: "W123 45.678"}}},
		{field: MyLatField, value: "S090 00.000", want: Valid, others: []validateTest{{field: MyLonField, value: "E123 45.678"}}},
		{field: MyLonField, value: "E009 00.000", want: Valid, others: []validateTest{{field: MyLatField, value: "S012 34.567"}}},
		{field: MyLonField, value: "W180 00.000", want: Valid, others: []validateTest{{field: MyLatField, value: "S012 34.567"}}},
		{field: LatField, value: "n012 34.567", want: Valid, others: []validateTest{{field: LonField, value: "E000 00.000"}}},
		{field: LatField, value: "s023 59.999", want: Valid, others: []validateTest{{field: LonField, value: "W000 00.000"}}},
		{field: LonField, value: "e000 00.001", want: Valid, others: []validateTest{{field: LatField, value: "S000 00.000"}}},
		{field: LonField, value: "w120 30.050", want: Valid, others: []validateTest{{field: LatField, value: "N000 00.000"}}},
		{field: MyLatField, value: "W012 34.567", want: InvalidError, others: []validateTest{{field: MyLonField, value: "E000 00.000"}}},
		{field: MyLonField, value: "S012 34.567", want: InvalidError, others: []validateTest{{field: MyLatField, value: "N000 00.000"}}},
		{field: LatField, value: "N091 00.000", want: InvalidError, others: []validateTest{{field: LonField, value: "E000 00.000"}}},
		{field: LonField, value: "W181 00.000", want: InvalidError, others: []validateTest{{field: LatField, value: "N000 00.000"}}},
		{field: MyLatField, value: "N090 00.001", want: InvalidError, others: []validateTest{{field: MyLonField, value: "E000 00.000"}}},
		{field: MyLonField, value: "W180 00.001", want: InvalidError, others: []validateTest{{field: MyLatField, value: "N000 00.000"}}},
		{field: LatField, value: "Equator", want: InvalidError},
		{field: LonField, value: "Greenwich", want: InvalidError},
		{field: LatField, value: "X012 34.567", want: InvalidError},
		{field: LatField, value: "N01234.567", want: InvalidError},
		{field: LatField, value: "N012 3.4567", want: InvalidError},
		{field: LatField, value: "S23 59.999", want: InvalidError},
		{field: LonField, value: "E0 00.001", want: InvalidError},
		{field: LonField, value: "W120 30.05", want: InvalidError},
		{field: LonField, value: "W123 4.567", want: InvalidError},
		{field: LonField, value: "120 30.050", want: InvalidError},
		{field: LonField, value: "E123 45.678extra", want: InvalidError},
		{field: MyLatField, value: "12.345678", want: InvalidError},
		{field: MyLatField, value: "-6.543", want: InvalidError},
		{field: MyLonField, value: "123 45.678E", want: InvalidError},
		{field: MyLonField, value: `123 45' 54.321"`, want: InvalidError},
		{field: MyLonField, value: `-123° 45' 54.321"`, want: InvalidError},
		// Warning if one coordinate field is set but not the other
		{field: LatField, value: "N012 34.567", want: InvalidWarning},
		{field: LonField, value: "E012 34.567", want: InvalidWarning},
		{field: MyLatField, value: "N012 34.567", want: InvalidWarning},
		{field: MyLonField, value: "E012 34.567", want: InvalidWarning},
	}
	for _, tc := range tests {
		testValidator(t, tc, validateTestCtx(tc), "ValidateLocation")
	}
}

func TestValidateIOTARef(t *testing.T) {
	tests := []validateTest{
		{field: IotaField, value: "", want: Valid},
		{field: IotaField, value: "AF-001", want: Valid},
		{field: IotaField, value: "AN-999", want: Valid},
		{field: IotaField, value: "AS-123", want: Valid},
		{field: IotaField, value: "EU-321", want: Valid},
		{field: IotaField, value: "na-500", want: Valid},
		{field: IotaField, value: "OC-842", want: Valid},
		{field: IotaField, value: "SA-767", want: Valid},
		{field: IotaField, value: "-", want: InvalidError},
		{field: IotaField, value: "987", want: InvalidError},
		{field: IotaField, value: "AF-1", want: InvalidError},
		{field: IotaField, value: "AN-9999", want: InvalidError},
		{field: IotaField, value: "AS123", want: InvalidError},
		{field: IotaField, value: "-321", want: InvalidError},
		{field: IotaField, value: "AQ-123", want: InvalidError},
		{field: IotaField, value: "EUROPE-555", want: InvalidError},
		{field: IotaField, value: "N.A-876", want: InvalidError},
		{field: IotaField, value: "SA 432", want: InvalidError},
		{field: IotaField, value: "AF-345,AF-678", want: InvalidError},
		{field: IotaField, value: "AS-12C", want: InvalidError},
		{field: IotaField, value: "OC-022 Bali", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateIOTARef")
	}
}

func TestValidatePOTARef(t *testing.T) {
	tests := []validateTest{
		{field: PotaRefField, value: "", want: Valid},
		// 2024 introduced ISO-3166 codes
		{field: MyPotaRefField, value: "US-0001", want: Valid},
		{field: PotaRefField, value: "CY-0123", want: Valid},
		{field: MyPotaRefField, value: "jp-1491", want: Valid},
		{field: PotaRefField, value: "AN-0001", want: Valid},
		{field: MyPotaRefField, value: "cl-0016@cl-at", want: Valid},
		{field: PotaRefField, value: "PR-0001@PR-PR", want: Valid},
		{field: MyPotaRefField, value: "MY-0129@MY-08", want: Valid},
		{field: PotaRefField, value: "GB-0086@GB-ENG", want: Valid},
		{field: MyPotaRefField, value: "US-0028,CA-0110", want: Valid},
		{field: PotaRefField, value: "FJ-0003,FJ-0004,FJ-0005", want: Valid},
		{field: MyPotaRefField, value: "CN-0004@CN-SN,CN-0004@CN-SX", want: Valid},
		// pre-2024 style
		{field: MyPotaRefField, value: "K-0001", want: Valid},
		{field: PotaRefField, value: "5B-0123", want: Valid},
		{field: MyPotaRefField, value: "ja-1491", want: Valid},
		{field: PotaRefField, value: "CE9-0001", want: Valid},
		{field: MyPotaRefField, value: "ca-0016@cl-at", want: Valid},
		{field: PotaRefField, value: "K-0103@US-KP4", want: Valid},
		{field: MyPotaRefField, value: "9M-0129@9M-WM", want: Valid},
		{field: PotaRefField, value: "G-0086@G-NORTHANTS", want: Valid},
		{field: MyPotaRefField, value: "K-0028,VE-0110", want: Valid},
		{field: PotaRefField, value: "3D2-0003,3D2-0004,3D2-0005", want: Valid},
		{field: MyPotaRefField, value: "BY-0004@CN-SA,BY-0004@CN-SX", want: Valid},
		// invalid
		{field: PotaRefField, value: "-", want: InvalidError},
		{field: MyPotaRefField, value: "K-1", want: InvalidError},
		{field: PotaRefField, value: "5B-123", want: InvalidError},
		{field: MyPotaRefField, value: "LA1234", want: InvalidError},
		{field: PotaRefField, value: "-1491", want: InvalidError},
		{field: MyPotaRefField, value: "CE9", want: InvalidError},
		{field: PotaRefField, value: "ca-0016@-at", want: InvalidError},
		{field: MyPotaRefField, value: "K-0103@KP4-PR", want: InvalidError},
		{field: PotaRefField, value: "4X-0040@4X-HA", want: InvalidError},
		{field: MyPotaRefField, value: "G-0048@G-NOTTINGHAMSHIRE", want: InvalidError},
		{field: PotaRefField, value: "K-0028 , VE-0110", want: InvalidError},
		{field: MyPotaRefField, value: "3D2-0003,3D2-0004,3D20005", want: InvalidError},
		{field: PotaRefField, value: "BY-0004@CNSA,BY-0004@CN-SX", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidatePOTARef")
	}
}

func TestValidateSOTARef(t *testing.T) {
	tests := []validateTest{
		{field: SotaRefField, value: "", want: Valid},
		{field: MySotaRefField, value: "F/AB-456", want: Valid},
		{field: SotaRefField, value: "ZS/WC-971", want: Valid},
		{field: MySotaRefField, value: "CE3/CO-099", want: Valid},
		{field: SotaRefField, value: "w0c/fr-226", want: Valid},
		{field: MySotaRefField, value: "9A/DH-173", want: Valid},
		{field: SotaRefField, value: "3Y/BV-001", want: Valid},
		{field: MySotaRefField, value: "/", want: InvalidError},
		{field: SotaRefField, value: "F-AB/456", want: InvalidError},
		{field: MySotaRefField, value: "WC-971", want: InvalidError},
		{field: SotaRefField, value: "CE3/CO-99", want: InvalidError},
		{field: MySotaRefField, value: "PY2/SAO-123", want: InvalidError},
		{field: SotaRefField, value: "VP8/SG-1234", want: InvalidError},
		{field: MySotaRefField, value: "ALGERIA/JK-123", want: InvalidError},
		{field: SotaRefField, value: "JA/NN-211 Fuji", want: InvalidError},
		{field: MySotaRefField, value: "W0C/FR-110,W0C/FR-118", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateSOTARef")
	}
}

func TestValidateWWFFRef(t *testing.T) {
	tests := []validateTest{
		{field: WwffRefField, value: "", want: Valid},
		{field: MyWwffRefField, value: "FFF-3210", want: Valid},
		{field: WwffRefField, value: "1SFF-0001", want: Valid},
		{field: MyWwffRefField, value: "cuff-0193", want: Valid},
		{field: WwffRefField, value: "3D2FF-0007", want: Valid},
		{field: MyWwffRefField, value: "-", want: InvalidError},
		{field: WwffRefField, value: "FFF/3210", want: InvalidError},
		{field: MyWwffRefField, value: "1SFF 0001", want: InvalidError},
		{field: WwffRefField, value: "cuff-12345", want: InvalidError},
		{field: MyWwffRefField, value: "3D2FF-7", want: InvalidError},
		{field: WwffRefField, value: "ZSAA-1234", want: InvalidError},
		{field: MyWwffRefField, value: "SHFF-0014 Serengeti", want: InvalidError},
		{field: WwffRefField, value: "KFF-0043,KFF-0042", want: InvalidError},
	}
	for _, tc := range tests {
		testValidator(t, tc, emptyCtx, "ValidateWWFFRef")
	}
}

func TestValidateCQZone(t *testing.T) {
	tests := []struct {
		validateTest
		dxcc, mydxcc, country, mycountry string
	}{
		{validateTest: validateTest{field: CqzField, value: "38", want: Valid}},
		{validateTest: validateTest{field: CqzField, value: "9", want: Valid}, dxcc: CountryTrinidadTobago.EntityCode},
		{validateTest: validateTest{field: CqzField, value: "09", want: Valid}, dxcc: CountryTrinidadTobago.EntityCode},
		{validateTest: validateTest{field: CqzField, value: "29", want: InvalidError}, dxcc: CountryTrinidadTobago.EntityCode},
		{validateTest: validateTest{field: MyCqZoneField, value: "40", want: Valid}, mydxcc: CountryIceland.EntityCode},
		{validateTest: validateTest{field: MyCqZoneField, value: "14", want: InvalidError}, mydxcc: CountryIceland.EntityCode},
		{validateTest: validateTest{field: CqzField, value: "24", want: Valid}, dxcc: CountryChina.EntityCode},
		{validateTest: validateTest{field: CqzField, value: "19", want: InvalidError}, dxcc: CountryChina.EntityCode},
		{validateTest: validateTest{field: CqzField, value: "37", want: Valid}, country: CountryYemen.EntityName},
		{validateTest: validateTest{field: MyCqZoneField, value: "39", want: InvalidError}, mycountry: CountryYemen.EntityName},
		{validateTest: validateTest{field: MyCqZoneField, value: "28", want: Valid}, mycountry: CountryAntarctica.EntityName},
		{validateTest: validateTest{field: MyCqZoneField, value: "1", want: InvalidError}, mycountry: CountryAntarctica.EntityName},
	}
	for _, tc := range tests {
		ctx := ValidationContext{FieldValue: func(name string) string {
			switch name {
			case DxccField.Name:
				return tc.dxcc
			case MyDxccField.Name:
				return tc.mydxcc
			case CountryField.Name:
				return tc.country
			case MyCountryField.Name:
				return tc.mycountry
			default:
				return ""
			}
		}}
		testValidator(t, tc.validateTest, ctx, "TestValidateCQZone")
	}
}

func TestValidateITUZone(t *testing.T) {
	tests := []struct {
		validateTest
		dxcc, mydxcc, country, mycountry string
	}{
		{validateTest: validateTest{field: ItuzField, value: "90", want: Valid}},
		{validateTest: validateTest{field: ItuzField, value: "1", want: Valid}, dxcc: CountryAlaska.EntityCode},
		{validateTest: validateTest{field: ItuzField, value: "02", want: Valid}, dxcc: CountryAlaska.EntityCode},
		{validateTest: validateTest{field: ItuzField, value: "12", want: InvalidError}, dxcc: CountryAlaska.EntityCode},
		{validateTest: validateTest{field: MyItuZoneField, value: "45", want: Valid}, mydxcc: CountryJapan.EntityCode},
		{validateTest: validateTest{field: MyItuZoneField, value: "25", want: InvalidError}, mydxcc: CountryJapan.EntityCode},
		{validateTest: validateTest{field: ItuzField, value: "33", want: Valid}, dxcc: CountryAsiaticRussia.EntityCode},
		{validateTest: validateTest{field: ItuzField, value: "42", want: InvalidError}, dxcc: CountryAsiaticRussia.EntityCode},
		{validateTest: validateTest{field: ItuzField, value: "36", want: Valid}, country: CountryCanaryIslands.EntityName},
		{validateTest: validateTest{field: MyItuZoneField, value: "37", want: InvalidError}, mycountry: CountryCanaryIslands.EntityName},
		{validateTest: validateTest{field: MyItuZoneField, value: "69", want: Valid}, mycountry: CountryAntarctica.EntityName},
		{validateTest: validateTest{field: MyItuZoneField, value: "75", want: InvalidError}, mycountry: CountryAntarctica.EntityName},
	}
	for _, tc := range tests {
		ctx := ValidationContext{FieldValue: func(name string) string {
			switch name {
			case DxccField.Name:
				return tc.dxcc
			case MyDxccField.Name:
				return tc.mydxcc
			case CountryField.Name:
				return tc.country
			case MyCountryField.Name:
				return tc.mycountry
			default:
				return ""
			}
		}}
		testValidator(t, tc.validateTest, ctx, "TestValidateCQZone")
	}
}

func TestValidateContinent(t *testing.T) {
	tests := []struct {
		validateTest
		dxcc, country string
	}{
		{validateTest: validateTest{field: ContField, value: "SA", want: Valid}, dxcc: CountryTrinidadTobago.EntityCode, country: ""},
		{validateTest: validateTest{field: ContField, value: "OC", want: Valid}, dxcc: "", country: CountryEastMalaysia.EntityName},
		{validateTest: validateTest{field: ContField, value: "AS", want: Valid}, dxcc: "", country: CountryWestMalaysia.EntityName},
		{validateTest: validateTest{field: ContField, value: "EU", want: InvalidWarning}, dxcc: CountryCanaryIslands.EntityCode, country: ""},
		{validateTest: validateTest{field: ContField, value: "EU", want: Valid}, dxcc: "", country: "Iceland"},
		{validateTest: validateTest{field: ContField, value: "AS", want: InvalidWarning}, dxcc: CountryBolivia.EntityCode, country: ""},
	}
	for _, tc := range tests {
		ctx := ValidationContext{FieldValue: func(name string) string {
			switch name {
			case DxccField.Name:
				return tc.dxcc
			case CountryField.Name:
				return tc.country
			default:
				return ""
			}
		}}
		testValidator(t, tc.validateTest, ctx, "TestValidateCountry")
	}
}

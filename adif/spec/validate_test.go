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

import "testing"

// TODO numbers, dates, times, enumerations

type validateTest struct {
	field   Field
	value   string
	wantErr bool
}

func testValidator(t *testing.T, tc validateTest, funcname string) {
	t.Helper()
	var ctx ValidationContext
	v := TypeValidators[tc.field.Type.Name]
	if err := v(tc.value, tc.field, ctx); (err != nil) != tc.wantErr {
		if tc.wantErr {
			t.Errorf("%s(%q, %s, ctx) got nil, want error", funcname, tc.value, tc.field.Name)
		} else {
			t.Errorf("%s(%q, %s, ctx) got error %v", funcname, tc.value, tc.field.Name, err)
		}
	}
}

func TestValidateBoolean(t *testing.T) {
	tests := []validateTest{
		{field: QsoRandomField, value: "Y", wantErr: false},
		{field: SilentKeyField, value: "y", wantErr: false},
		{field: ForceInitField, value: "N", wantErr: false},
		{field: SwlField, value: "n", wantErr: false},
		{field: QsoRandomField, value: "YES", wantErr: true},
		{field: SilentKeyField, value: "true", wantErr: true},
		{field: ForceInitField, value: "F", wantErr: true},
		{field: SwlField, value: "false", wantErr: true},
	}
	for _, tc := range tests {
		testValidator(t, tc, "ValidateBoolean")
	}
}

func TestValidateString(t *testing.T) {
	tests := []validateTest{
		{field: ProgramidField, value: "Log & Operate", wantErr: false},
		{field: CallField, value: "P5K2JI", wantErr: false},
		{field: CommentField, value: `~!@#$%^&*()_+=-{}[]|\:;"'<>,.?/`, wantErr: false},
		{field: EmailField, value: "flann.o-brien@example.com (Flann O`Brien)", wantErr: false},
		{field: MyCityField, value: "", wantErr: false},
		{field: ProgramidField, value: "🌲file", wantErr: true},
		{field: CallField, value: "Я7ПД", wantErr: true},
		{field: CommentField, value: `full width ．`, wantErr: true},
		{field: EmailField, value: "thor@planet.earþ", wantErr: true},
		{field: MyCityField, value: "\n", wantErr: true},
		{field: CommentField, value: "Good\r\nchat", wantErr: true},
	}
	for _, tc := range tests {
		testValidator(t, tc, "ValidateString")
	}
}

func TestValidateMultilineString(t *testing.T) {
	tests := []validateTest{
		{field: AddressField, value: "1600 Pennsylvania Ave\r\nWashington, DC 25000\r\n", wantErr: false},
		{field: NotesField, value: "They were\non a boat", wantErr: false},
		{field: QslmsgField, value: "\r", wantErr: false},
		{field: RigField, value: "5 watts and a wire", wantErr: false},
		{field: AddressField, value: "1600 Pennsylvania Ave\r\nWashington, ⏦ 25000\r\n", wantErr: true},
		{field: NotesField, value: "Vertical\vtab", wantErr: true},
		{field: QslmsgField, value: "🧶", wantErr: true},
		{field: RigField, value: "Non‑breaking‑hyphen", wantErr: true},
	}
	for _, tc := range tests {
		testValidator(t, tc, "ValidateMultilineString")
	}
}

func TestValidateIntlString(t *testing.T) {
	tests := []validateTest{
		{field: CommentIntlField, value: "१٢௩Ⅳ໖⁷၈🄊〸", wantErr: false},
		{field: CountryIntlField, value: "🇭🇰", wantErr: false},
		{field: MyAntennaIntlField, value: "null\0000character", wantErr: false},
		{field: MyCityIntlField, value: "مكة المكرمة, Makkah the Honored", wantErr: false},
		{field: MyCountryIntlField, value: "", wantErr: false},
		{field: MyNameIntlField, value: "zero​width	space", wantErr: false},
		{field: CommentIntlField, value: "new\nline", wantErr: true},
		{field: CountryIntlField, value: "carriage\rreturn", wantErr: true},
		{field: MyAntennaIntlField, value: "blank\r\n\r\nline", wantErr: true},
		{field: MyCityIntlField, value: "line end\n", wantErr: true},
		{field: MyCountryIntlField, value: "\r\n", wantErr: true},
	}
	for _, tc := range tests {
		testValidator(t, tc, "ValidateIntlString")
	}
}

func TestValidateIntlMultilineString(t *testing.T) {
	tests := []validateTest{
		{field: AddressIntlField, value: "१٢௩Ⅳ໖⁷၈🄊〸", wantErr: false},
		{field: NotesIntlField, value: "🇭🇰", wantErr: false},
		{field: QslmsgIntlField, value: "null\0000character", wantErr: false},
		{field: RigIntlField, value: "مكة المكرمة, Makkah the Honored", wantErr: false},
		{field: AddressIntlField, value: "", wantErr: false},
		{field: NotesIntlField, value: "zero​width	space", wantErr: false},
		{field: QslmsgIntlField, value: "new\nline", wantErr: false},
		{field: RigIntlField, value: "carriage\rreturn", wantErr: false},
		{field: AddressIntlField, value: "blank\r\n\r\nline", wantErr: false},
		{field: NotesIntlField, value: "line end\n", wantErr: false},
		{field: QslmsgIntlField, value: "\r\n", wantErr: false},
		{field: RigIntlField, value: "hello\tworld\r\n", wantErr: false},
	}
	for _, tc := range tests {
		testValidator(t, tc, "ValidateIntlMultilineString")
	}
}

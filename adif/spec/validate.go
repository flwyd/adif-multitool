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
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var allNumeric = regexp.MustCompile("^[0-9]*$")

type ValidationContext struct{}

type FieldValidator func(value string, f Field, ctx ValidationContext) error

var TypeValidators = map[string]FieldValidator{
	"Boolean":                  ValidateBoolean,
	"Character":                ValidateCharacter,
	"Enumeration":              ValidateEnumeration,
	"IntlCharacter":            ValidateIntlCharacter,
	"Date":                     ValidateDate,
	"Digit":                    ValidateDigit,
	"Integer":                  ValidateNumber,
	"IntlString":               ValidateIntlString,
	"IntlMultilineString":      ValidateIntlString,
	"MultilineString":          ValidateString,
	"Number":                   ValidateNumber,
	"PositiveInteger":          ValidateNumber,
	"String":                   ValidateString,
	"Time":                     ValidateTime,
	"AwardList":                ValidateNoop, // TODO
	"CreditList":               ValidateNoop, // TODO
	"GridSquare":               ValidateNoop, // TODO
	"GridSquareExt":            ValidateNoop, // TODO
	"GridSquareList":           ValidateNoop, // TODO
	"IOTARefNo":                ValidateNoop, // TODO
	"Location":                 ValidateNoop, // TODO
	"POTARef":                  ValidateNoop, // TODO
	"POTARefList":              ValidateNoop, // TODO
	"SOTARef":                  ValidateNoop, // TODO
	"SecondarySubdivisionList": ValidateNoop, // TODO
	"SponsoredAwardList":       ValidateNoop, // TODO
	"WWFFRef":                  ValidateNoop, // TODO
}

func ValidateNoop(value string, f Field, ctx ValidationContext) error { return nil }

func ValidateBoolean(val string, f Field, ctx ValidationContext) error {
	switch val {
	case "Y", "N", "y", "n":
		return nil
	default:
		return fmt.Errorf("%s invalid boolean %q", f.Name, val)
	}
}

func ValidateCharacter(val string, f Field, ctx ValidationContext) error {
	if len(val) != 1 {
		return fmt.Errorf("%s not a single character %q", f.Name, val)
	}
	if !isASCIIChar(rune(val[0])) {
		return fmt.Errorf("%s not a printable ASCII character %q", f.Name, val)
	}
	return nil
}

func ValidateString(val string, f Field, ctx ValidationContext) error {
	for _, c := range val {
		if c == '\n' || c == '\r' {
			if !strings.Contains(f.Type.Name, "Multiline") {
				return fmt.Errorf("%s contains newlines but is not a MultilineString %q", f.Name, val)
			}
		} else if !isASCIIChar(c) {
			return fmt.Errorf("%s not a printable ASCII string %q", f.Name, val)
		}
	}
	return nil
}

func ValidateIntlCharacter(val string, f Field, ctx ValidationContext) error {
	if utf8.RuneCountInString(val) != 1 {
		return fmt.Errorf("%s not a single character %q", f.Name, val)
	}
	c, _ := utf8.DecodeRuneInString(val)
	if c == utf8.RuneError {
		return fmt.Errorf("%s invalid Unicode encoding %q", f.Name, val)
	}
	if c == '\n' || c == '\r' {
		return fmt.Errorf("%s line break not allowed in Character fields %q", f.Name, val)
	}
	return nil
}

func ValidateIntlString(val string, f Field, ctx ValidationContext) error {
	for _, c := range val {
		if c == '\n' || c == '\r' {
			if !strings.Contains(f.Type.Name, "Multiline") {
				return fmt.Errorf("%s contains newlines but is not a MultilineIntlString %q", f.Name, val)
			}
		}
	}
	// TODO investigate whether control characters and other Unicode characters are intentionally allowed
	return nil
}

func ValidateDigit(val string, f Field, ctx ValidationContext) error {
	if len(val) != 1 {
		return fmt.Errorf("%s not a single digit %q", f.Name, val)
	}
	d := val[0]
	if d < '0' || d > '9' {
		return fmt.Errorf("%s not an ASCII digit %q", f.Name, val)
	}
	return nil
}

func ValidateNumber(val string, f Field, ctx ValidationContext) error {
	var num float64
	if strings.IndexRune(val, '.') >= 0 {
		if strings.Contains(f.Type.Name, "Integer") {
			return fmt.Errorf("%s invalid integer %q", f.Name, val)
		}
		n, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return fmt.Errorf("%s invalid decimal %q: %v", f.Name, val, err)
		}
		num = n
	} else {
		n, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return fmt.Errorf("%s invalid integer %q: %v", f.Name, val, err)
		}
		num = float64(n)
	}
	minstr := f.Minimum
	if minstr == "" {
		minstr = f.Type.Minimum
	}
	if minstr != "" || f.Type.Name == "PositiveInteger" {
		min, err := strconv.ParseFloat(minstr, 64)
		if err != nil {
			return fmt.Errorf("specification error! invalid minimum %q in field %v: %v", minstr, f, err)
		}
		if num < min {
			return fmt.Errorf("%s value %s below minimum %s", f.Name, val, minstr)
		}
	}
	maxstr := f.Maximum
	if maxstr == "" {
		maxstr = f.Type.Maximum
	}
	if maxstr != "" {
		max, err := strconv.ParseFloat(maxstr, 64)
		if err != nil {
			return fmt.Errorf("specification error! invalid maximum %q in field %v: %v", maxstr, f, err)
		}
		if num > max {
			return fmt.Errorf("%s value %s above maximum %s", f.Name, val, maxstr)
		}
	}
	return nil
}

func ValidateDate(val string, f Field, ctx ValidationContext) error {
	if !allNumeric.MatchString(val) {
		return fmt.Errorf("%s invalid date %q", f.Name, val)
	}
	if len(val) != 8 {
		return fmt.Errorf("%s not an 8-digit date %q", f.Name, val)
	}
	return nil
}

func ValidateTime(val string, f Field, ctx ValidationContext) error {
	if !allNumeric.MatchString(val) {
		return fmt.Errorf("%s invalid time %q", f.Name, val)
	}
	switch len(val) {
	case 4:
		_, err := time.Parse("1504", val)
		if err != nil {
			return fmt.Errorf("%s time out of HH:MM range %q", f.Name, val)
		}
	case 6:
		_, err := time.Parse("150406", val)
		if err != nil {
			return fmt.Errorf("%s time out of HH:MM:SS range %q", f.Name, val)
		}
	default:
		return fmt.Errorf("%s not an 4- or 6-digit time %q", f.Name, val)
	}
	return nil
}

func ValidateEnumeration(val string, f Field, ctx ValidationContext) error {
	if val == "" {
		return nil
	}
	e := f.Enum()
	if e.Name == "" {
		if f.Name == "DARC_DOK" {
			// ADIF spec says type is Enumeration but doesn't provide an enum
			// See https://www.darc.de/der-club/referate/conteste/wag-contest/en/service/districtsdoks/
			// It's supposed to be LNN (L = letter N = number) but special occasions
			// allow values like ILERA and 50ESKILSTUNA so just accept anything
			return ValidateString(val, f, ctx)
		}
		return fmt.Errorf("%s unknown enumeration %q", f.Name, f.EnumName)
	}
	vals := e.Value(val)
	if len(vals) == 0 {
		if e.Name == "Secondary_Administrative_Subdivision" {
			// TODO add ValidateCounties; ADIF spec only lists Alaska values but has
			// references to formats and other sources in III.B.12
			return nil
		}
		return fmt.Errorf("%s unknown value %q for enumeration %s", f.Name, val, e.Name)
	}
	// TODO check EnumScope
	return nil
}

func isASCIIChar(c rune) bool { return c >= 32 && c <= 126 }

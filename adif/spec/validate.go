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
	"unicode"
	"unicode/utf8"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

type Validity int

const (
	Valid Validity = iota
	InvalidError
	InvalidWarning
)

func (v Validity) String() string {
	switch v {
	case Valid:
		return "Valid"
	case InvalidError:
		return "InvalidError"
	case InvalidWarning:
		return "InvalidWarning"
	}
	panic(v)
}

type Validation struct {
	Validity
	Message string
}

func (v Validation) String() string { return v.Message }

func valid() Validation { return Validation{Validity: Valid} }
func errorf(format string, a ...any) Validation {
	return Validation{Validity: InvalidError, Message: fmt.Sprintf(format, a...)}
}
func warningf(format string, a ...any) Validation {
	return Validation{Validity: InvalidWarning, Message: fmt.Sprintf(format, a...)}
}

var (
	allNumeric  = regexp.MustCompile("^[0-9]*$")
	locationPat = regexp.MustCompile(`^(?i)([NESW])([0-9]{3}) ([0-9]{2}\.[0-9]{3})$`)
	// POTA reference A-B@C-D. A is the program identifier (generally an ITU callsign prefix).
	// B is a 4- or 5-digit park number, zero padded.  @C-D is an optional location qualifier.
	// ADIF says C-D follows ISO 3166-2 which would be 2-3 letters, but POTA has
	// some nonstandard country codes: F (France), G (Great Britain), and 9M (Malaysia).
	// There are also non-standard subdivisions, including long English county abbreviations.
	potaPat = regexp.MustCompile(
		"^(?i)[A-Z0-9]{1,4}-[0-9]{4,5}(?:@(?:[A-Z]{2}|F|G|9M)-[A-Z0-9]{2,9})?$")
	sotaPat = regexp.MustCompile("^(?i)[A-Z0-9]{1,4}/[A-Z]{2}-[0-9]{3}$")
	iotaPat = regexp.MustCompile("^(?i)(AF|AN|AS|EU|NA|OC|SA)-[0-9]{3}$")
	wwffPat = regexp.MustCompile("^(?i)[A-Z0-9]{1,4}FF-[0-9]{4}$")
)

type ValidationContext struct {
	UnknownEnumValueWarning bool      // if true, values not in an enumeration are a warning, otherwise an error
	Now                     time.Time // comparison point for times-in-the-future checks
	FieldValue              func(name string) string
}

type FieldValidator func(value string, f Field, ctx ValidationContext) Validation

var TypeValidators = map[string]FieldValidator{
	"Boolean":                  ValidateBoolean,
	"Character":                ValidateCharacter,
	"Enumeration":              ValidateEnumeration,
	"IntlCharacter":            ValidateIntlCharacter,
	"Date":                     ValidateDate,
	"Digit":                    ValidateDigit,
	"GridSquare":               gridsquareValidator(8),
	"GridSquareExt":            gridsquareValidator(4),
	"GridSquareList":           listValidator(gridsquareValidator(8)),
	"Integer":                  ValidateNumber,
	"IntlString":               ValidateIntlString,
	"IntlMultilineString":      ValidateIntlString,
	"IOTARefNo":                formatValidator("IOTA reference", iotaPat),
	"Location":                 ValidateLocation,
	"MultilineString":          ValidateString,
	"Number":                   ValidateNumber,
	"POTARef":                  formatValidator("POTA reference", potaPat),
	"POTARefList":              listValidator(formatValidator("POTA reference", potaPat)),
	"PositiveInteger":          ValidateNumber,
	"SOTARef":                  formatValidator("SOTA reference", sotaPat),
	"String":                   ValidateString,
	"Time":                     ValidateTime,
	"WWFFRef":                  formatValidator("WWFF reference", wwffPat),
	"AwardList":                ValidateNoop, // TODO
	"CreditList":               ValidateNoop, // TODO
	"SecondarySubdivisionList": ValidateNoop, // TODO
	"SponsoredAwardList":       ValidateNoop, // TODO
}

func ValidateNoop(value string, f Field, ctx ValidationContext) Validation { return valid() }

func ValidateBoolean(val string, f Field, ctx ValidationContext) Validation {
	switch val {
	case "Y", "N", "y", "n":
		return valid()
	default:
		return errorf("%s invalid boolean %q", f.Name, val)
	}
}

func ValidateCharacter(val string, f Field, ctx ValidationContext) Validation {
	if len(val) != 1 {
		return errorf("%s not a single character %q", f.Name, val)
	}
	if !isASCIIChar(rune(val[0])) {
		return errorf("%s not a printable ASCII character %q", f.Name, val)
	}
	return valid()
}

func ValidateString(val string, f Field, ctx ValidationContext) Validation {
	for _, c := range val {
		if c == '\n' || c == '\r' {
			if !strings.Contains(f.Type.Name, "Multiline") {
				return errorf("%s contains newlines but is not a MultilineString %q", f.Name, val)
			}
		} else if !isASCIIChar(c) {
			return errorf("%s not a printable ASCII string %q", f.Name, val)
		}
	}
	if f.EnumName != "" {
		// CONTEST_ID and SUBMODE are string fiields with an enumeration; mismatches are warnings not errors
		ctx.UnknownEnumValueWarning = true
		return ValidateEnumeration(val, f, ctx)
	}
	return valid()
}

func ValidateIntlCharacter(val string, f Field, ctx ValidationContext) Validation {
	if utf8.RuneCountInString(val) != 1 {
		return errorf("%s not a single character %q", f.Name, val)
	}
	c, _ := utf8.DecodeRuneInString(val)
	if c == utf8.RuneError {
		return errorf("%s invalid Unicode encoding %q", f.Name, val)
	}
	if c == '\n' || c == '\r' {
		return errorf("%s line break not allowed in Character fields %q", f.Name, val)
	}
	return valid()
}

func ValidateIntlString(val string, f Field, ctx ValidationContext) Validation {
	for _, c := range val {
		if c == '\n' || c == '\r' {
			if !strings.Contains(f.Type.Name, "Multiline") {
				return errorf("%s contains newlines but is not a MultilineIntlString %q", f.Name, val)
			}
		}
	}
	// TODO investigate whether control characters and other Unicode characters are intentionally allowed
	if f.EnumScope != "" {
		return ValidateEnumScope(val, f, ctx)
	}
	return valid()
}

func ValidateDigit(val string, f Field, ctx ValidationContext) Validation {
	if len(val) != 1 {
		return errorf("%s not a single digit %q", f.Name, val)
	}
	d := val[0]
	if !between(d, '0', '9') {
		return errorf("%s not an ASCII digit %q", f.Name, val)
	}
	return valid()
}

func ValidateNumber(val string, f Field, ctx ValidationContext) Validation {
	if val == "" {
		return errorf("%s empty number", f.Name)
	}
	for _, r := range val {
		if !between(r, '0', '9') && (r != '-' && r != '.') {
			// ADIF spec doesn't allow +123 or 1.23e7 as numbers
			return errorf("%s invalid number %q", f.Name, val)
		}
	}
	var num float64
	if strings.IndexRune(val, '.') >= 0 {
		if strings.Contains(f.Type.Name, "Integer") {
			return errorf("%s invalid integer %q", f.Name, val)
		}
		n, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return errorf("%s invalid decimal %q: %v", f.Name, val, err)
		}
		num = n
	} else {
		n, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return errorf("%s invalid integer %q: %v", f.Name, val, err)
		}
		num = float64(n)
	}
	minstr := f.Minimum
	if minstr == "" {
		minstr = f.Type.Minimum
	}
	if minstr != "" || f.Type.Name == PositiveIntegerDataType.Name {
		min, err := strconv.ParseFloat(minstr, 64)
		if err != nil {
			return errorf("specification error! invalid minimum %q in field %v: %v", minstr, f, err)
		}
		if num < min {
			return errorf("%s value %s below minimum %s", f.Name, val, minstr)
		}
	}
	maxstr := f.Maximum
	if maxstr == "" {
		maxstr = f.Type.Maximum
	}
	if maxstr != "" {
		max, err := strconv.ParseFloat(maxstr, 64)
		if err != nil {
			return errorf("specification error! invalid maximum %q in field %v: %v", maxstr, f, err)
		}
		if num > max {
			return errorf("%s value %s above maximum %s", f.Name, val, maxstr)
		}
	}

	iscq := f.Name == CqzField.Name || f.Name == MyCqZoneField.Name
	isitu := f.Name == ItuzField.Name || f.Name == MyItuZoneField.Name
	if iscq || isitu {
		var dxcc string
		if f.Name == CqzField.Name || f.Name == ItuzField.Name {
			dxcc = ctx.FieldValue(DxccField.Name)
			if dxcc == "" {
				dxcc = ctx.FieldValue(CountryField.Name)
			}
		} else if f.Name == MyCqZoneField.Name || f.Name == MyItuZoneField.Name {
			dxcc = ctx.FieldValue(MyDxccField.Name)
			if dxcc == "" {
				dxcc = ctx.FieldValue(MyCountryField.Name)
			}
		}
		if dxcc != "" {
			var zs []int
			if iscq {
				zs = CQZoneFor(dxcc)
			} else if isitu {
				zs = ITUZoneFor(dxcc)
			}
			if len(zs) > 0 && !slices.Contains(zs, int(num)) {
				return errorf("%s value %s is not valid for DXCC %s, not in %v", f.Name, val, dxcc, zs)
			}
		}
		// TODO if len(zs) > 1 check (My)StateField too, but ADIF CQ Zone data on
		// PrimaryAdministrativeSubdivision is incomplete for China/Australia/Yemen
		// and incomplete/inaccurate for several ITU zones
	}
	return valid()
}

func ValidateDate(val string, f Field, ctx ValidationContext) Validation {
	if !allNumeric.MatchString(val) {
		return errorf("%s invalid date %q", f.Name, val)
	}
	if len(val) != 8 {
		return errorf("%s not an 8-digit date %q", f.Name, val)
	}
	d, err := time.Parse("20060102", val)
	if err != nil {
		return errorf("%s invalid date %q", f.Name, val)
	}
	if d.Year() < 1930 {
		return errorf("%s year before 1930 %q", f.Name, val)
	}
	if !ctx.Now.IsZero() && d.After(ctx.Now) {
		return warningf("%s value %q later than today", f.Name, val)
	}
	return valid()
}

func ValidateTime(val string, f Field, ctx ValidationContext) Validation {
	if !allNumeric.MatchString(val) {
		return errorf("%s invalid time %q", f.Name, val)
	}
	var dateField, d string
	if f.Name == TimeOnField.Name {
		dateField = QsoDateField.Name
	} else if f.Name == TimeOffField.Name {
		dateField = QsoDateOffField.Name
	}
	if dateField != "" {
		d = ctx.FieldValue(dateField)
	}
	var dtfmt = "20060102"
	switch len(val) {
	case 4:
		_, err := time.Parse("1504", val)
		if err != nil {
			return errorf("%s time out of HH:MM range %q", f.Name, val)
		}
		dtfmt += "1504"
	case 6:
		_, err := time.Parse("150405", val)
		if err != nil {
			return errorf("%s time out of HH:MM:SS range %q", f.Name, val)
		}
		dtfmt += "150405"
	default:
		return errorf("%s not a 4- or 6-digit time %q", f.Name, val)
	}
	if d != "" && !ctx.Now.IsZero() {
		t, err := time.ParseInLocation(dtfmt, d+val, time.UTC)
		if err == nil && t.After(ctx.Now) && t.Truncate(24*time.Hour) == ctx.Now.Truncate(24*time.Hour) {
			return warningf("%s time %q is later than now, %s=%q", f.Name, val, dateField, d)
		}
	}
	return valid()
}

func ValidateLocation(val string, f Field, ctx ValidationContext) Validation {
	if val == "" {
		return valid()
	}
	g := locationPat.FindStringSubmatch(val)
	if g == nil {
		return errorf("%s invalid location format, make sure to zero-pad %q", f.Name, val)
	}
	var lat, lon bool
	var other string
	switch f.Name {
	case LatField.Name:
		lat = true
		other = LonField.Name
	case LonField.Name:
		lon = true
		other = LatField.Name
	case MyLatField.Name:
		lat = true
		other = MyLonField.Name
	case MyLonField.Name:
		lon = true
		other = MyLatField.Name
	}
	var deg int64
	var err error
	if deg, err = strconv.ParseInt(g[2], 10, 64); err != nil {
		return errorf("%s non-numeric degrees in %q", f.Name, val)
	} else if !between(deg, 0, 180) || lat && !between(deg, 0, 90) {
		return errorf("%s degrees out of range in %q", f.Name, val)
	}
	if min, err := strconv.ParseFloat(g[3], 10); err != nil {
		return errorf("%s non-numeric minutes in %q", f.Name, val)
	} else if !between(min, 0.0, 60.0) {
		return errorf("%s minutes out of range in %q", f.Name, val)
	} else if min > 0.0 && (deg == 180 && lon || deg == 90 && lat) {
		return errorf("%s minutes out of range in %q", f.Name, val)
	}
	if lat && !strings.ContainsAny(g[1], "NSns") {
		return errorf("%s should be north or south, not %q", f.Name, val)
	}
	if lon && !strings.ContainsAny(g[1], "EWew") {
		return errorf("%s should be east or west, not %q", f.Name, val)
	}
	if other != "" && ctx.FieldValue(other) == "" {
		return warningf("%s is set but %s is not set; latitude and longitude usually come together", f.Name, other)
	}
	return valid()
}

func ValidateEnumeration(val string, f Field, ctx ValidationContext) Validation {
	if val == "" {
		return valid()
	}
	e := f.Enum()
	if e.Name == "" {
		if f.Name == DarcDokField.Name || f.Name == MyDarcDokField.Name {
			// ADIF spec says type is Enumeration but doesn't provide an enum
			// See https://www.darc.de/der-club/referate/conteste/wag-contest/en/service/districtsdoks/
			// It's supposed to be LNN (L = letter N = number) but special occasions
			// allow values like ILERA and 50ESKILSTUNA so just accept anything
			// TODO Ask adifdev why this enumeration isn't in the spec
			return ValidateString(val, f, ctx)
		}
		return errorf("%s unknown enumeration %q", f.Name, f.EnumName)
	}
	if f.EnumScope != "" {
		return ValidateEnumScope(val, f, ctx)
	}
	vals := e.Value(val)
	if len(vals) == 0 {
		fn := errorf
		if ctx.UnknownEnumValueWarning {
			fn = warningf
		}
		return fn("%s unknown value %q for enumeration %s", f.Name, val, e.Name)
	}
	if f.Name == ContField.Name {
		if d := ctx.FieldValue(DxccField.Name); d != "" {
			if c := ContinentFor(d); !strings.EqualFold(val, c.Abbreviation) {
				return warningf("continent %s does not match DXCC %s continent %s", val, d, c.Abbreviation)
			}
		}
		if d := ctx.FieldValue(CountryField.Name); d != "" {
			if c := ContinentFor(d); !strings.EqualFold(val, c.Abbreviation) {
				return warningf("continent %s does not match country %s continent %s", val, d, c.Abbreviation)
			}
		}
	}
	return valid()
}

func ValidateEnumScope(val string, f Field, ctx ValidationContext) Validation {
	if val == "" || f.EnumScope == "" {
		return valid()
	}
	e := f.Enum()
	var sval string
	if ctx.FieldValue != nil {
		sval = ctx.FieldValue(f.EnumScope)
	}
	if sval == "" {
		vals := e.Value(val)
		if len(vals) == 0 {
			if e.Name == SecondaryAdministrativeSubdivisionEnumeration.Name {
				// TODO add ValidateCounties; ADIF spec only lists Alaska values but has
				// references to formats and other sources in III.B.12
				return valid()
			}
			return warningf("%s unknown value %q for enumeration %s, %s is not set", f.Name, val, e.Name, f.EnumScope)
		}
		return warningf("%s has value %q but %s is not set", f.Name, val, f.EnumScope)
	}
	var match bool
	prop := e.ScopeProperty()
	if prop == "" {
		return warningf("%s config error! %s doesn't have a ScopeProperty", f.Name, e.Name)
	}
	svals := e.ScopeValues(sval)
	if len(svals) == 0 {
		if e.Name == SecondaryAdministrativeSubdivisionEnumeration.Name {
			// This enum only has values defined for Alaska
			return valid()
		}
		return warningf("%s has value %q but %s doesn't define any values for %s=%q", f.Name, val, e.Name, f.EnumScope, sval)
	}
	for _, v := range svals {
		if strings.EqualFold(val, v.String()) {
			match = true
			break
		}
	}
	if !match {
		fn := errorf
		if e.Name == SubmodeEnumeration.Name {
			// ADIF spec suggests "use enumeration values for interoperability"
			// but a new submode might've been developed since ADIF spec
			fn = warningf
		}
		return fn("%s value %q is not valid for %s=%q", f.Name, val, f.EnumScope, sval)
	}
	return valid()
}

func gridsquareValidator(maxLen int) FieldValidator {
	return func(val string, f Field, ctx ValidationContext) Validation {
		if val == "" {
			return valid()
		}
		if len(val) > maxLen {
			return errorf("%s grid square too long (max %d) %q", f.Name, maxLen, val)
		}
		if len(val)%2 != 0 {
			return errorf("%s odd grid square length %q", f.Name, val)
		}
		// Maidenhead locator alternates two letters, two digits
		for i, c := range val {
			if !isASCIIChar(c) {
				return errorf("%s invalid grid square %q", f.Name, val)
			}
			switch i % 4 {
			case 0, 1:
				if !unicode.IsLetter(c) {
					return errorf("%s non-letter in position %d %q", f.Name, i, val)
				}
			case 2, 3:
				if !unicode.IsNumber(c) {
					return errorf("%s non-digit in position %d %q", f.Name, i, val)
				}
			}
		}
		return valid()
	}
}

func listValidator(fv FieldValidator) FieldValidator {
	return func(val string, f Field, ctx ValidationContext) Validation {
		if val == "" {
			return valid()
		}
		for _, v := range strings.Split(val, ",") {
			if res := fv(v, f, ctx); res.Validity != Valid {
				return res
			}
		}
		return valid()
	}
}

func formatValidator(name string, p *regexp.Regexp) FieldValidator {
	return func(val string, f Field, ctx ValidationContext) Validation {
		if val == "" {
			return valid()
		}
		if !p.MatchString(val) {
			return errorf("%s invalid %s format %q", f.Name, name, val)
		}
		return valid()
	}
}

func isASCIIChar(c rune) bool { return between(c, 32, 126) }

func between[T constraints.Ordered](val, low, high T) bool {
	return val >= low && val <= high
}

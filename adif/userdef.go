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
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type UserdefField struct {
	Name       string
	Type       DataType
	Min, Max   float64
	EnumValues []string
}

func (u UserdefField) String() string {
	var sb strings.Builder
	sb.WriteString(strings.ToUpper(u.Name))
	if u.Type != TypeUnspecified {
		sb.WriteRune(':')
		sb.WriteString(u.Type.Indicator())
	}
	if u.Min != 0.0 || u.Max != 0.0 {
		sb.WriteString(fmt.Sprintf("={%v,%v}", u.Min, u.Max))
	}
	if len(u.EnumValues) > 0 {
		sb.WriteString(fmt.Sprintf("={%s}", strings.Join(u.EnumValues, ",")))
	}
	return sb.String()
}

func (u UserdefField) ValidateSelf() error {
	// ADIF spec says user-defined field can't ccontain comma, colon, or
	// open/close angle/curly bracket and can't start/end with a space.
	// This allows for some pretty weird fields names like ")*(" and "1 - .2 = &"
	// Userdef is defined as a String field, thogh, so unicode and control
	// characters aren't allowed in userdef field names.
	if u.Name == "" {
		return errors.New("empty USERDEF field name")
	}
	if u.Name[0] == ' ' || u.Name[len(u.Name)-1] == ' ' {
		return fmt.Errorf("USERDEF fields cannot begin or end with space: %q", u.Name)
	}
	if strings.HasPrefix(strings.ToUpper(u.Name), "APP_") {
		return fmt.Errorf("USERDEF fields cannot begin with %q: %q", "APP_", u.Name)
	}
	invalid := func(r rune) bool {
		switch r {
		case ',', ':', '<', '>', '{', '}':
			return true
		}
		return r < 32 || r > 126
	}
	if i := strings.IndexFunc(u.Name, invalid); i != -1 {
		return fmt.Errorf("invalid characters in USERDEF field %q", u.Name)
	}
	if math.IsNaN(u.Min) || math.IsNaN(u.Max) || u.Min > u.Max {
		return fmt.Errorf("invalid range {%f:%f} for USERDEF field %q", u.Min, u.Max, u.Name)
	}
	for _, v := range u.EnumValues {
		enumInvalid := func(r rune) bool {
			return r < 32 || r > 126 || r == ',' || r == '{' || r == '}'
		}
		if i := strings.IndexFunc(v, enumInvalid); i != -1 {
			return fmt.Errorf("invalid enum value %q for USERDEF field %q", v, u.Name)
		}
	}
	return nil
}

func (u UserdefField) Validate(f Field) error {
	if u.Type == TypeNumber && (u.Min != 0.0 || u.Max != 0.0) {
		v, err := strconv.ParseFloat(f.Value, 64)
		if err != nil {
			return fmt.Errorf("%s: number parse error: %w", u.Name, err)
		}
		if v < u.Min || v > u.Max {
			return fmt.Errorf("%s: %v is not in range {%v,%v}", u.Name, v, u.Min, u.Max)
		}
	}
	if len(u.EnumValues) > 0 {
		var match bool
		for _, e := range u.EnumValues {
			if strings.EqualFold(e, f.Value) {
				match = true
				break
			}
		}
		if !match {
			return fmt.Errorf("%s: %q is not in enumeration {%s}", u.Name, f.Value, strings.Join(u.EnumValues, ","))
		}
	}
	return nil
}

// formatRange returns the "{Min:Max}" format used by ADI and ADX.
func formatRange(u UserdefField) string {
	// fmt doesn't seem to have a "no decimal point if it's an integer" option
	return fmt.Sprintf(
		"{%s:%s}",
		strconv.FormatFloat(u.Min, 'f', -1, 64),
		strconv.FormatFloat(u.Max, 'f', -1, 64),
	)
}

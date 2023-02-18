// Copyright 2022 Google LLC
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
)

type DataType int

const (
	TypeUnspecified DataType = iota
	TypeMultilineString
	TypeIntlMultilineString
	TypeString
	TypeIntlString
	TypeEnumeration
	TypeBoolean
	TypeNumber
	TypeDate
	TypeTime
	TypeLocation
	// TODO add types without data type indicators? e.g. Character, Digit, GridSquare
)

var typeIndicators = map[string]DataType{
	"":  TypeUnspecified,
	"M": TypeMultilineString,
	"G": TypeIntlMultilineString,
	"S": TypeString,
	"I": TypeIntlString,
	"E": TypeEnumeration,
	"B": TypeBoolean,
	"N": TypeNumber,
	"D": TypeDate,
	"T": TypeTime,
	"L": TypeLocation,
}

func DataTypeFromIndicator(id string) (DataType, error) {
	t, ok := typeIndicators[strings.ToUpper(id)]
	if !ok {
		return TypeUnspecified, fmt.Errorf("unknown data type indicator %s", id)
	}
	return t, nil
}

func (t DataType) Indicator() string {
	switch t {
	case TypeUnspecified:
		return ""
	case TypeMultilineString:
		return "M"
	case TypeIntlMultilineString:
		return "G"
	case TypeString:
		return "S"
	case TypeIntlString:
		return "I"
	case TypeEnumeration:
		return "E"
	case TypeBoolean:
		return "B"
	case TypeNumber:
		return "N"
	case TypeDate:
		return "D"
	case TypeTime:
		return "T"
	case TypeLocation:
		return "L"
	default:
		panic(fmt.Sprintf("unknown DataType %d", t))
	}
}

func (t DataType) String() string {
	return t.Indicator() // TODO use go generate
}

type Field struct {
	Name  string
	Value string
	Type  DataType
}

func (f Field) String() string {
	if f.Type == TypeUnspecified {
		return fmt.Sprintf("%s=%s", strings.ToUpper(f.Name), f.Value)
	}
	return fmt.Sprintf("%s:%s=%s", strings.ToUpper(f.Name), f.Type.Indicator(), f.Value)
}

func (f Field) IsAppDefined() bool {
	return strings.HasPrefix(strings.ToUpper(f.Name), "APP_")
}

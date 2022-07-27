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

// CONSIDER an enum prefix or a subpackage
type DataType int

const (
	Unspecified DataType = iota
	MultilineString
	IntlMultilineString
	String
	IntlString
	Enumeration
	Boolean
	Number
	Date
	Time
	Location
	// TODO add types without data type identifiers? e.g. Character, Digit, GridSquare
)

var typeIdentifiers = map[string]DataType{
	"":  Unspecified,
	"M": MultilineString, "G": IntlMultilineString,
	"S": String, "I": IntlString,
	"E": Enumeration, "B": Boolean, "N": Number,
	"D": Date, "T": Time,
	"L": Location}

func DataTypeFromIdentifier(id string) (DataType, error) {
	t, ok := typeIdentifiers[strings.ToUpper(id)]
	if !ok {
		return Unspecified, fmt.Errorf("unknown data type identifier %s", id)
	}
	return t, nil
}

func (t DataType) Identifier() string {
	switch t {
	case Unspecified:
		return ""
	case MultilineString:
		return "M"
	case IntlMultilineString:
		return "G"
	case String:
		return "S"
	case IntlString:
		return "I"
	case Enumeration:
		return "E"
	case Boolean:
		return "B"
	case Number:
		return "N"
	case Date:
		return "D"
	case Time:
		return "T"
	case Location:
		return "L"
	default:
		panic(fmt.Sprintf("unknown DataType %d", t))
	}
}

func (t DataType) String() string {
	return t.Identifier() // TODO use go generate
}

type Field struct {
	Name  string
	Value string
	Type  DataType
}

func (f Field) String() string {
	if f.Type == Unspecified {
		return fmt.Sprintf("%s=%s", strings.ToUpper(f.Name), f.Value)
	}
	return fmt.Sprintf("%s:%s=%s", strings.ToUpper(f.Name), f.Type.Identifier(), f.Value)
}

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
	"strings"
)

type EnumValue interface {
	Property(string) string
	fmt.Stringer
}

type Enumeration struct {
	Name       string
	Properties []string
	Values     []EnumValue
}

func (e Enumeration) String() string { return e.Name }

func (e Enumeration) Value(val string) []EnumValue {
	res := make([]EnumValue, 0, 2)
	for _, v := range e.Values {
		// Band is lower case, most others upper
		if strings.EqualFold(val, v.String()) {
			res = append(res, v)
		}
	}
	return res
}

func (e Enumeration) ScopeValues(scope string) []EnumValue {
	p := e.ScopeProperty()
	if p == "" {
		return []EnumValue{}
	}
	res := make([]EnumValue, 0)
	for _, v := range e.Values {
		if strings.EqualFold(scope, v.Property(p)) {
			res = append(res, v)
		}
	}
	return res
}

func (e Enumeration) ScopeProperty() string {
	switch e.Name {
	case "Primary_Administrative_Subdivision":
		return "DXCC Entity Code"
	case "Secondary_Administrative_Subdivision":
		return "DXCC Entity Code"
	case "Submode":
		return "Mode"
	default:
		return ""
	}
}

var Enumerations = make(map[string]Enumeration)

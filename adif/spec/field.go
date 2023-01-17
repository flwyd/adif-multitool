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

type Field struct {
	Name       string
	Type       DataType
	EnumName   string
	EnumScope  string
	Header     bool
	Minimum    string
	Maximum    string
	ImportOnly bool
}

func (f Field) Enum() Enumeration {
	if f.EnumName == "" {
		return Enumeration{}
	}
	return Enumerations[f.EnumName]
}

var Fields = make(map[string]Field)

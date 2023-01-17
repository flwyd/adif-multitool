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
	"testing"
)

func TestFieldsDeclared(t *testing.T) {
	if len(Fields) == 0 {
		t.Errorf("len(Fields) == 0, make sure to run go generate aidf/spec")
	}
}

func TestFieldTypes(t *testing.T) {
	for name, f := range Fields {
		if name != f.Name {
			t.Errorf("Fields[%q].Name got %q, want %q", name, f.Name, name)
		}
		if f.Type.Name == "" {
			t.Errorf("Fields[%q] has empty data type %v", name, f.Type)
		}
	}
}

func TestFieldEnums(t *testing.T) {
	var failed bool
	for name, f := range Fields {
		if f.EnumName != "" {
			if got := f.Enum(); got.Name != f.EnumName {
				t.Errorf("Fields[%q].Enum() got %q, want name %q", name, got, f.EnumName)
				failed = true
			}
		}
	}
	if failed {
		var sb strings.Builder
		for name, e := range Enumerations {
			sb.WriteString(fmt.Sprintf("%q = %v\n", name, e))
		}
		t.Errorf("Enumerations list:\n%s", sb.String())
	}
}

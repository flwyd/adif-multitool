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

func TestEnumerationDeclared(t *testing.T) {
	if len(Enumerations) == 0 {
		t.Errorf("len(Enumerations) == 0, make sure to run go generate aidf/spec")
	}
}

func TestEnumerationNames(t *testing.T) {
	for name, e := range Enumerations {
		if name != e.Name {
			t.Errorf("Enumerations[%q].Name got %q, want %q", name, e.Name, name)
		}
	}
}

func TestEnumerationValues(t *testing.T) {
	for name, e := range Enumerations {
		seen := make(map[string]bool)
		for _, v := range e.Values {
			prop := "Enumeration Name"
			if got := v.Property(prop); got != name {
				t.Errorf("%s %s.Property(%q) got %q, want %q", e, v, prop, got, name)
			}
			// Primary_Administrative_Subdivision has lots of duplicate abbreviations
			// Country has a couple entities that were deleted and re-added
			// Region has three KO Kosovo entries for different time ranges
			dupsok := map[string]bool{"Country": true, "Region": true, "Primary_Administrative_Subdivision": true}
			if seen[v.String()] && !dupsok[e.Name] {
				t.Errorf("%s has multiple enum values named %q", e, v)
			}
			seen[v.String()] = true
		}
	}
}

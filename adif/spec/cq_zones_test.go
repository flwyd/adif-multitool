// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spec

import "testing"

func TestAllCountriesCQZone(t *testing.T) {
	var active, inactive []CountryEnum
	for _, e := range CountryEnumeration.Values {
		c := e.(CountryEnum)
		if c == CountryNone {
			continue
		}
		if c.Deleted == "true" {
			inactive = append(inactive, c)
		} else {
			active = append(active, c)
		}
	}
	tests := []struct {
		name      string
		countries []CountryEnum
	}{
		{name: "Active", countries: active},
		{name: "Inactive", countries: inactive},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for _, c := range tc.countries {
				if len(CQZoneFor(c.EntityName)) == 0 {
					t.Errorf("CQZoneFor(%q) is missing", c.EntityName)
				}
				if len(CQZoneFor(c.EntityCode)) == 0 {
					t.Errorf("CQZoneFor(%q) is missing", c.EntityCode)
				}
			}
		})
	}
}

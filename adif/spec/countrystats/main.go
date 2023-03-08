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

// countrystats prints miscellaneous information about DXCC entities and
// ISO 3166-1 codes, particularly highlighting potential problems.
// Output format and contents are subject to change.
package main

import (
	"fmt"

	"github.com/flwyd/adif-multitool/adif/spec"
)

func main() {
	entitiesWithISO := make(map[spec.CountryEnum]int)
	fmt.Println("ISO 3166-1 codes associated with a deleted DXCC entity:")
	for _, c := range spec.ISO3166Countries {
		for _, d := range c.DXCC {
			entitiesWithISO[d] += 1
			if d.Deleted == "true" {
				fmt.Printf("%s/%s %s : %s %s\n", c.Alpha2, c.Alpha3, c.EnglishName, d.EntityCode, d.EntityName)
			}
		}
	}
	fmt.Println()
	fmt.Println("ISO 3166-1 codes associated with multiple DXCC entities:")
	for _, c := range spec.ISO3166Countries {
		switch len(c.DXCC) {
		case 1:
		case 0:
			fmt.Printf("%s/%s %s (no DXCC association)\n", c.Alpha2, c.Alpha3, c.EnglishName)
		default:
			for _, d := range c.DXCC {
				fmt.Printf("%s/%s %s : %s %s\n", c.Alpha2, c.Alpha3, c.EnglishName, d.EntityCode, d.EntityName)
			}
		}
	}
	fmt.Println()
	fmt.Println("ISO 3166-1 codes missing from a lookup:")
	for _, c := range spec.ISO3166Countries {
		if spec.ISO3166Alpha[c.Alpha2].EnglishName != c.EnglishName {
			fmt.Printf("%s %s\n", c.Alpha2, c.EnglishName)
		}
		if spec.ISO3166Alpha[c.Alpha3].EnglishName != c.EnglishName {
			fmt.Printf("%s %s\n", c.Alpha3, c.EnglishName)
		}
	}
	fmt.Println()
	fmt.Println("ISO 3166-1 codes sharing the name of a DXCC entity:")
	for _, c := range spec.ISO3166Countries {
		for _, a := range []string{c.Alpha2, c.Alpha3} {
			if e := spec.CountryEnumeration.Value(a); len(e) != 0 {
				fmt.Printf("%s %s : %v\n", a, c.EnglishName, e)
			}
		}
	}
	fmt.Println()
	fmt.Println("DXCC entities without an ISO code:")
	for _, e := range spec.CountryEnumeration.Values {
		c := e.(spec.CountryEnum)
		if entitiesWithISO[c] == 0 && c.Deleted != "true" {
			fmt.Printf("%s %s\n", c.EntityCode, c.EntityName)
		}
	}
	fmt.Println()
	fmt.Println("DXCC entities with more than one ISO code:")
	for c, n := range entitiesWithISO {
		if n > 1 {
			fmt.Printf("%s %s\n", c.EntityCode, c.EntityName)
		}
	}
}

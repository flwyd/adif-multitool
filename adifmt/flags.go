// Copyright 2024 Google LLC
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

package main

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

type mapValue struct {
	m       map[string]string
	key     string
	allowed []string
}

func (v *mapValue) String() string {
	return v.m[v.key]
}

func (v *mapValue) Set(s string) error {
	val := strings.ToUpper(s)
	if !slices.Contains(v.allowed, val) {
		return fmt.Errorf("%q is not in %s", val, strings.Join(v.allowed, ", "))
	}
	v.m[v.key] = val
	return nil
}

func (v *mapValue) Get() any {
	return v.String()
}

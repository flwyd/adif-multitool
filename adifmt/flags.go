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
	"unicode/utf8"

	"golang.org/x/exp/slices"
	"golang.org/x/text/language"
)

type runeValue struct{ r *rune }

func (v runeValue) String() string {
	if v.r == nil {
		return ""
	}
	return fmt.Sprintf("%q", *v.r)
}

func (v runeValue) Set(s string) error {
	switch utf8.RuneCountInString(s) {
	case 0:
		return nil
	case 1:
		r, _ := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError {
			return fmt.Errorf("invalid UTF-8 encoding %q", s)
		}
		*v.r = r
		return nil
	default:
		return fmt.Errorf("expecting one character, not %q", s)
	}
}

func (v runeValue) Get() rune { return *v.r }

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

type languageValue struct{ *language.Tag }

func (v *languageValue) Set(s string) error {
	t, err := language.Parse(s)
	if err != nil {
		return err
	}
	*v.Tag = t
	return nil
}

func (v *languageValue) Get() any { return v.Tag }

func (v *languageValue) String() string {
	if v.Tag == nil {
		return language.Und.String()
	}
	return v.Tag.String()
}

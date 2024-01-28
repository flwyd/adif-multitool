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

package adif

import (
	"fmt"
	"strings"
)

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func ensureCRLF(s string) string {
	first := -1
	var wasCR bool
	for i, r := range s {
		if wasCR && r != '\n' {
			first = i - 1
			break
		}
		if r == '\n' && !wasCR {
			first = i
			break
		}
		wasCR = r == '\r'
	}
	if wasCR && first < 0 {
		first = len(s) - 1
	}
	if first < 0 && !wasCR {
		return s // don't rebuild the string if it doesn't need changes
	}
	var res strings.Builder
	res.WriteString(s[0:first])
	wasCR = false
	for _, r := range s[first:] {
		if wasCR && r != '\n' {
			res.WriteRune('\n')
		} else if !wasCR && r == '\n' {
			res.WriteRune('\r')
		}
		res.WriteRune(r)
		wasCR = r == '\r'
	}
	if wasCR {
		res.WriteRune('\n')
	}
	return res.String()
}

func fieldValues(l *Logfile, fname string) map[string]int {
	res := make(map[string]int)
	for _, r := range l.Records {
		if f, ok := r.Get(fname); ok {
			res[f.Value]++
		}
	}
	return res
}

func isAllDigits(s string) bool {
	// TODO when upgrading to Go 1.20+, use
	// return !strings.ContainsFunc(s, func(r rune) bool { return r < '0' || r > '9' })
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// cutPrefix is a polyfill for strings.CutPrefix which was introduced in Go 1.20.
func cutPrefix(s, prefix string) (after string, found bool) {
	// TODO when upgrading to Go 1.20+ remove this polyfill
	if strings.HasPrefix(s, prefix) {
		var before string
		before, after, found = strings.Cut(s, prefix)
		if !found || before != "" {
			panic(fmt.Sprintf("strings.Cut(%q, %q) got %q, %q, %v", s, prefix, before, after, found))
		}
		return
	}
	return s, false
}

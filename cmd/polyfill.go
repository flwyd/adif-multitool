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

package cmd

import (
	"fmt"
	"strings"
)

// errorsJoin is a polyfill for Go's errors.Join function (added in 1.20).
func errorsJoin(errs ...error) error {
	f := make([]string, 0, len(errs))
	e := make([]any, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			f = append(f, "%w")
			e = append(e, err)
		}
	}
	if len(e) == 0 {
		return nil
	}
	return fmt.Errorf(strings.Join(f, "\n"), e...)
}

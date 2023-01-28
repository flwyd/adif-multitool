// Copyright 2022 Google LLC
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

	"github.com/flwyd/adif-multitool/adif"
)

var Cat = Command{Name: "cat", Run: runCat,
	Description: "Concatenate all input files to standard output"}

func runCat(ctx *Context, args []string) error {
	// TODO add any needed flags
	srcs := argSources(ctx, args...)
	out := adif.NewLogfile()
	for _, f := range srcs {
		l, err := readSource(ctx, f)
		if err != nil {
			return fmt.Errorf("error reading %s: %v", f, err)
		}
		updateFieldOrder(out, l.FieldOrder)
		// TODO merge headers and comments
		out.Records = append(out.Records, l.Records...)
	}
	return write(ctx, out)
}

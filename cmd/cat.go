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
	"github.com/flwyd/adif-multitool/adif"
)

var Cat = Command{Name: "cat", Run: runCat,
	Description: "Concatenate all input files to standard output"}

func runCat(ctx *Context, args []string) error {
	// TODO add any needed flags
	acc := accumulator{Out: adif.NewLogfile(), Ctx: ctx}
	for _, f := range filesOrStdin(args) {
		l, err := acc.read(f)
		if err != nil {
			return err
		}
		updateFieldOrder(acc.Out, l.FieldOrder)
		// TODO merge headers and comments
		acc.Out.Records = append(acc.Out.Records, l.Records...)
	}
	if err := acc.prepare(); err != nil {
		return err
	}
	return write(ctx, acc.Out)
}

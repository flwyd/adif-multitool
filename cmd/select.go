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

	"github.com/flwyd/adif-multitool/adif"
)

var Select = Command{Name: "select", Run: runSelect,
	Description: "Print only specific fields from the input; skip records with no matching fields"}

type SelectContext struct {
	Fields FieldList
}

func runSelect(ctx *Context, args []string) error {
	con := ctx.CommandCtx.(*SelectContext)
	if len(con.Fields) == 0 {
		return fmt.Errorf("no fields provided, try %s select -fields CALL,BAND", ctx.ProgramName)
	}
	out := adif.NewLogfile()
	out.FieldOrder = con.Fields
	for _, f := range filesOrStdin(args) {
		l, err := readFile(ctx, f)
		if err != nil {
			return err
		}
		// TODO merge headers and comments
		for _, r := range l.Records {
			fields := make([]adif.Field, 0, len(con.Fields))
			for _, name := range con.Fields {
				if f, ok := r.Get(name); ok {
					fields = append(fields, f)
				}
			}
			if len(fields) > 0 {
				out.Records = append(out.Records, adif.NewRecord(fields...))
			}
		}
	}
	return write(ctx, out)
}

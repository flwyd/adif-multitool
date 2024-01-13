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
	"os"
	"path/filepath"

	"github.com/flwyd/adif-multitool/adif"
)

var Select = Command{Name: "select", Run: runSelect, Help: helpSelect,
	Description: "Print only specific fields from the input"}

type SelectContext struct {
	Fields FieldList
}

func helpSelect() string {
	return "Records with no matching fields will be skipped in the output.\n"
}

func runSelect(ctx *Context, args []string) error {
	con := ctx.CommandCtx.(*SelectContext)
	if len(con.Fields) == 0 {
		return fmt.Errorf("no fields provided, try %s select -fields CALL,BAND", filepath.Base(os.Args[0]))
	}
	out := adif.NewLogfile()
	out.FieldOrder = con.Fields
	acc := accumulator{Out: out, Ctx: ctx}
	for _, f := range filesOrStdin(args) {
		l, err := acc.read(f)
		if err != nil {
			return err
		}
		for _, r := range l.Records {
			fields := make([]adif.Field, 0, len(con.Fields))
			for _, name := range con.Fields {
				if f, ok := r.Get(name); ok {
					fields = append(fields, f)
				}
			}
			if len(fields) > 0 {
				out.AddRecord(adif.NewRecord(fields...))
			}
		}
	}
	if err := acc.prepare(); err != nil {
		return err
	}
	return write(ctx, out)
}

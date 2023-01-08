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
	"flag"
	"fmt"

	"github.com/flwyd/adif-multitool/adif"
)

var Select = Command{Name: "select", Run: runSelect, AddFlags: selectFlags,
	Description: "Print only specific fields from the input; skip records with no matching fields"}

type selectContext struct {
	fields fieldList
}

func selectFlags(ctx *Context, fs *flag.FlagSet) {
	con := selectContext{fields: make(fieldList, 0, 16)}
	fs.Var(&con.fields, "fields", "Comma-separated or multiple instance field names to include in output")
	ctx.CommandCtx = &con
}

func runSelect(ctx *Context, args []string) error {
	con := ctx.CommandCtx.(*selectContext)
	srcs := argSources(args...)
	out := adif.NewLogfile("")
	for _, f := range srcs {
		l, err := readSource(ctx, f)
		if err != nil {
			return fmt.Errorf("error reading %s: %v", f, err)
		}
		// TODO merge headers and comments
		for _, r := range l.Records {
			fields := make([]adif.Field, 0, len(con.fields))
			for _, name := range con.fields {
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

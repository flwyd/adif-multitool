// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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

var Edit = Command{Name: "edit", Run: runEdit, AddFlags: editFlags,
	Description: "Add, remove, or change field values"}

type editContext struct {
	add         fieldAssignments
	set         fieldAssignments
	remove      fieldList
	removeBlank bool
}

func editFlags(ctx *Context, fs *flag.FlagSet) {
	cctx := editContext{
		add:    newFieldAssignments(validateAlphanumName),
		set:    newFieldAssignments(validateAlphanumName),
		remove: make(fieldList, 0)}
	fs.Var(&cctx.add, "add", "Add `field=value` if field is not already in a record (repeatable)")
	fs.Var(&cctx.set, "set", "Set `field=value` for all records (repeatable)")
	fs.Var(&cctx.remove, "remove", "Remove fields from all records (comma-separated, repeatable)")
	fs.BoolVar(&cctx.removeBlank, "remove-blank", false, "Remove all blank fields")
	ctx.CommandCtx = &cctx
}

func runEdit(ctx *Context, args []string) error {
	cctx := ctx.CommandCtx.(*editContext)
	remove := make(map[string]bool)
	for _, n := range cctx.remove {
		remove[n] = true
	}
	set := make(map[string]adif.Field)
	for _, f := range cctx.set.values {
		if remove[f.Name] {
			return fmt.Errorf("%q in both -set and -remove", f.Name)
		}
		set[f.Name] = f
	}
	for _, f := range cctx.add.values {
		if remove[f.Name] {
			return fmt.Errorf("%q in both -add and -remove, use -set to change values", f.Name)
		}
		if _, ok := set[f.Name]; ok {
			return fmt.Errorf("%q in both -set and -add", f.Name)
		}
	}
	srcs := argSources(ctx, args...)
	out := adif.NewLogfile("")
	for _, f := range srcs {
		l, err := readSource(ctx, f)
		if err != nil {
			return fmt.Errorf("error reading %s: %v", f, err)
		}
		updateFieldOrder(out, l.FieldOrder)
		// TODO merge headers and comments
		for _, r := range l.Records {
			seen := make(map[string]bool)
			old := r.Fields()
			fields := make([]adif.Field, 0, len(old))
			for _, f := range old {
				if remove[f.Name] {
					continue
				}
				if cctx.removeBlank && f.Value == "" {
					continue
				}
				seen[f.Name] = true
				if v, ok := set[f.Name]; ok {
					f = v
				}
				fields = append(fields, f)
			}
			for _, f := range cctx.set.values {
				if !seen[f.Name] {
					fields = append(fields, f)
				}
				seen[f.Name] = true
			}
			for _, f := range cctx.add.values {
				if !seen[f.Name] {
					fields = append(fields, f)
				}
				seen[f.Name] = true
			}
			if len(fields) > 0 {
				out.Records = append(out.Records, adif.NewRecord(fields...))
			}
		}
	}
	return write(ctx, out)
}

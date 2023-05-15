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
	"errors"
	"sort"
	"strings"

	"github.com/flwyd/adif-multitool/adif"
	"github.com/flwyd/adif-multitool/adif/spec"
)

var Sort = Command{Name: "sort", Run: runSort, Help: helpSort,
	Description: "Sort records by a list of fields"}

type SortContext struct {
	Fields FieldList
}

func helpSort() string {
	return "Prefix a field with - for descending order, e.g. -FREQ,-QSO_DATE\n"
}

func runSort(ctx *Context, args []string) error {
	cctx := ctx.CommandCtx.(*SortContext)
	comps := make([]spec.FieldComparator, len(cctx.Fields))
	mults := make([]int, len(cctx.Fields))
	fields := make([]string, len(cctx.Fields))
	for i, n := range cctx.Fields {
		if n == "" {
			return errors.New("empty field name")
		}
		if n[0] == '-' {
			mults[i] = -1
			n = n[1:]
		} else {
			mults[i] = 1
		}
		fields[i] = n
		n = strings.ToUpper(n)
		if f, ok := spec.Fields[n]; ok {
			comps[i] = spec.ComparatorForField(f, ctx.Locale)
		} // else resolve dynamically
	}
	out := adif.NewLogfile()
	acc := accumulator{Out: out, Ctx: ctx}
	for _, f := range filesOrStdin(args) {
		l, err := acc.read(f)
		if err != nil {
			return err
		}
		for _, r := range l.Records {
			out.AddRecord(r)
		}
	}
	if err := acc.prepare(); err != nil {
		return err
	}
	sort.SliceStable(out.Records, func(i, j int) bool {
		a := out.Records[i]
		b := out.Records[j]
		for k, comp := range comps {
			if comp == nil {
				if uf, _ := out.GetUserdef(fields[k]); uf.Type.Indicator() != "" {
					f := spec.Field{Name: fields[k], Type: spec.DataTypes[uf.Type.Indicator()]}
					comp = spec.ComparatorForField(f, ctx.Locale)
				} else {
					var t spec.DataType
					fa, _ := a.Get(fields[k])
					fb, _ := b.Get(fields[k])
					if fa.Type == fb.Type {
						if fa.Type == adif.TypeUnspecified {
							t = spec.StringDataType
						} else {
							t = spec.DataTypes[fa.Type.Indicator()]
						}
					} else if fa.Type == adif.TypeUnspecified {
						t = spec.DataTypes[fb.Type.Indicator()]
					} else if fb.Type == adif.TypeUnspecified {
						t = spec.DataTypes[fa.Type.Indicator()]
					} else {
						t = spec.StringDataType // type confusion, sort as strings
					}
					f := spec.Field{Name: fields[k], Type: t}
					comp = spec.ComparatorForField(f, ctx.Locale)
				}
			}
			af, _ := a.Get(fields[k])
			bf, _ := b.Get(fields[k])
			c, err := comp(af.Value, bf.Value)
			if err != nil {
				return false // if can't compare, treat as equal
			}
			c *= mults[k]
			if c < 0 {
				return true
			}
			if c > 0 {
				return false
			}
		}
		return false
	})
	return write(ctx, out)
}

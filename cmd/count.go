// Copyright 2025 Google LLC
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

package cmd

import (
	"sort"
	"strconv"
	"strings"

	"github.com/flwyd/adif-multitool/adif"
	"github.com/flwyd/adif-multitool/adif/spec"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

var Count = Command{Name: "count", Run: runCount, Help: helpCount,
	Description: "Count records or unique field combinations"}

type CountContext struct {
	CountFieldName string
	Fields         FieldList
}

func helpCount() string {
	return `If no fields are specified:
  Outputs a single record with a single field with the number of records in all
  input files.
If fields are specified:
  Outputs each unique combination of those fields with the number of times the
  combination occurs in the input records.  Record order is unspecified.
`
}

func runCount(ctx *Context, args []string) error {
	cctx := ctx.CommandCtx.(*CountContext)
	countName := cctx.CountFieldName
	if countName == "" {
		countName = "COUNT"
	}
	names := make([]string, len(cctx.Fields))
	for i, f := range cctx.Fields {
		names[i] = strings.ToUpper(f)
	}
	acc, err := newAccumulator(ctx)
	if err != nil {
		return err
	}
	all := make([]*adif.Record, 0, 128)
	for _, file := range filesOrStdin(args) {
		l, err := acc.read(file)
		if err != nil {
			return err
		}
		for _, r := range l.Records {
			cr := adif.NewRecord()
			for _, n := range cctx.Fields {
				f, _ := r.Get(n)
				cr.Set(f)
			}
			all = append(all, cr)
		}
	}
	if len(all) == 0 {
		r := adif.NewRecord(adif.Field{Name: countName, Value: "0", Type: adif.TypeNumber})
		for _, n := range cctx.Fields {
			r.Set(adif.Field{Name: n, Value: ""})
		}
		acc.Out.AddRecord(r)
		if err := acc.prepare(); err != nil {
			return err
		}
		return write(ctx, acc.Out)
	}
	comps := make([]spec.FieldComparator, len(cctx.Fields))
	for i, n := range cctx.Fields {
		if f, ok := spec.FieldNamed(n); ok {
			comps[i] = spec.ComparatorForField(f, ctx.Locale)
		} else if u, ok := acc.Out.GetUserdef(n); ok {
			f := spec.Field{Name: n, Type: spec.DataTypes[u.Type.Indicator()]}
			comps[i] = spec.ComparatorForField(f, ctx.Locale)
		} else {
			comps[i] = spec.ComparatorForField(spec.Field{Name: n, Type: spec.StringDataType}, ctx.Locale)
		}
	}
	col := collate.New(language.Und, collate.IgnoreCase)
	comp := func(a, b *adif.Record) int {
		for i, c := range comps {
			n := cctx.Fields[i]
			af, _ := a.Get(n)
			bf, _ := b.Get(n)
			v, err := c(af.Value, bf.Value)
			if err != nil {
				v = col.CompareString(af.Value, bf.Value)
			}
			if v != 0 {
				return v
			}
		}
		return 0
	}
	sort.SliceStable(all, func(i, j int) bool { return comp(all[i], all[j]) < 0 })
	for i := 0; i < len(all); {
		cur := all[i]
		vals := make([]map[adif.Field]int, len(cctx.Fields))
		for j, n := range cctx.Fields {
			vals[j] = make(map[adif.Field]int)
			if f, ok := cur.Get(n); ok {
				vals[j][f] = 1
			}
		}
		num := 1
		i++
		for i < len(all) {
			if comp(cur, all[i]) != 0 {
				break
			}
			for j, n := range cctx.Fields {
				if f, ok := all[i].Get(n); ok {
					vals[j][f]++
				}
			}
			num++
			i++
		}
		r := adif.NewRecord(adif.Field{Name: countName, Value: strconv.Itoa(num), Type: adif.TypeNumber})
		for j, m := range vals {
			if len(m) == 0 {
				r.Set(adif.Field{Name: cctx.Fields[j], Value: ""})
			} else {
				r.Set(mostFrequent(m))
			}
		}
		acc.Out.AddRecord(r)
	}
	if err := acc.prepare(); err != nil {
		return err
	}
	return write(ctx, acc.Out)
}

func mostFrequent(counts map[adif.Field]int) adif.Field {
	var f adif.Field
	var m int
	for k, v := range counts {
		if v > m || (m == v && k.Value > f.Value) {
			f = k
			m = v
		}
	}
	return f
}

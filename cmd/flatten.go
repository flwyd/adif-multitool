// Copyright 2024 Google LLC
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
	"errors"
	"fmt"
	"strings"

	"github.com/flwyd/adif-multitool/adif"
	"github.com/flwyd/adif-multitool/adif/spec"
)

var Flatten = Command{Name: "flatten", Run: runFlatten, Help: helpFlatten,
	Description: "Flatten multi-instance fields to multiple records"}

type FlattenContext struct {
	Fields     FieldList
	Delimiters FieldDelimiters
}

func helpFlatten() string {
	return "If multiple fields are given, a Cartesian combination will be output.\n"
}

func runFlatten(ctx *Context, args []string) error {
	cctx := ctx.CommandCtx.(*FlattenContext)
	delims := make(map[string]string)
	for name, delim := range cctx.Delimiters {
		delims[name] = delim
	}
	for _, n := range cctx.Fields {
		if n == "" {
			return errors.New("empty field name")
		}
		n = strings.ToUpper(n)
		if delims[n] != "" {
			continue
		}
		f, ok := spec.Fields[n]
		if !ok {
			return fmt.Errorf("unknown field %q", n)
		}
		d := typeDelims[f.Type]
		if d == "" {
			return fmt.Errorf("don't know delimiter for field %q of type %s", n, f.Type.Name)
		}
		delims[n] = d
	}

	out := adif.NewLogfile()
	acc := accumulator{Out: out, Ctx: ctx}
	for _, f := range filesOrStdin(args) {
		l, err := acc.read(f)
		if err != nil {
			return err
		}
		for _, r := range l.Records {
			expn := []*adif.Record{r}
			for _, n := range cctx.Fields {
				d := delims[strings.ToUpper(n)]
				if f, ok := r.Get(n); ok && f.Value != "" {
					l := strings.Split(f.Value, d)
					if len(l) > 1 {
						more := make([]*adif.Record, 0, len(expn)*len(l))
						for _, e := range expn {
							for _, v := range l {
								fv := f
								fv.Value = v
								c := adif.NewRecord(e.Fields()...)
								c.Set(fv)
								c.SetComment(e.GetComment())
								more = append(more, c)
							}
						}
						expn = more
					}
				}
			}
			for _, e := range expn {
				out.AddRecord(e)
			}
		}
	}
	if err := acc.prepare(); err != nil {
		return err
	}
	return write(ctx, out)
}

var typeDelims = map[spec.DataType]string{
	spec.AwardListDataType:                ",",
	spec.CreditListDataType:               ",",
	spec.GridSquareListDataType:           ",",
	spec.POTARefListDataType:              ",",
	spec.SecondarySubdivisionListDataType: ":", // NV,Clark:UT,Washington
	spec.SponsoredAwardListDataType:       ",",
}

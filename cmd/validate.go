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
	"strings"
	"time"

	"github.com/flwyd/adif-multitool/adif"
	"github.com/flwyd/adif-multitool/adif/spec"
)

var Validate = Command{Name: "validate", Run: runValidate, Help: helpValidate,
	Description: "Validate field values; non-zero exit and no stdout if invalid"}

type ValidateContext struct {
	RequiredFields FieldList
	Cond           ConditionValue
}

func helpValidate() string {
	return "Non-failure warnings are added as comments in ADI and ADX output.\n\n" + helpCondition("validate --required-fields=AGE")
}

func runValidate(ctx *Context, args []string) error {
	cctx := ctx.CommandCtx.(*ValidateContext)
	now := time.Now().UTC() // consistent for the whole log
	cond := cctx.Cond.Get()
	log := os.Stderr
	var errors, warnings int
	appFields := make(map[string]adif.DataType)
	acc, err := newAccumulator(ctx)
	if err != nil {
		return err
	}
	for _, f := range filesOrStdin(args) {
		l, err := acc.read(f)
		if err != nil {
			return err
		}
		updateFieldOrder(acc.Out, l.FieldOrder)
		for i, r := range l.Records {
			vctx := spec.ValidationContext{
				Now: now,
				FieldValue: func(name string) string {
					f, _ := r.Get(name)
					return f.Value
				}}
			var msgs []string
			if cond.Evaluate(recordEvalContext{record: r, lang: ctx.Locale}) {
				missing := make([]string, 0)
				for _, x := range cctx.RequiredFields {
					if f, ok := r.Get(x); !ok || f.Value == "" {
						missing = append(missing, x)
					}
				}
				if len(missing) > 0 {
					errors++
					fmt.Fprintf(log, "ERROR on %s record %d: missing fields %s\n", l, i+1, strings.Join(missing, ", "))
				}
			}
			for _, f := range r.Fields() {
				name := strings.ToUpper(f.Name)
				if f.IsAppDefined() {
					if adt := appFields[name]; adt == adif.TypeUnspecified {
						appFields[name] = f.Type
					} else if f.Type != adif.TypeUnspecified && f.Type != adt {
						warnings++
						fmt.Fprintf(log, "WARNING on %s record %d: inconsistent types for %s\n", l, i+1, f.Name)
					}
				}
				if f.Value == "" {
					continue
				}
				validateSpec := func(fv spec.FieldValidator, fs spec.Field) {
					if fv != nil {
						switch v := fv(f.Value, fs, vctx); v.Validity {
						case spec.InvalidError:
							errors++
							fmt.Fprintf(log, "ERROR on %s record %d: %s\n", l, i+1, v)
						case spec.InvalidWarning:
							warnings++
							fmt.Fprintf(log, "WARNING on %s record %d: %s\n", l, i+1, v)
							msgs = append(msgs, fmt.Sprintf("%s: %s", f.Name, v.Message))
						}
					}
				}
				if fs, ok := spec.FieldNamed(f.Name); ok {
					validateSpec(spec.TypeValidators[fs.Type.Name], fs)
				} else if u, ok := acc.Out.GetUserdef(f.Name); ok {
					if len(u.EnumValues) > 0 || u.Min != 0.0 || u.Max != 0.0 {
						if err := u.Validate(f); err != nil {
							errors++
							fmt.Fprintf(log, "ERROR on %s record %d: %s\n", l, i+1, err)
						}
					} else { // spec enum validator can't handle userdef enums
						dt := spec.DataTypes[u.Type.Indicator()]
						fs := spec.Field{Name: u.Name, Type: dt}
						validateSpec(spec.TypeValidators[dt.Name], fs)
					}
				} else if f.IsAppDefined() {
					fs := spec.Field{Name: f.Name, Type: spec.DataTypes[appFields[name].Indicator()]}
					validateSpec(spec.TypeValidators[fs.Type.Name], fs)
				}
				if len(msgs) > 0 {
					r.SetComment("adif-multitool: validate warnings: " + strings.Join(msgs, "; "))
				}
			}
			acc.Out.AddRecord(r)
		}
	}
	if errors > 0 {
		return fmt.Errorf("validate got %d errors and %d warnings", errors, warnings)
	}
	if err := acc.prepare(); err != nil {
		return err
	}
	err = write(ctx, acc.Out)
	if warnings > 0 {
		fmt.Fprintf(log, "validate got %d warnings\n", warnings)
	}
	return err
}

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

	"github.com/flwyd/adif-multitool/adif"
	"github.com/flwyd/adif-multitool/adif/spec"
)

// TODO write command tests
var Validate = Command{Name: "validate", Run: runValidate,
	Description: "Validate field values; non-zero exit and no stdout if invalid"}

func runValidate(ctx *Context, args []string) error {
	// TODO add any needed flags
	log := os.Stderr
	var errors, warnings int
	appFields := make(map[string]adif.DataType)
	out := adif.NewLogfile()
	acc := accumulator{Out: out, Ctx: ctx}
	for _, f := range filesOrStdin(args) {
		l, err := acc.read(f)
		if err != nil {
			return err
		}
		updateFieldOrder(out, l.FieldOrder)
		for i, r := range l.Records {
			vctx := spec.ValidationContext{}
			var msgs []string
			for _, f := range r.Fields() {
				name := strings.ToUpper(f.Name)
				if f.IsAppDefined() {
					if adt := appFields[name]; adt == adif.TypeUnspecified {
						appFields[name] = f.Type
					} else if f.Type != adif.TypeUnspecified && f.Type != adt {
						warnings++
						fmt.Fprintf(log, "WARNING on %s record %d: inconsistent types for %s\n", l, i, f.Name)
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
							fmt.Fprintf(log, "ERROR on %s record %d: %s\n", l, i, v)
						case spec.InvalidWarning:
							warnings++
							fmt.Fprintf(log, "WARNING on %s record %d: %s\n", l, i, v)
							msgs = append(msgs, fmt.Sprintf("%s: %s", f.Name, v.Message))
						}
					}
				}
				if fs, ok := spec.Fields[f.Name]; ok {
					validateSpec(spec.TypeValidators[fs.Type.Name], fs)
				} else if u, ok := acc.Out.GetUserdef(f.Name); ok {
					if len(u.EnumValues) > 0 || u.Min != 0.0 || u.Max != 0.0 {
						if err := u.Validate(f); err != nil {
							errors++
							fmt.Fprintf(log, "ERROR on %s record %d: %s\n", l, i, err)
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
			out.Records = append(out.Records, r)
		}
	}
	if errors > 0 {
		return fmt.Errorf("validate got %d errors and %d warnings", errors, warnings)
	}
	if err := acc.prepare(); err != nil {
		return err
	}
	err := write(ctx, out)
	if warnings > 0 {
		fmt.Fprintf(log, "validate got %d warnings\n", warnings)
	}
	return err
}

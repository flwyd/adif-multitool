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
	out := adif.NewLogfile()
	var comments commentCatcher
	for _, f := range filesOrStdin(args) {
		l, err := readFile(ctx, f)
		if err != nil {
			return err
		}
		updateFieldOrder(out, l.FieldOrder)
		// TODO merge headers and comments
		for i, r := range l.Records {
			vctx := spec.ValidationContext{}
			var msgs []string
			for _, f := range r.Fields() {
				if f.Value == "" {
					continue
				}
				if fs, ok := spec.Fields[f.Name]; ok {
					if dtv := spec.TypeValidators[fs.Type.Name]; dtv != nil {
						switch v := dtv(f.Value, fs, vctx); v.Validity {
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
				if len(msgs) > 0 {
					r.SetComment("adif-multitool: validate warnings: " + strings.Join(msgs, "; "))
				}
			}
			out.Records = append(out.Records, r)
		}
		comments.read(l, f)
	}
	if errors > 0 {
		return fmt.Errorf("validate got %d errors and %d warnings", errors, warnings)
	}
	comments.write(out)
	err := write(ctx, out)
	if warnings > 0 {
		fmt.Fprintf(log, "validate got %d warnings\n", warnings)
	}
	return err
}

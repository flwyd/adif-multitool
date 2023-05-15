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
	"github.com/flwyd/adif-multitool/adif"
)

var Find = Command{Name: "find", Run: runFind, Help: helpFind,
	Description: "Include only records matching a condition"}

type FindContext struct {
	Cond ConditionValue
}

func helpFind() string {
	return `Condition syntax and examples:
  field = value : Case-insensitive equality, contest_id=ARRL-field-day
  field < value : Less than, freq<29.701
  field <= value : Less than or equal, band<=10m
  field > value : Greater than, tx_pwr>100
  field >= value : Greater than or equal, qso_date>=20200101

Fields can be compared to other fields by enclosing in '{' and '}':
  gridsquare={my_gridsquare} : contact in the same maidenhead grid
  freq<{freq_rx} : operating split below other station

Conditions can match multiple values separated by '|' characters:
  mode=SSB|FM|AM|DIGITALVOICE : any phone mode
  arrl_sect={my_arrl_sect}|ENY|NLI|NNY|WNY : In same section or New York

Conditions match list fields if any value matches:
  pota_ref=K-4556 : Matches "K-0034,K-4556"

Empty or absent fields can be matched by omitting value:
  operator= : OPERATOR field not set
  my_sig_info> : MY_SIG_INFO field is set ("greater than empty")

Use quotes so operators are not treated as special shell characters:
  find --if 'freq>=7' --if-not 'mode=CW' --or-if 'tx_pwr<=5'
`
}

func runFind(ctx *Context, args []string) error {
	cctx := ctx.CommandCtx.(*FindContext)
	cond := cctx.Cond.Get()
	out := adif.NewLogfile()
	acc := accumulator{Out: out, Ctx: ctx}
	for _, f := range filesOrStdin(args) {
		l, err := acc.read(f)
		if err != nil {
			return err
		}
		updateFieldOrder(out, l.FieldOrder)
		for _, r := range l.Records {
			eval := recordEvalContext{record: r, lang: ctx.Locale}
			if cond.Evaluate(eval) {
				out.AddRecord(r)
			}
		}
	}
	if err := acc.prepare(); err != nil {
		return err
	}
	return write(ctx, out)
}

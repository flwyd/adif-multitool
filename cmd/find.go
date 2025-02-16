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

var Find = Command{Name: "find", Run: runFind, Help: helpFind,
	Description: "Include only records matching a condition"}

type FindContext struct {
	Cond ConditionValue
}

func helpFind() string {
	return helpCondition("find")
}

func runFind(ctx *Context, args []string) error {
	cctx := ctx.CommandCtx.(*FindContext)
	cond := cctx.Cond.Get()
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
		for _, r := range l.Records {
			eval := recordEvalContext{record: r, lang: ctx.Locale}
			if cond.Evaluate(eval) {
				acc.Out.AddRecord(r)
			}
		}
	}
	if err := acc.prepare(); err != nil {
		return err
	}
	return write(ctx, acc.Out)
}

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

package main

import (
	"flag"

	"github.com/flwyd/adif-multitool/cmd"
)

type cmdConfig struct {
	cmd.Command
	Configure func(*cmd.Context, *flag.FlagSet)
}

var (
	catConf = cmdConfig{Command: cmd.Cat,
		Configure: func(*cmd.Context, *flag.FlagSet) {}}

	editConf = cmdConfig{Command: cmd.Edit,
		Configure: func(ctx *cmd.Context, fs *flag.FlagSet) {
			cctx := cmd.EditContext{
				Add:    cmd.NewFieldAssignments(cmd.ValidateAlphanumName),
				Set:    cmd.NewFieldAssignments(cmd.ValidateAlphanumName),
				Remove: make(cmd.FieldList, 0)}
			fs.Var(&cctx.Add, "add", "Add `field=value` if field is not already in a record (repeatable)")
			fs.Var(&cctx.Set, "set", "Set `field=value` for all records (repeatable)")
			fs.Var(&cctx.Remove, "remove", "Remove `fields` from all records (comma-separated, repeatable)")
			fs.BoolVar(&cctx.RemoveBlank, "remove-blank", false, "Remove all blank fields")
			fs.Var(&cctx.FromZone, "time-zone-from", "Adjust times and dates from this time `zone` into -time-zone-to (default UTC)")
			fs.Var(&cctx.ToZone, "time-zone-to", "Adjust times and dates into this time `zone` from -time-zone-from (default UTC)")
			ctx.CommandCtx = &cctx
		}}

	fixConf = cmdConfig{Command: cmd.Fix,
		Configure: func(*cmd.Context, *flag.FlagSet) {}}

	selectConf = cmdConfig{Command: cmd.Select,
		Configure: func(ctx *cmd.Context, fs *flag.FlagSet) {
			cctx := cmd.SelectContext{Fields: make(cmd.FieldList, 0, 16)}
			fs.Var(&cctx.Fields, "fields", "Comma-separated or multiple instance field names to include in output")
			ctx.CommandCtx = &cctx
		}}

	validateConf = cmdConfig{Command: cmd.Validate,
		Configure: func(*cmd.Context, *flag.FlagSet) {}}

	cmds = []cmdConfig{catConf, editConf, fixConf, selectConf, validateConf}
)

func commandNamed(name string) (cmdConfig, bool) {
	for _, c := range cmds {
		if c.Name == name {
			return c, true
		}
	}
	return cmdConfig{}, false
}

func commandNames() []string {
	res := make([]string, len(cmds))
	for i, c := range cmds {
		res[i] = c.Name
	}
	return res
}

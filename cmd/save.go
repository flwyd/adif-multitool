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

	"github.com/flwyd/adif-multitool/adif"
)

var Save = Command{Name: "save", Run: runSave,
	Description: "Save standard input to file with format inferred by extension"}

type SaveContext struct {
	OverwriteExisting bool
	WriteIfEmpty      bool
}

func runSave(ctx *Context, args []string) error {
	cctx := ctx.CommandCtx.(*SaveContext)
	if len(args) != 1 {
		return fmt.Errorf("save expects 1 output file, got %v", args)
	}
	fname := args[0]
	fs := ctx.fs
	if fs == nil {
		fs = osFilesystem{}
	}
	if !cctx.OverwriteExisting && fs.Exists(fname) {
		return fmt.Errorf("output file %s already exists", fname)
	}
	format := ctx.OutputFormat
	if !format.IsValid() {
		f, err := adif.GuessFormatFromName(fname)
		if err != nil {
			return fmt.Errorf("unknown output format, set -output: %w", err)
		}
		format = f
	}
	l, err := readFile(ctx, os.Stdin.Name())
	if err != nil {
		return err
	}
	if len(l.Records) == 0 {
		if !cctx.WriteIfEmpty {
			return fmt.Errorf("no records in input, not saving to %s", fname)
		}
		fmt.Fprintf(os.Stderr, "Warning: saving %s with no records", fname)
	}
	for _, u := range ctx.UserdefFields {
		l.AddUserdef(u)
	}
	out, err := fs.Create(fname)
	if err != nil {
		return err
	}
	defer out.Close()
	ctx.Out = out
	ctx.OutputFormat = format
	return write(ctx, l)
}

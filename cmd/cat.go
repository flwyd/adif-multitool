// Copyright 2022 Google LLC
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
	"path/filepath"
	"strings"

	"github.com/flwyd/adif-multitool/adif"
)

var Cat = Command{Name: "cat", Run: runCat,
	Description: "Concatenate all input files to standard output"}

func runCat(ctx *Context, args []string) error {
	// TODO add any needed flags
	srcs := argSources(args...)
	out := adif.NewLogfile("")
	for _, f := range srcs {
		src, err := f.Open()
		if err != nil {
			return err
		}
		ext := strings.TrimPrefix(filepath.Ext(src.Name()), ".")
		format, err := adif.ParseFormat(ext)
		if err != nil {
			format = ctx.InputFormat
		}
		r := ctx.Readers[format]
		l, err := r.Read(src)
		if err != nil {
			return fmt.Errorf("error reading %s: %v", f, err)
		}
		// TODO merge headers and comments
		out.Records = append(out.Records, l.Records...)
	}
	ctx.SetHeaders(out)
	w, ok := ctx.Writers[ctx.OutputFormat]
	if !ok {
		return fmt.Errorf("unknown output format %q", ctx.OutputFormat)
	}
	// TODO flag to save to file rather than stdout?
	w.Write(out, os.Stdout)
	return nil
}

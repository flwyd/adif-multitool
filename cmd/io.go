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
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/flwyd/adif-multitool/adif"
)

func write(ctx *Context, l *adif.Logfile) error {
	ctx.SetHeaders(l)
	w, ok := ctx.Writers[ctx.OutputFormat]
	if !ok {
		return fmt.Errorf("unknown output format %q", ctx.OutputFormat)
	}
	w.Write(l, ctx.Out)
	return nil
}

func filesOrStdin(args []string) []string {
	if len(args) == 0 {
		return []string{"-"}
	}
	return args
}

func readFile(ctx *Context, filename string) (*adif.Logfile, error) {
	fs := ctx.fs
	if fs == nil {
		fs = osFilesystem{}
	}
	f, err := fs.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	ext := strings.TrimPrefix(filepath.Ext(f.Name()), ".")
	format, err := adif.ParseFormat(ext)
	if err != nil {
		format = ctx.InputFormat
	}
	r := ctx.Readers[format]
	l, err := r.Read(f)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %w", f.Name(), err)
	}
	l.Filename = f.Name()
	return l, nil
}

// NamedReader is an io.Reader with a name.  os.File implements this interface
// and StringReader is provided for testing.
type NamedReader interface {
	io.ReadCloser
	Name() string
}

type filesystem interface {
	// Open opens a file with the given name with the semantics of os.File.
	Open(name string) (NamedReader, error)
}

type osFilesystem struct{}

func (_ osFilesystem) Open(name string) (NamedReader, error) { return os.Open(name) }

func updateFieldOrder(l *adif.Logfile, fields []string) {
	seen := make(map[string]bool)
	for _, f := range l.FieldOrder {
		seen[strings.ToUpper(f)] = true
	}
	for _, f := range fields {
		n := strings.ToUpper(f)
		if !seen[n] {
			l.FieldOrder = append(l.FieldOrder, f)
			seen[n] = true
		}
	}
}

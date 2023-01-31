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
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/flwyd/adif-multitool/adif"
)

func write(ctx *Context, l *adif.Logfile) error {
	ctx.SetHeaders(l)
	format := ctx.OutputFormat
	if !format.IsValid() {
		format = adif.FormatADI
	}
	w, ok := ctx.Writers[format]
	if !ok {
		return fmt.Errorf("unknown output format %q", format)
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
	ior := bufio.NewReader(f)
	format := ctx.InputFormat
	if !format.IsValid() {
		format, err = adif.GuessFormatFromName(f.Name())
		if err != nil {
			format, err = adif.GuessFormatFromContent(ior)
			if err != nil {
				return nil, fmt.Errorf("could not determine type of %s: %w", f.Name(), err)
			}
		}
	}
	r := ctx.Readers[format]
	l, err := r.Read(ior)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %w", f.Name(), err)
	}
	l.Filename = f.Name()
	return l, nil
}

// NamedReader is an io.Reader with a name.  os.File implements this interface
// and stringReader is provided for testing.
type NamedReader interface {
	io.ReadCloser
	Name() string
}

type filesystem interface {
	// Exists returns true if the named file is known to exist, false otherwise.
	Exists(name string) bool
	// Open opens a file with the given name with the semantics of os.File.
	Open(name string) (NamedReader, error)
	// Create creates a file and opens it for writing, truncating the file if it
	// alrready exists.  See os.Create for more details.
	Create(name string) (io.WriteCloser, error)
}

type osFilesystem struct{}

func (_ osFilesystem) Exists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func (_ osFilesystem) Open(name string) (NamedReader, error) {
	if name == "-" || name == os.Stdin.Name() {
		return os.Stdin, nil
	}
	return os.Open(name)
}

func (_ osFilesystem) Create(name string) (io.WriteCloser, error) { return os.Create(name) }

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

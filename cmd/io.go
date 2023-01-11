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

func write(ctx *Context, l *adif.Logfile) error {
	ctx.SetHeaders(l)
	w, ok := ctx.Writers[ctx.OutputFormat]
	if !ok {
		return fmt.Errorf("unknown output format %q", ctx.OutputFormat)
	}
	w.Write(l, ctx.Out)
	return nil
}

func argSources(ctx *Context, filenames ...string) []argSource {
	fs := ctx.fs
	if fs == nil {
		fs = osFilesystem{}
	}
	if len(filenames) == 0 {
		return []argSource{fs.Lookup("-")}
	}
	s := make([]argSource, len(filenames))
	for i, f := range filenames {
		s[i] = fs.Lookup(f)
	}
	return s
}

func readSource(ctx *Context, f argSource) (*adif.Logfile, error) {
	src, err := f.Open()
	if err != nil {
		return nil, err
	}
	ext := strings.TrimPrefix(filepath.Ext(src.Name()), ".")
	format, err := adif.ParseFormat(ext)
	if err != nil {
		format = ctx.InputFormat
	}
	r := ctx.Readers[format]
	l, err := r.Read(src)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %v", f, err)
	}
	return l, nil
}

type argSource interface{ Open() (adif.Source, error) }

type fileSource struct{ filename string }

func (s fileSource) Open() (adif.Source, error) { return os.Open(s.filename) }

func (s fileSource) String() string { return s.filename }

type stdinSource struct{}

func (s stdinSource) Open() (adif.Source, error) { return os.Stdin, nil }

func (s stdinSource) String() string { return os.Stdin.Name() }

type filesystem interface{ Lookup(name string) argSource }

type osFilesystem struct{}

func (_ osFilesystem) Lookup(name string) argSource {
	if name == "-" || name == os.Stdin.Name() {
		return stdinSource{}
	}
	return fileSource{filename: name}
}

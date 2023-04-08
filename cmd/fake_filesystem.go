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
	"io"
	"strings"
)

// stringReader implements NamedReader with a strings.Reader to aid testing.
type stringReader struct {
	*strings.Reader
	Filename string
}

func (s *stringReader) Close() error { return nil }

func (s *stringReader) Name() string { return s.Filename }

func (s *stringReader) String() string { return s.Filename }

type stringWriter struct {
	strings.Builder
	closeCB func(*stringWriter) error
}

func (w *stringWriter) Close() error { return w.closeCB(w) }

type fakeFilesystem struct{ files map[string]string }

func (fs fakeFilesystem) Exists(name string) bool {
	_, ok := fs.files[name]
	return ok
}

func (fs fakeFilesystem) Open(name string) (NamedReader, error) {
	f, ok := fs.files[name]
	if !ok {
		names := make([]string, 0, len(fs.files))
		for k := range fs.files {
			names = append(names, k)
		}
		return nil, fmt.Errorf("%s does not exist, not one of %v", name, names)
	}
	return &stringReader{Reader: strings.NewReader(f), Filename: name}, nil
}

func (fs fakeFilesystem) Create(name string) (io.WriteCloser, error) {
	fs.files[name] = ""
	return &stringWriter{closeCB: func(w *stringWriter) error {
		fs.files[name] = w.String()
		return nil
	}}, nil
}

func (fs fakeFilesystem) MkdirAll(dir string) error {
	// currently not worying about enforcing directories
	return nil
}

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
	"strings"
)

// StringReader implements NamedReader with a strings.Reader to aid testing.
type StringReader struct {
	*strings.Reader
	Filename string
}

func (s *StringReader) Close() error { return nil }

func (s *StringReader) Name() string { return s.Filename }

func (s *StringReader) String() string { return s.Filename }

type fakeFilesystem struct{ files map[string]string }

func (fs fakeFilesystem) Open(name string) (NamedReader, error) {
	f, ok := fs.files[name]
	if !ok {
		names := make([]string, 0, len(fs.files))
		for k, _ := range fs.files {
			names = append(names, k)
		}
		return nil, fmt.Errorf("%s does not exist, not one of %v", name, names)
	}
	return &StringReader{Reader: strings.NewReader(f), Filename: name}, nil
}

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

type fakeSource struct {
	name   string
	reader *strings.Reader
}

func (s fakeSource) Open() (NamedReader, error) {
	return StringReader{Reader: s.reader, Filename: s.name}, nil
}

func (s fakeSource) String() string { return s.name }

type errorSource struct {
	err error
}

func (s errorSource) Open() (NamedReader, error) { return nil, s.err }

type fakeFilesystem struct{ files map[string]string }

func (fs fakeFilesystem) Lookup(name string) argSource {
	f, ok := fs.files[name]
	if !ok {
		return errorSource{fmt.Errorf("error opening %s", name)}
	}
	return fakeSource{name: name, reader: strings.NewReader(f)}
}

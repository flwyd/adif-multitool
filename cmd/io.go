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
	"os"

	"github.com/flwyd/adif-multitool/adif"
)

func argSources(filenames ...string) []argSource {
	if len(filenames) == 0 {
		return []argSource{stdinSource{}}
	}
	s := make([]argSource, 0, len(filenames))
	for _, f := range filenames {
		s = append(s, fileSource{f})
	}
	return s
}

type argSource interface{ Open() (adif.Source, error) }

type fileSource struct {
	filename string
}

func (s fileSource) Open() (adif.Source, error) {
	return os.Open(s.filename)
}

func (s fileSource) String() string {
	return s.filename
}

type stdinSource struct{}

func (s stdinSource) Open() (adif.Source, error) {
	return os.Stdin, nil
}

func (s stdinSource) String() string {
	return os.Stdin.Name()
}

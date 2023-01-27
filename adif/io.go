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

package adif

import (
	"fmt"
	"io"
	"strings"
)

type NamedReader interface {
	io.Reader
	Name() string
}

type Reader interface {
	Read(NamedReader) (*Logfile, error)
}

type Writer interface {
	Write(*Logfile, io.Writer) error
}

type ReadWriter interface {
	Reader
	Writer
	fmt.Stringer
}

// StringReader implements NamedReader with a strings.Reader to aid testing.
type StringReader struct {
	Reader   *strings.Reader
	Filename string
}

func (s StringReader) Read(p []byte) (int, error) {
	return s.Reader.Read(p)
}

func (s StringReader) Name() string { return s.Filename }

func (s StringReader) String() string { return s.Filename }

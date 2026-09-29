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

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type Reader interface {
	Read(io.Reader) (*Logfile, error)
}

type Writer interface {
	Write(*Logfile, io.Writer) error
}

type ReadWriter interface {
	Reader
	Writer
	fmt.Stringer
}

func utf8NormalizingReader(r io.Reader) io.Reader {
	// TODO Provide a way to specify non-UTF encodings e.g. charmap.Windows1251
	return transform.NewReader(r, unicode.BOMOverride(unicode.UTF8.NewDecoder()))
}

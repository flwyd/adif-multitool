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

//go:generate go run github.com/abice/go-enum -f=$GOFILE --nocase --flag --names
package adif

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// ENUM(ADI, ADX, Cabrillo, CSV, JSON, TSV)
type Format string

// GuessFormatFromName guesses a file's Format based on its extension.
// If filename doesn't match any known format, Format("") and an error are
// returned.
func GuessFormatFromName(filename string) (Format, error) {
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")
	if ext == "" {
		return Format(""), fmt.Errorf("no file extension in %q", filename)
	}
	f, err := ParseFormat(ext)
	if err != nil {
		switch strings.ToLower(ext) {
		case "cbr", "log":
			f, err = FormatCabrillo, nil
		}
	}
	return f, err
}

var (
	// ADI files can start with an arbitrary-length comment
	firstADITagPat = regexp.MustCompile(`^[^<]*<\w+:\d+(:\w)?>`)
	// Require CSV and TSV files to have a header with purely alphanumeric columns
	csvHeaderPat = regexp.MustCompile(`^\w+(,\w+)+[\r\n]`)
	tsvHeaderPat = regexp.MustCompile(`^\w+(\t\w+)+[\r\n]`)
)

const contentPeekSize = 4096

// GuessFormatFromContent inspects the beginning bytes of r to guess which
// Format the data is in.  Returns Format("") and an error  if no heuristic
// matched the content.
func GuessFormatFromContent(r *bufio.Reader) (Format, error) {
	buf, err := r.Peek(contentPeekSize)
	if err != nil && !errors.Is(err, bufio.ErrBufferFull) && !errors.Is(err, io.EOF) {
		return Format(""), err
	}
	start := bytes.TrimLeftFunc(buf, unicode.IsSpace)
	if len(start) == 0 {
		return Format(""), fmt.Errorf("could not determine data format, input is empty")
	}
	if bytes.HasPrefix(start, []byte("<?xml")) || bytes.HasPrefix(start, []byte("<ADX>")) {
		return FormatADX, nil
	}
	if bytes.HasPrefix(start, []byte("START-OF-LOG:")) {
		return FormatCabrillo, nil
	}
	if firstADITagPat.Find(start) != nil {
		return FormatADI, nil
	}
	if start[0] == '{' {
		return FormatJSON, nil
	}
	if csvHeaderPat.Find(start) != nil {
		return FormatCSV, nil
	}
	if tsvHeaderPat.Find(start) != nil {
		return FormatTSV, nil
	}
	return Format(""), fmt.Errorf("could not determine data format, use the -input option")
}

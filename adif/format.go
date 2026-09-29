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

	"golang.org/x/text/transform"
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
	firstADITagPat = regexp.MustCompile(`^[^<]*<(?i:\w+:\d+(:\w)?|EO[HR])>`)
	// Require CSV and TSV files to have a header with purely alphanumeric columns
	csvHeaderPat = regexp.MustCompile(`^\w+(,\w+)+[\r\n]`)
	tsvHeaderPat = regexp.MustCompile(`^\w+(\t\w+)+[\r\n]`)
)

const contentPeekSize = 4096

// GuessFormatFromContent inspects the beginning bytes of r to guess which
// Format the data is in.  Returns Format("") and an error  if no heuristic
// matched the content.
func GuessFormatFromContent(r *bufio.Reader) (Format, error) {
	rawbuf, err := r.Peek(contentPeekSize)
	if err != nil && !errors.Is(err, bufio.ErrBufferFull) && !errors.Is(err, io.EOF) {
		return Format(""), err
	}
	var buf bytes.Buffer
	tr := utf8NormalizingReader(bytes.NewReader(rawbuf))
	if _, err := io.Copy(&buf, tr); err != nil {
		// It's okay if the 4096th byte is in the middle of a multi-byte Unicode sequence
		if !errors.Is(err, transform.ErrShortDst) && !errors.Is(err, transform.ErrShortSrc) && !errors.Is(err, transform.ErrEndOfSpan) {
			return Format(""), err
		}
	}
	start := bytes.TrimLeftFunc(buf.Bytes(), unicode.IsSpace)
	if len(start) == 0 {
		return Format(""), fmt.Errorf("could not determine data format, input is empty: %v", buf)
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
	return Format(""), fmt.Errorf("could not determine data format, use the --input option")
}

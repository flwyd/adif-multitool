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

import "fmt"

// ENUM(empty, space, tab, newline, 2newline, crlf, 2crlf)
type Separator int

func (d Separator) Val() string {
	switch d {
	case SeparatorEmpty:
		return ""
	case SeparatorSpace:
		return " "
	case SeparatorTab:
		return "\t"
	case SeparatorNewline:
		return "\n"
	case Separator2Newline:
		return "\n\n"
	case SeparatorCrlf:
		return "\r\n"
	case Separator2Crlf:
		return "\r\n\r\n"
	}
	panic(fmt.Sprintf("unhandled Separator %s", d))
}

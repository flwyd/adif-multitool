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

package adif

import (
	"testing"
)

func TestEnsureCRLF(t *testing.T) {
	tests := []struct {
		name, from, want string
	}{
		{name: "empty string", from: "", want: ""},
		{name: "no line breaks", from: "Lorem ipsum", want: "Lorem ipsum"},
		{name: "leading CRLF", from: "\r\nLorem ipsum", want: "\r\nLorem ipsum"},
		{name: "leading CR", from: "\rLorem ipsum", want: "\r\nLorem ipsum"},
		{name: "leading LF", from: "\nLorem ipsum", want: "\r\nLorem ipsum"},
		{name: "trailing CRLF", from: "Lorem ipsum\r\n", want: "Lorem ipsum\r\n"},
		{name: "trailing CR", from: "Lorem ipsum\r", want: "Lorem ipsum\r\n"},
		{name: "trailing LF", from: "Lorem ipsum\n", want: "Lorem ipsum\r\n"},
		{name: "CRLF, no trailing", from: "Lorem\r\nipsum", want: "Lorem\r\nipsum"},
		{name: "CRLF multiple present",
			from: "Lorem\r\nipsum\r\nsic dolor amet\r\n",
			want: "Lorem\r\nipsum\r\nsic dolor amet\r\n"},
		{name: "just multiple CR",
			from: "Lorem\ripsum\rsic dolor amet\r",
			want: "Lorem\r\nipsum\r\nsic dolor amet\r\n"},
		{name: "just multiple LF",
			from: "Lorem\nipsum\nsic dolor amet\n",
			want: "Lorem\r\nipsum\r\nsic dolor amet\r\n"},
		{name: "repeated CRLF",
			from: "\r\n\r\nLorem\r\n\r\nipsum\r\n\r\nsic dolor amet\r\n\r\n",
			want: "\r\n\r\nLorem\r\n\r\nipsum\r\n\r\nsic dolor amet\r\n\r\n"},
		{name: "repeated CR",
			from: "\r\rLorem\r\ripsum\r\rsic dolor amet\r\r",
			want: "\r\n\r\nLorem\r\n\r\nipsum\r\n\r\nsic dolor amet\r\n\r\n"},
		{name: "repeated LF",
			from: "\n\nLorem\n\nipsum\n\nsic dolor amet\n\n",
			want: "\r\n\r\nLorem\r\n\r\nipsum\r\n\r\nsic dolor amet\r\n\r\n"},
		{name: "UTF two-byte CRLF",
			from: "Tromsø\r\nÅland\r\nCuraçao\r\n",
			want: "Tromsø\r\nÅland\r\nCuraçao\r\n"},
		{name: "UTF two-byte CR",
			from: "Tromsø\rÅland\rCuraçao\r",
			want: "Tromsø\r\nÅland\r\nCuraçao\r\n"},
		{name: "UTF two-byte LF",
			from: "Tromsø\nÅland\nCuraçao\n",
			want: "Tromsø\r\nÅland\r\nCuraçao\r\n"},
		{name: "UTF three-byte CRLF",
			from: "ግዕዝ\r\nລາວ\r\n中文\r\n",
			want: "ግዕዝ\r\nລາວ\r\n中文\r\n"},
		{name: "UTF three-byte CR",
			from: "ግዕዝ\rລາວ\r中文\r",
			want: "ግዕዝ\r\nລາວ\r\n中文\r\n"},
		{name: "UTF three-byte CRLF",
			from: "ግዕዝ\nລາວ\n中文\n",
			want: "ግዕዝ\r\nລາວ\r\n中文\r\n"},
		{name: "UTF four-byte CRLF",
			from: "📻\r\n🏴‍☠️\r\n",
			want: "📻\r\n🏴‍☠️\r\n"},
		{name: "UTF four-byte CR",
			from: "📻\r🏴‍☠️\r",
			want: "📻\r\n🏴‍☠️\r\n"},
		{name: "UTF four-byte LF",
			from: "📻\n🏴‍☠️\n",
			want: "📻\r\n🏴‍☠️\r\n"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ensureCRLF(tc.from); got != tc.want {
				t.Errorf("ensureCRLF(%q) = %q, want %q", tc.from, got, tc.want)
			}
		})
	}
}

func TestIsAllDigits(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{s: "", want: true},
		{s: "0", want: true},
		{s: "9", want: true},
		{s: "123", want: true},
		{s: "9876543210", want: true},
		{s: "/123", want: false},
		{s: "42:", want: false},
		{s: "54 46", want: false},
		{s: "5N5", want: false},
		{s: "X", want: false},
		{s: "௭௩", want: false},
		{s: "⑦", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.s, func(t *testing.T) {
			if got := isAllDigits(tc.s); got != tc.want {
				t.Errorf("isAllDigits(%q) = %v, want %v", tc.s, got, tc.want)
			}
		})
	}
}

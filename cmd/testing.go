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

	"github.com/flwyd/adif-multitool/adif"
)

func readers(rs ...adif.ReadWriter) map[adif.Format]adif.Reader {
	res := make(map[adif.Format]adif.Reader)
	for _, r := range rs {
		f, err := adif.ParseFormat(r.String())
		if err != nil {
			panic(fmt.Sprintf("Unknown adif.Format %q: %v", r, err))
		}
		res[f] = r.(adif.Reader)
	}
	return res
}

func writers(rs ...adif.ReadWriter) map[adif.Format]adif.Writer {
	res := make(map[adif.Format]adif.Writer)
	for _, r := range rs {
		f, err := adif.ParseFormat(r.String())
		if err != nil {
			panic(fmt.Sprintf("Unknown adif.Format %q: %v", r, err))
		}
		res[f] = r.(adif.Writer)
	}
	return res
}

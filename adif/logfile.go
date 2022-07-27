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

type Logfile struct {
	Records  []*Record
	Header   *Record
	Comment  string
	Filename string
}

func NewLogfile(name string) *Logfile {
	return &Logfile{Records: make([]*Record, 0), Header: NewRecord(), Filename: name}
}

func (f *Logfile) String() string {
	if len(f.Filename) == 0 {
		return "(standard input)"
	}
	return f.Filename
}

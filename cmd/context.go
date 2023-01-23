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
	"io"

	"github.com/flwyd/adif-multitool/adif"
)

type Context struct {
	ProgramName    string
	ProgramVersion string
	ADIFVersion    string
	InputFormat    adif.Format
	OutputFormat   adif.Format
	Readers        map[adif.Format]adif.Reader
	Writers        map[adif.Format]adif.Writer
	Out            io.Writer
	CommandCtx     any
	fs             filesystem
}

func (c *Context) SetHeaders(l *adif.Logfile) {
	l.Header.Set(adif.Field{Name: "ADIF_VER", Value: c.ADIFVersion})
	l.Header.Set(adif.Field{Name: "PROGRAMID", Value: c.ProgramName})
	l.Header.Set(adif.Field{Name: "PROGRAMVERSION", Value: c.ProgramVersion})
}

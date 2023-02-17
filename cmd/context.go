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
	"github.com/flwyd/adif-multitool/adif/spec"
)

type Context struct {
	InputFormat   adif.Format
	OutputFormat  adif.Format
	Readers       map[adif.Format]adif.Reader
	Writers       map[adif.Format]adif.Writer
	Out           io.Writer
	CommandCtx    any
	UserdefFields UserdefFieldList
	Prepare       func(*adif.Logfile)
	fs            filesystem
}

func testPrepare(comment, adifVer, progName, progVer string) func(l *adif.Logfile) {
	return func(l *adif.Logfile) {
		l.Header.SetComment(comment)
		l.Header.Set(adif.Field{Name: spec.AdifVerField.Name, Value: adifVer})
		l.Header.Set(adif.Field{Name: spec.ProgramidField.Name, Value: progName})
		l.Header.Set(adif.Field{Name: spec.ProgramversionField.Name, Value: progVer})
	}
}

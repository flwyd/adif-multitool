// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/flwyd/adif-multitool/adif"
	"github.com/flwyd/adif-multitool/adif/spec"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"adifmt": func() int { return runMain(testPrepare) },
	}))
}

func TestScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata",
	})
}

func testPrepare(l *adif.Logfile) {
	l.Header.SetComment(fmt.Sprintf("Generated with %d records by %s", len(l.Records), helpUrl))
	l.Header.Set(adif.Field{Name: spec.AdifVerField.Name, Value: spec.ADIFVersion})
	l.Header.Set(adif.Field{Name: spec.CreatedTimestampField.Name, Value: "23450607 080910"})
	l.Header.Set(adif.Field{Name: spec.ProgramidField.Name, Value: programName})
	l.Header.Set(adif.Field{Name: spec.ProgramversionField.Name, Value: version})
}

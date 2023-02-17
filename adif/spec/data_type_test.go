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

package spec

import (
	"testing"

	"github.com/flwyd/adif-multitool/adif"
)

func TestDataTypeDeclared(t *testing.T) {
	if len(DataTypes) == 0 {
		t.Errorf("len(DataTypes) == 0, make sure to run go generate aidf/spec")
	}
}

func TestDataTypeNames(t *testing.T) {
	for name, d := range DataTypes {
		if len(name) != 1 && name != d.Name {
			t.Errorf("DataTypes[%q].Name got %q, want %q", name, d.Name, name)
		}
	}
}

func TestDataTypeIndicators(t *testing.T) {
	for _, id := range []string{"S", "I", "M", "G", "E", "B", "N", "D", "T", "L"} {
		want, err := adif.DataTypeFromIndicator(id)
		if err != nil {
			t.Fatalf("%q is not an adif.DataType indicator", id)
		}
		if got, ok := DataTypes[id]; !ok {
			t.Errorf("DataTypes[%q] missing", id)
		} else if got.Indicator != want.Indicator() {
			t.Errorf("DataTypes[%q].Indicator got %q, want %q", id, got.Indicator, want.Indicator())
		}
	}
}

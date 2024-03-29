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

// Package spec contains data types, field definitions, and enumerations defined
// by the ADIF specification from https://adif.org.uk/
// Most structures in this package are automatically generated.
package spec

//go:generate go run ./mkspec https://adif.org.uk/314/ADIF_314_released_exports_2022_12_06.zip

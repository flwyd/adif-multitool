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
	"regexp"
	"strings"
	"time"

	"github.com/flwyd/adif-multitool/adif"
)

type fieldList []string

func (f *fieldList) String() string {
	return strings.Join(*f, ",")
}

func (f *fieldList) Get() fieldList { return *f }

func (f *fieldList) Set(s string) error {
	fields := strings.Split(s, ",")
	for _, x := range fields {
		x = strings.ToUpper(strings.TrimSpace(x))
		if x == "" {
			return fmt.Errorf("empty field name in %q", fields)
		}
		*f = append(*f, x)
	}
	return nil
}

type fieldAssignments struct {
	values   []adif.Field
	validate func(k, v string) error
}

// TODO allow name:type fields
var adifNamePat = regexp.MustCompile(`^\w+$`)

func validateAlphanumName(name, value string) error {
	if !adifNamePat.MatchString(name) {
		return fmt.Errorf("invalid ADIF field name %q", name)
	}
	return nil
}

func newFieldAssignments(validate func(k, v string) error) fieldAssignments {
	return fieldAssignments{values: make([]adif.Field, 0), validate: validate}
}

func (a *fieldAssignments) add(k, v string) error {
	k = strings.ToUpper(strings.TrimSpace(k))
	if err := a.validate(k, v); err != nil {
		return err
	}
	a.values = append(a.values, adif.Field{Name: k, Value: v})
	return nil
}

func (a *fieldAssignments) String() string {
	res := make([]string, len(a.values))
	for i, f := range a.values {
		res[i] = f.String()
	}
	return strings.Join(res, ";;")
}

func (a *fieldAssignments) Get() fieldAssignments { return *a }

func (f *fieldAssignments) Set(s string) error {
	chunks := strings.Split(s, ";;")
	a := newFieldAssignments(validateAlphanumName)
	for _, c := range chunks {
		c = strings.TrimSpace(c)
		kv := strings.SplitN(c, "=", 2)
		if len(kv) == 1 || len(kv[0]) == 0 {
			return fmt.Errorf(`expected "name=value", got %q`, c)
		}
		if err := a.add(kv[0], kv[1]); err != nil {
			return fmt.Errorf("validation error on %q: %v", c, err)
		}
	}
	f.values = append(f.values, a.values...)
	return nil
}

type timeZone struct {
	tz *time.Location
}

func (z *timeZone) String() string {
	if z.tz == nil {
		return time.UTC.String()
	}
	return z.tz.String()
}

func (z *timeZone) Get() *time.Location {
	if z.tz == nil {
		return time.UTC
	}
	return z.tz
}

func (z *timeZone) Set(s string) error {
	l, err := time.LoadLocation(s)
	if err != nil {
		return fmt.Errorf("invalid time zone %q: %v", s, err)
	}
	z.tz = l
	return nil
}

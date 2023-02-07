// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"time"

	"github.com/flwyd/adif-multitool/adif"
	"github.com/flwyd/adif-multitool/adif/spec"
)

var Edit = Command{Name: "edit", Run: runEdit,
	Description: "Add, change, remove, or adjust field values"}

type EditContext struct {
	Add         FieldAssignments
	Set         FieldAssignments
	Remove      FieldList
	RemoveBlank bool
	FromZone    TimeZone
	ToZone      TimeZone
}

func runEdit(ctx *Context, args []string) error {
	cctx := ctx.CommandCtx.(*EditContext)
	remove := make(map[string]bool)
	for _, n := range cctx.Remove {
		remove[n] = true
	}
	set := make(map[string]adif.Field)
	for _, f := range cctx.Set.values {
		if remove[f.Name] {
			return fmt.Errorf("%q in both -set and -remove", f.Name)
		}
		set[f.Name] = f
	}
	for _, f := range cctx.Add.values {
		if remove[f.Name] {
			return fmt.Errorf("%q in both -add and -remove, use -set to change values", f.Name)
		}
		if _, ok := set[f.Name]; ok {
			return fmt.Errorf("%q in both -set and -add", f.Name)
		}
	}
	fromTz := cctx.FromZone.Get()
	toTz := cctx.ToZone.Get()
	adjustTz := fromTz.String() != toTz.String()
	out := adif.NewLogfile()
	var comments commentCatcher
	for _, f := range filesOrStdin(args) {
		l, err := readFile(ctx, f)
		if err != nil {
			return err
		}
		updateFieldOrder(out, l.FieldOrder)
		// TODO merge headers
		for _, r := range l.Records {
			seen := make(map[string]bool)
			old := r.Fields()
			fields := make([]adif.Field, 0, len(old))
			for _, f := range old {
				if remove[f.Name] {
					continue
				}
				if cctx.RemoveBlank && f.Value == "" {
					continue
				}
				seen[f.Name] = true
				if v, ok := set[f.Name]; ok {
					f = v
				}
				fields = append(fields, f)
			}
			for _, f := range cctx.Set.values {
				if !seen[f.Name] {
					fields = append(fields, f)
				}
				seen[f.Name] = true
			}
			for _, f := range cctx.Add.values {
				if !seen[f.Name] {
					fields = append(fields, f)
				}
				seen[f.Name] = true
			}
			if len(fields) > 0 {
				rec := adif.NewRecord(fields...)
				if adjustTz {
					if err := adjustTimeZone(rec, fromTz, toTz); err != nil {
						return fmt.Errorf("could not adjust time zone: %w", err)
					}
				}
				out.Records = append(out.Records, rec)
			}
		}
		comments.read(l, f)
	}
	comments.write(out)
	return write(ctx, out)
}

func adjustTimeZone(r *adif.Record, from, to *time.Location) error {
	dayfmt := "20060102"
	adjust := func(timef, dayf, dayfallback adif.Field) error {
		timefmt := "150405"
		if len(timef.Value) == 4 {
			timefmt = "1504"
		}
		day := dayf
		if dayf.Value == "" {
			day = dayfallback
		}
		t, err := time.ParseInLocation(dayfmt+timefmt, day.Value+timef.Value, from)
		if err != nil {
			return fmt.Errorf("invalid %s %s: %w", day, timef, err)
		}
		t = t.In(to)
		timef.Value = t.Format(timefmt)
		r.Set(timef)
		if dayf.Value != "" {
			dayf.Value = t.Format(dayfmt)
			r.Set(dayf)
		}
		return nil
	}
	ton, tonok := r.Get(spec.TimeOnField.Name)
	toff, toffok := r.Get(spec.TimeOffField.Name)
	don, donok := r.Get(spec.QsoDateField.Name)
	doff, doffok := r.Get(spec.QsoDateOffField.Name)
	if !donok && !doffok {
		return fmt.Errorf("no %s or %s field for %s %s",
			spec.QsoDateField.Name, spec.QsoDateOffField.Name, ton, toff)
	}
	if tonok {
		if err := adjust(ton, don, doff); err != nil {
			return err
		}
	}
	if toffok {
		if err := adjust(toff, doff, don); err != nil {
			return err
		}
	}
	return nil
}

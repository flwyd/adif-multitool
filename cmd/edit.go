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
	"strings"
	"time"

	"github.com/flwyd/adif-multitool/adif"
	"github.com/flwyd/adif-multitool/adif/spec"
)

var Edit = Command{Name: "edit", Run: runEdit, Help: helpEdit,
	Description: "Add, change, remove, or adjust field values"}

type EditContext struct {
	Add         FieldAssignments
	Set         FieldAssignments
	Rename      FieldAssignments
	Remove      FieldList
	RemoveBlank bool
	Cond        ConditionValue
	FromZone    TimeZone
	ToZone      TimeZone
}

func helpEdit() string {
	return fmt.Sprintf("Time zone adjustments affect %s, %s, %s, and %s.\n",
		spec.TimeOnField.Name, spec.TimeOffField.Name,
		spec.QsoDateField.Name, spec.QsoDateOffField.Name) +
		"Renames can be circular, e.g. --rename my_lat=my_lon --rename my_lon=my_lat\n\n" +
		helpCondition("edit")
}

func runEdit(ctx *Context, args []string) error {
	cctx := ctx.CommandCtx.(*EditContext)
	remove := make(map[string]bool)
	for _, n := range cctx.Remove {
		remove[n] = true
	}
	set := make(map[string]adif.Field)
	rename := make(map[string]string)
	renameFrom := make(map[string]string)
	for _, f := range cctx.Set.values {
		if remove[f.Name] {
			return fmt.Errorf("%q in both --set and --remove", f.Name)
		}
		set[f.Name] = f
	}
	for _, f := range cctx.Add.values {
		if remove[f.Name] {
			return fmt.Errorf("%q in both --add and --remove, use --set to change values", f.Name)
		}
		if _, ok := set[f.Name]; ok {
			return fmt.Errorf("%q in both --set and --add", f.Name)
		}
	}
	for _, f := range cctx.Rename.values {
		v := strings.ToUpper(f.Value)
		if e := renameFrom[v]; e != "" {
			return fmt.Errorf("duplicate rename target --rename %s=%s and --rename %s=%s", f.Name, v, e, v)
		}
		rename[f.Name] = v
		renameFrom[v] = f.Name
		if remove[f.Name] {
			return fmt.Errorf("%q in both --rename and --remove, rename will leave field unset", f.Name)
		}
		if _, ok := set[v]; ok {
			return fmt.Errorf("%q in --set and --rename %s=%s, set would override rename", v, f.Name, f.Value)
		}
	}
	fromTz := cctx.FromZone.Get()
	toTz := cctx.ToZone.Get()
	adjustTz := fromTz.String() != toTz.String()
	cond := cctx.Cond.Get()
	acc, err := newAccumulator(ctx)
	if err != nil {
		return err
	}
	for _, f := range filesOrStdin(args) {
		l, err := acc.read(f)
		if err != nil {
			return err
		}
		updateFieldOrder(acc.Out, l.FieldOrder)
		for _, r := range l.Records {
			eval := recordEvalContext{record: r, lang: ctx.Locale}
			if !cond.Evaluate(eval) {
				acc.Out.AddRecord(r) // edit condition doesn't match, pass through
				continue
			}
			seen := make(map[string]string)
			old := r.Fields()
			fields := make([]adif.Field, 0, len(old))
			for _, f := range old {
				if remove[f.Name] {
					continue
				}
				if cctx.RemoveBlank && f.Value == "" {
					continue
				}
				if v, ok := set[f.Name]; ok {
					f = v
				}
				if dest := rename[f.Name]; dest != "" {
					if s := seen[dest]; s != "" {
						if f.Value != "" {
							return fmt.Errorf("X rename %s to %s would overwrite value %q with %q; to overwrite all use --remove %s --rename %s=%s", f.Name, dest, s, f.Value, dest, f.Name, dest)
						} else {
							continue
						}
					}
					f.Name = dest
				} else if src := renameFrom[f.Name]; src != "" {
					if s := seen[f.Name]; s != "" {
						if f.Value != "" {
							return fmt.Errorf("Y rename %s to %s would overwrite value %q with %q; to overwrite all use --remove %s --rename %s=%s", src, f.Name, f.Value, s, f.Name, src, f.Name)
						} else {
							continue
						}
					}
				}
				seen[f.Name] = f.Value
				fields = append(fields, f)
			}
			for _, f := range cctx.Set.values {
				if _, ok := seen[f.Name]; !ok {
					fields = append(fields, f)
				}
				seen[f.Name] = f.Value
			}
			for _, f := range cctx.Add.values {
				if v, ok := seen[f.Name]; !ok || v == "" {
					fields = append(fields, f)
					seen[f.Name] = f.Value
				}
			}
			if len(fields) > 0 {
				rec := adif.NewRecord(fields...)
				if adjustTz {
					if err := adjustTimeZone(rec, fromTz, toTz); err != nil {
						return fmt.Errorf("could not adjust time zone: %w", err)
					}
				}
				acc.Out.AddRecord(rec)
			}
		}
	}
	if err := acc.prepare(); err != nil {
		return err
	}
	return write(ctx, acc.Out)
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

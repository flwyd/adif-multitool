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
	"flag"
	"fmt"
	"regexp"
	"strings"

	"github.com/flwyd/adif-multitool/adif"
	"github.com/flwyd/adif-multitool/adif/spec"
	"golang.org/x/text/language"
)

type EvaluationContext interface {
	Compare(a, b adif.Field) (int, error)
	Get(name string) adif.Field
	Cast(name, value string) adif.Field
}

type recordEvalContext struct {
	record *adif.Record
	lang   language.Tag
}

func (r recordEvalContext) Get(name string) adif.Field {
	rec, _ := r.record.Get(name)
	return rec
}

func (r recordEvalContext) Cast(name, value string) adif.Field {
	f, ok := r.record.Get(name)
	f.Value = value
	if !ok {
		if sf, ok := spec.Fields[strings.ToUpper(name)]; ok {
			f.Type, _ = adif.DataTypeFromIndicator(sf.Type.Indicator)
		}
	}
	return f
}

func (r recordEvalContext) Compare(a, b adif.Field) (int, error) {
	var f spec.Field
	var ok bool
	if f, ok = spec.Fields[strings.ToUpper(a.Name)]; !ok {
		if f, ok = spec.Fields[strings.ToUpper(b.Name)]; !ok {
			f = spec.Field{Name: strings.ToUpper(a.Name), Type: spec.DataTypes[a.Type.String()]}
		}
	}
	comp := spec.ComparatorForField(f, r.lang)
	return comp(a.Value, b.Value)
}

type Condition interface {
	Evaluate(EvaluationContext) bool // maybe (bool, error)?
	String() string
}

type operator string // TODO enum
const (
	OpEqual            operator = "="
	OpLessThan                  = "<"
	OpLessThanEqual             = "<="
	OpGreaterThan               = ">"
	OpGreaterThanEqual          = ">="
	// TODO OpGlob
)

type comparison struct {
	Op        operator
	FieldName string
	Operands  []string
	Negate    bool
}

func (c comparison) String() string {
	not := ""
	if c.Negate {
		not = "NOT "
	}
	return fmt.Sprintf("%s%s%s%s", not, c.FieldName, c.Op, strings.Join(c.Operands, "|"))
}

func (c comparison) Evaluate(e EvaluationContext) bool {
	result := func(b bool) bool {
		if c.Negate {
			return !b
		}
		return b
	}
	f := e.Get(c.FieldName)
	for _, o := range c.Operands {
		var v adif.Field
		if isSurrounded(o, "{", "}") {
			v = e.Get(trimSurrounding(o, "{", "}"))
		} else {
			v = e.Cast(c.FieldName, o)
		}
		comp, err := e.Compare(f, v)
		if err != nil {
			return false
		}
		switch c.Op {
		case OpEqual:
			if comp == 0 {
				return result(true)
			}
		case OpLessThan:
			if comp < 0 {
				return result(true)
			}
		case OpLessThanEqual:
			if comp <= 0 {
				return result(true)
			}
		case OpGreaterThan:
			if comp > 0 {
				return result(true)
			}
		case OpGreaterThanEqual:
			if comp >= 0 {
				return result(true)
			}
		default:
			panic("Unknown operator " + c.Op)
		}
	}
	return result(false)
}

type junction struct {
	Terms []Condition
	Any   bool // if true, or logic, otherwise all (and logic)
}

func (j junction) String() string {
	op := " AND "
	if j.Any {
		op = " OR "
	}
	s := make([]string, len(j.Terms))
	for i, t := range j.Terms {
		s[i] = t.String()
	}
	return strings.Join(s, op)
}

func (j junction) Evaluate(e EvaluationContext) bool {
	if len(j.Terms) == 0 {
		return true
	}
	for _, t := range j.Terms {
		v := t.Evaluate(e)
		if v && j.Any {
			return true
		}
		if !v && !j.Any {
			return false
		}
	}
	return !j.Any
}

type ConditionValue struct {
	cur  junction
	done []junction
}

func (cv *ConditionValue) IfFlag() flag.Value {
	return &ifValue{cv: cv, negate: false}
}

func (cv *ConditionValue) IfNotFlag() flag.Value {
	return &ifValue{cv: cv, negate: true}
}

func (cv *ConditionValue) OrIfFlag() flag.Value {
	return &orIfValue{cv: cv, negate: false}
}

func (cv *ConditionValue) OrIfNotFlag() flag.Value {
	return &orIfValue{cv: cv, negate: true}
}

func (cv *ConditionValue) Get() junction {
	j := junction{Any: true}
	if len(cv.cur.Terms) > 0 {
		j.Terms = []Condition{cv.cur}
	}
	for _, t := range cv.done {
		j.Terms = append(j.Terms, t)
	}
	return j
}

var conditionalPat = regexp.MustCompile(`^(\w+)(=|[<>]=?)(.*)`)

type ifValue struct {
	cv     *ConditionValue
	negate bool
}

func (v *ifValue) String() string { return "" }

func (i *ifValue) Set(s string) error {
	g := conditionalPat.FindStringSubmatch(s)
	if g == nil {
		return fmt.Errorf("invalid if condition: %q", s)
	}
	opts := strings.Split(g[3], "|")
	if g[3] == "" {
		if g[2] != string(OpEqual) && g[2] != string(OpGreaterThan) {
			return fmt.Errorf("cannot use %s with empty string: %q", g[2], s)
		}
		opts = []string{""}
	}
	c := comparison{Op: operator(g[2]), FieldName: g[1], Operands: opts, Negate: i.negate}
	i.cv.cur.Terms = append(i.cv.cur.Terms, c)
	return nil
}

type orIfValue struct {
	cv     *ConditionValue
	negate bool
}

func (i *orIfValue) String() string { return "" }

func (i *orIfValue) Set(s string) error {
	if len(i.cv.cur.Terms) > 0 {
		i.cv.done = append(i.cv.done, i.cv.cur)
	}
	i.cv.cur = junction{}
	return (&ifValue{cv: i.cv, negate: i.negate}).Set(s)
}

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
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

type FieldComparator func(a, b string) (int, error)

func ComparatorForField(f Field, locale language.Tag) FieldComparator {
	var c FieldComparator
	switch f.Type.Name {
	case StringDataType.Name, MultilineStringDataType.Name, CharacterDataType.Name,
		StringDataType.Indicator, MultilineStringDataType.Indicator:
		c = compareStringsBasic
	case IntlStringDataType.Name, IntlMultilineStringDataType.Name, IntlCharacterDataType.Name,
		IntlStringDataType.Indicator, IntlMultilineStringDataType.Indicator:
		c = compareStringsLocale(locale, collate.OptionsFromTag(locale), collate.IgnoreCase, collate.IgnoreWidth)
	case NumberDataType.Name, IntegerDataType.Name, PositiveIntegerDataType.Name, DigitDataType.Name,
		NumberDataType.Indicator, IntegerDataType.Indicator:
		c = compareNumbers
	case DateDataType.Name, DateDataType.Indicator:
		c = compareDates
	case TimeDataType.Name, TimeDataType.Indicator:
		c = compareTimes
	case BooleanDataType.Name, BooleanDataType.Indicator:
		c = compareBooleans
	case GridsquareField.Name, GridsquareExtField.Name:
		c = compareStringsBasic
	case LocationDataType.Name, LocationDataType.Indicator:
		c = compareLocations
	case EnumerationDataType.Name, EnumerationDataType.Indicator:
		switch f.EnumName {
		case BandEnumeration.Name:
			c = compareBands
		case DxccEntityCodeEnumeration.Name:
			c = compareNumbers
		default:
			c = compareStringsBasic
		}
	case IOTARefNoDataType.Name, POTARefDataType.Name, SOTARefDataType.Name, WWFFRefDataType.Name:
		c = compareStringsBasic
	case AwardListDataType.Name, CreditListDataType.Name, GridSquareListDataType.Name, POTARefListDataType.Name, SponsoredAwardListDataType.Name:
		c = compareStringLists(",")
	case SecondarySubdivisionListDataType.Name:
		c = compareStringLists(":")
	default:
		c = compareStringsBasic
	}
	return compareEmptyFirst(c)
}

func compareEmptyFirst(c FieldComparator) FieldComparator {
	return func(a, b string) (int, error) {
		if a == "" && b == "" {
			return 0, nil
		}
		if a == "" {
			return -1, nil
		}
		if b == "" {
			return 1, nil
		}
		return c(a, b)
	}
}

func compareStringsLocale(l language.Tag, opts ...collate.Option) FieldComparator {
	col := collate.New(l, opts...)
	return func(a, b string) (int, error) { return col.CompareString(a, b), nil }
}

var compareStringsBasic = compareStringsLocale(language.Und, collate.Loose)

func compareNumbers(a, b string) (int, error) {
	an, err := strconv.ParseFloat(a, 10)
	if err != nil {
		return 0, err
	}
	bn, err := strconv.ParseFloat(b, 10)
	if err != nil {
		return 0, err
	}
	if an == bn {
		return 0, nil
	}
	if an < bn {
		return -1, nil
	}
	return 1, nil
}

func compareDates(a, b string) (int, error) {
	at, err := time.Parse("20060102", a)
	if err != nil {
		return 0, err
	}
	bt, err := time.Parse("20060102", b)
	if err != nil {
		return 0, err
	}
	if at.Equal(bt) {
		return 0, nil
	}
	if at.Before(bt) {
		return -1, nil
	}
	return 1, nil
}

func compareTimes(a, b string) (int, error) {
	var at, bt time.Time
	var err error
	switch len(a) {
	case 4:
		if at, err = time.Parse("1504", a); err != nil {
			return 0, err
		}
	case 6:
		if at, err = time.Parse("150405", a); err != nil {
			return 0, err
		}
	default:
		return 0, fmt.Errorf("invalid time format %q", a)
	}
	switch len(b) {
	case 4:
		if bt, err = time.Parse("1504", b); err != nil {
			return 0, err
		}
	case 6:
		if bt, err = time.Parse("150405", b); err != nil {
			return 0, err
		}
	default:
		return 0, fmt.Errorf("invalid time format %q", b)
	}
	if at.Equal(bt) {
		return 0, nil
	}
	if at.Before(bt) {
		return -1, nil
	}
	return 1, nil
}

func compareBooleans(a, b string) (int, error) {
	ab := a == "Y" || a == "y"
	if !ab && a != "N" && a != "n" {
		return 0, fmt.Errorf("invalid boolean value %q", a)
	}
	bb := b == "Y" || b == "y"
	if !bb && b != "N" && b != "n" {
		return 0, fmt.Errorf("invalid boolean value %q", b)
	}
	if ab == bb {
		return 0, nil
	}
	if !ab {
		return -1, nil
	}
	return 1, nil
}

func parseLocation(s string) (dir rune, degrees int64, minutes float64, err error) {
	g := locationPat.FindStringSubmatch(s)
	if g == nil {
		err = fmt.Errorf("invalid location format %q", s)
		return
	}
	switch g[1] {
	case "W", "w":
		dir = 'W'
	case "E", "e":
		dir = 'E'
	case "S", "s":
		dir = 'S'
	case "N", "n":
		dir = 'N'
	default:
		err = fmt.Errorf("invalid location direction %q", s)
		return
	}
	degrees, err = strconv.ParseInt(g[2], 10, 64)
	if err != nil {
		return
	}
	minutes, err = strconv.ParseFloat(g[3], 10)
	return
}

var directionOrder = map[[2]rune]int{
	{'W', 'W'}: 0, {'E', 'E'}: 0, {'S', 'S'}: 0, {'N', 'N'}: 0,
	// west before all others
	{'W', 'E'}: -1, {'W', 'S'}: -1, {'W', 'N'}: -1,
	// east before north/south
	{'E', 'W'}: 1, {'E', 'S'}: -1, {'E', 'N'}: -1,
	// south before north
	{'S', 'W'}: 1, {'S', 'E'}: 1, {'S', 'N'}: -1,
	{'N', 'W'}: 1, {'N', 'E'}: 1, {'N', 'S'}: 1,
}

func compareLocations(a, b string) (int, error) {
	adir, adeg, amin, err := parseLocation(a)
	if err != nil {
		return 0, err
	}
	bdir, bdeg, bmin, err := parseLocation(b)
	if err != nil {
		return 0, err
	}
	// Compare west-to-east then south-to-north so sorting by gridsquare is
	// the same as sorting by lon/lat.
	if dircomp := directionOrder[[2]rune{adir, bdir}]; dircomp != 0 {
		if adeg == 0 && bdeg == 0 && amin == 0.0 && bmin == 0.0 &&
			(adir == 'W' && bdir == 'E' || adir == 'E' && bdir == 'W' || adir == 'S' && bdir == 'N' || adir == 'N' && bdir == 'S') {
			return 0, nil // equator or prime meridian
		}
		if adeg == 180 && bdeg == 180 && amin == 0.0 && bmin == 0.0 &&
			(adir == 'W' && bdir == 'E' || adir == 'E' && bdir == 'W') {
			return 0, nil // antiprime meridian
		}
		return dircomp, nil
	}
	// degrees and minuts are always positive, but larger values of west/south
	// comee before smaller values whiile smaller east/north values come first
	adec := float64(adeg) + (amin / 60.0)
	bdec := float64(bdeg) + (bmin / 60.0)
	switch adir {
	case 'W', 'S':
		if adec == bdec {
			return 0, nil
		}
		if adec > bdec {
			return -1, nil // a is to the west/south of b
		}
		return 1, nil // a is to the east/north of b
	case 'E', 'N':
		if adec == bdec {
			return 0, nil
		}
		if adec < bdec {
			return -1, nil // a is to the west/south of b
		}
		return 1, nil // a is to the east/north of b
	default:
		panic(fmt.Sprintf("unknown direction %q in %v", adir, a))
	}
}

func compareBands(a, b string) (int, error) {
	ab := BandEnumeration.Value(a)
	if len(ab) != 1 {
		return 0, fmt.Errorf("unknown band %q", a)
	}
	bb := BandEnumeration.Value(b)
	if len(bb) != 1 {
		return 0, fmt.Errorf("unknown band %q", b)
	}
	aband := ab[0].(BandEnum)
	bband := bb[0].(BandEnum)
	if aband.Band == bband.Band {
		return 0, nil
	}
	alow, err := strconv.ParseFloat(aband.LowerFreqMhz, 10)
	if err != nil {
		return 0, err
	}
	blow, err := strconv.ParseFloat(bband.LowerFreqMhz, 10)
	if err != nil {
		return 0, err
	}
	if alow < blow {
		return -1, nil
	}
	return 1, nil
}

func compareStringLists(sep string) FieldComparator {
	return func(a, b string) (int, error) {
		alist := strings.Split(strings.ToUpper(a), sep)
		blist := strings.Split(strings.ToUpper(b), sep)
		sort.Strings(alist) // treat foo,bar as equal to bar,foo
		sort.Strings(blist)
		min := len(alist)
		if len(blist) < len(alist) {
			min = len(blist)
		}
		for i := 0; i < min; i++ {
			if c, err := compareStringsBasic(alist[i], blist[i]); c != 0 || err != nil {
				return c, err
			}
		}
		return len(alist) - len(blist), nil
	}
}

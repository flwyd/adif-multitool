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
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/flwyd/adif-multitool/adif"
	"github.com/flwyd/adif-multitool/adif/spec"
)

var Infer = Command{Name: "infer", Run: runInfer, Help: helpInfer,
	Description: "Add missing fields based on present fields"}

type InferContext struct {
	Fields     FieldList
	CommentLog bool
}

type inferrer func(*adif.Record, string) bool

var inferrers = map[string]inferrer{
	spec.BandField.Name:            inferBand,
	spec.BandRxField.Name:          inferBand,
	spec.ModeField.Name:            inferMode,
	spec.CountryField.Name:         inferCountry,
	spec.MyCountryField.Name:       inferCountry,
	spec.DxccField.Name:            inferDxcc,
	spec.MyDxccField.Name:          inferDxcc,
	spec.GridsquareField.Name:      inferGridsquare,
	spec.GridsquareExtField.Name:   inferGridsquare,
	spec.MyGridsquareField.Name:    inferGridsquare,
	spec.MyGridsquareExtField.Name: inferGridsquare,
	spec.LatField.Name:             inferLatLon,
	spec.MyLatField.Name:           inferLatLon,
	spec.LonField.Name:             inferLatLon,
	spec.MyLonField.Name:           inferLatLon,
	spec.OperatorField.Name:        inferStation,
	spec.StationCallsignField.Name: inferStation,
	spec.OwnerCallsignField.Name:   inferStation,
	spec.UsacaCountiesField.Name:   inferUSCounty,
	spec.MyUsacaCountiesField.Name: inferUSCounty,
	spec.CntyField.Name:            inferUSCounty, // US is the only secondary subdivision
	spec.MyCntyField.Name:          inferUSCounty, // with a special field
	spec.SigInfoField.Name:         inferSigInfo,
	spec.MySigInfoField.Name:       inferSigInfo,
	spec.IotaField.Name:            inferProgramRef("IOTA"),
	spec.MyIotaField.Name:          inferProgramRef("IOTA"),
	spec.PotaRefField.Name:         inferProgramRef("POTA"),
	spec.MyPotaRefField.Name:       inferProgramRef("POTA"),
	spec.SotaRefField.Name:         inferProgramRef("SOTA"),
	spec.MySotaRefField.Name:       inferProgramRef("SOTA"),
	spec.WwffRefField.Name:         inferProgramRef("WWFF"),
	spec.MyWwffRefField.Name:       inferProgramRef("WWFF"),
}

func helpInfer() string {
	res := &strings.Builder{}
	res.WriteString("Inferable fields:\n")
	fromfmt := "  %s from %s\n"
	fmt.Fprintf(res, fromfmt, spec.BandField.Name, spec.FreqField.Name)
	fmt.Fprintf(res, fromfmt, spec.BandRxField.Name, spec.FreqRxField.Name)
	fmt.Fprintf(res, fromfmt, spec.ModeField.Name, spec.SubmodeField.Name)
	fmt.Fprintf(res, fromfmt, spec.CountryField.Name, spec.DxccField.Name)
	fmt.Fprintf(res, fromfmt, spec.MyCountryField.Name, spec.MyDxccField.Name)
	fmt.Fprintf(res, fromfmt, spec.DxccField.Name, spec.CountryField.Name)
	fmt.Fprintf(res, fromfmt, spec.MyDxccField.Name, spec.MyCountryField.Name)
	fmt.Fprintf(res, fromfmt, spec.CntyField.Name, spec.UsacaCountiesField.Name)
	fmt.Fprintf(res, fromfmt, spec.MyCntyField.Name, spec.MyUsacaCountiesField.Name)
	fmt.Fprintf(res, fromfmt, spec.UsacaCountiesField.Name, spec.CntyField.Name)
	fmt.Fprintf(res, fromfmt, spec.MyUsacaCountiesField.Name, spec.MyCntyField.Name)

	gsfmt := "  %s and %s from %s/%s\n"
	fmt.Fprintf(res, gsfmt, spec.GridsquareField.Name, spec.GridsquareExtField.Name, spec.LatField.Name, spec.LonField.Name)
	fmt.Fprintf(res, gsfmt, spec.MyGridsquareField.Name, spec.MyGridsquareExtField.Name, spec.MyLatField.Name, spec.MyLonField.Name)
	llfmt := "  %s/%s from %s and optionally %s\n"
	fmt.Fprintf(res, llfmt, spec.LatField.Name, spec.LonField.Name, spec.GridsquareField.Name, spec.GridsquareExtField.Name)
	fmt.Fprintf(res, llfmt, spec.MyLatField.Name, spec.MyLonField.Name, spec.MyGridsquareField.Name, spec.MyGridsquareExtField.Name)

	fmt.Fprintf(res, fromfmt, spec.OperatorField.Name, spec.GuestOpField.Name)
	fmt.Fprintf(res, "  %s from %s or %s\n", spec.StationCallsignField.Name, spec.OperatorField.Name, spec.GuestOpField.Name)
	fmt.Fprintf(res, "  %s from %s, %s, or %s\n", spec.OwnerCallsignField.Name, spec.StationCallsignField.Name, spec.OperatorField.Name, spec.GuestOpField.Name)

	sifmt := "  %s from one of %s, %s, %s, %s based on %s\n"
	sigfmt := "    (sets %s if unset and only one of the others is set)\n"
	fmt.Fprintf(res, sifmt, spec.SigInfoField.Name, spec.IotaField.Name, spec.PotaRefField.Name, spec.SotaRefField.Name, spec.WwffRefField.Name, spec.SigField.Name)
	fmt.Fprintf(res, sigfmt, spec.SigField.Name)
	fmt.Fprintf(res, sifmt, spec.MySigInfoField.Name, spec.MyIotaField.Name, spec.MyPotaRefField.Name, spec.MySotaRefField.Name, spec.MyWwffRefField.Name, spec.MySigField.Name)
	fmt.Fprintf(res, sigfmt, spec.MySigField.Name)
	progfmt := "  %s from %s if %s is %q\n"
	progs := []struct {
		field spec.Field
		prog  string
	}{
		{field: spec.IotaField, prog: "IOTA"},
		{field: spec.PotaRefField, prog: "POTA"},
		{field: spec.SotaRefField, prog: "SOTA"},
		{field: spec.WwffRefField, prog: "WWFF"},
	}
	for _, p := range progs {
		fmt.Fprintf(res, progfmt, p.field.Name, spec.SigInfoField.Name, spec.SigField.Name, p.prog)
		fmt.Fprintf(res, progfmt, "MY_"+p.field.Name, spec.MySigInfoField.Name, spec.MySigField.Name, p.prog)
	}
	return res.String()
}

func runInfer(ctx *Context, args []string) error {
	cctx := ctx.CommandCtx.(*InferContext)
	todo := make([]string, len(cctx.Fields))
	for i, f := range cctx.Fields {
		todo[i] = strings.ToUpper(f)
		if inferrers[todo[i]] == nil {
			return fmt.Errorf("don't know how to infer field %s\n%s", todo[i], helpInfer())
		}
	}
	out := adif.NewLogfile()
	acc := accumulator{Out: out, Ctx: ctx}
	for _, f := range filesOrStdin(args) {
		l, err := acc.read(f)
		if err != nil {
			return err
		}
		updateFieldOrder(out, l.FieldOrder)
		for _, r := range l.Records {
			did := make([]string, 0, len(todo))
			for _, t := range todo {
				if inferrers[t] != nil {
					if f, ok := r.Get(t); !ok || f.Value == "" {
						if inferrers[t](r, t) {
							did = append(did, t)
						}
					}
				}
			}
			if cctx.CommentLog && len(did) > 0 {
				c := "adif-multitool infered value for " + strings.Join(did, ", ")
				if r.GetComment() == "" {
					r.SetComment(c)
				} else {
					r.SetComment(r.GetComment() + "\n" + c)
				}
			}
			out.AddRecord(r)
		}
	}
	if err := acc.prepare(); err != nil {
		return err
	}
	return write(ctx, out)
}

func inferBand(r *adif.Record, name string) bool {
	freqname := spec.FreqField.Name
	if name == spec.BandRxField.Name {
		freqname = spec.FreqRxField.Name
	}
	f, ok := r.Get(freqname)
	if !ok || f.Value == "" {
		return false
	}
	freq, err := strconv.ParseFloat(f.Value, 64)
	if err != nil {
		return false
	}
	for _, b := range spec.BandEnumeration.Values {
		bb := b.(spec.BandEnum)
		min, err := strconv.ParseFloat(bb.LowerFreqMhz, 64)
		if err != nil {
			return false
		}
		max, err := strconv.ParseFloat(bb.UpperFreqMhz, 64)
		if err != nil {
			return false
		}
		if min <= freq && freq <= max {
			r.Set(adif.Field{Name: name, Value: bb.Band})
			return true
		}
	}
	return false
}

func inferCountry(r *adif.Record, name string) bool {
	my := func(s string) string { return s }
	if strings.HasPrefix(name, "MY_") {
		my = func(s string) string { return "MY_" + s }
	}
	code, ok := r.Get(my(spec.DxccField.Name))
	if !ok || code.Value == "" || code.Value == "0" {
		return false
	}
	for _, e := range spec.DxccEntityCodeEnumeration.Value(code.Value) {
		ee := e.(spec.DxccEntityCodeEnum)
		if ee.Deleted == "true" {
			continue
		}
		r.Set(adif.Field{Name: name, Value: ee.EntityName})
		return true
	}
	return false
}

func inferDxcc(r *adif.Record, name string) bool {
	my := func(s string) string { return s }
	if strings.HasPrefix(name, "MY_") {
		my = func(s string) string { return "MY_" + s }
	}
	c, ok := r.Get(my(spec.CountryField.Name))
	if !ok || c.Value == "" {
		return false
	}
	for _, e := range spec.CountryEnumeration.Value(c.Value) {
		ee := e.(spec.CountryEnum)
		if ee.Deleted == "true" {
			continue
		}
		r.Set(adif.Field{Name: name, Value: ee.EntityCode})
		return true
	}
	return false
}

func inferMode(r *adif.Record, name string) bool {
	s, ok := r.Get(spec.SubmodeField.Name)
	if !ok || s.Value == "" {
		return false
	}
	for _, e := range spec.SubmodeEnumeration.Value(s.Value) {
		ee := e.(spec.SubmodeEnum)
		r.Set(adif.Field{Name: name, Value: ee.Mode})
		return true
	}
	return false
}

func inferSigInfo(r *adif.Record, name string) bool {
	my := func(s string) string { return s }
	if strings.HasPrefix(name, "MY_") {
		my = func(s string) string { return "MY_" + s }
	}
	islota, iotaok := r.Get(my(spec.IotaField.Name))
	pota, potaok := r.Get(my(spec.PotaRefField.Name))
	sota, sotaok := r.Get(my(spec.SotaRefField.Name))
	wwff, wwffok := r.Get(my(spec.WwffRefField.Name))
	var unknownSig = false
	if f, ok := r.Get(my(spec.SigField.Name)); ok && f.Value != "" {
		var v adif.Field
		switch strings.ToUpper(f.Value) {
		default:
			unknownSig = true
		case "IOTA":
			v = islota
		case "POTA":
			v = pota
		case "SOTA":
			v = sota
		case "WWFF":
			v = wwff
		}
		if v.Value != "" {
			r.Set(adif.Field{Name: my(spec.SigInfoField.Name), Value: v.Value})
			return true
		}
	} else if !unknownSig {
		// SIG/MY_SIG not set, guess which program was active
		var v adif.Field
		var sig string
		var got int
		if iotaok && islota.Value != "" {
			v = islota
			sig = "IOTA"
			got++
		}
		if potaok && pota.Value != "" {
			v = pota
			sig = "POTA"
			got++
		}
		if sotaok && sota.Value != "" {
			v = sota
			sig = "SOTA"
			got++
		}
		if wwffok && wwff.Value != "" {
			v = wwff
			sig = "WWFF"
			got++
		}
		if got == 1 {
			r.Set(adif.Field{Name: my(spec.SigField.Name), Value: sig})
			r.Set(adif.Field{Name: my(spec.SigInfoField.Name), Value: v.Value})
			return true
		}
	}
	return false
}

func inferProgramRef(wantSig string) inferrer {
	return func(r *adif.Record, name string) bool {
		my := func(s string) string { return s }
		if strings.HasPrefix(name, "MY_") {
			my = func(s string) string { return "MY_" + s }
		}
		siginfo, ok := r.Get(my(spec.SigInfoField.Name))
		if !ok || siginfo.Value == "" {
			return false
		}
		sig, ok := r.Get(my(spec.SigField.Name))
		if !ok || sig.Value == "" {
			return false
		}
		if !strings.EqualFold(sig.Value, wantSig) {
			return false
		}
		f, ok := spec.FieldNamed(name)
		if !ok {
			return false
		}
		v := spec.TypeValidators[f.Type.Name]
		if v(siginfo.Value, f, spec.ValidationContext{}).Validity != spec.Valid {
			return false
		}
		r.Set(adif.Field{Name: name, Value: siginfo.Value})
		return true
	}
}

func inferStation(r *adif.Record, name string) bool {
	var order []spec.Field
	if name == spec.OperatorField.Name {
		order = []spec.Field{spec.GuestOpField}
	} else if name == spec.StationCallsignField.Name {
		order = []spec.Field{spec.OperatorField, spec.GuestOpField}
	} else if name == spec.OwnerCallsignField.Name {
		order = []spec.Field{spec.StationCallsignField, spec.OperatorField, spec.GuestOpField}
	}
	for _, f := range order {
		if v, ok := r.Get(f.Name); ok && v.Value != "" {
			r.Set(adif.Field{Name: name, Value: v.Value})
			return true
		}
	}
	return false
}

var usCountyPattern = regexp.MustCompile(`^[A-Z]{2},[A-Za-z '.-]+$`) // doesn't match county-line lists

func inferUSCounty(r *adif.Record, name string) bool {
	if f, ok := r.Get(name); ok && f.Value != "" {
		return false
	}
	my := func(s string) string { return s }
	if strings.HasPrefix(name, "MY_") {
		my = func(s string) string { return "MY_" + s }
	}
	if dx, ok := r.Get(my(spec.DxccField.Name)); ok && dx.Value != "" {
		if !spec.CountryCodeUSA.IncludesDXCC(dx.Value) {
			return false // not a US contact
		}
	}
	var src string
	switch strings.ToUpper(name) {
	case spec.UsacaCountiesField.Name, spec.MyUsacaCountiesField.Name:
		src = my(spec.CntyField.Name)
	case spec.CntyField.Name, spec.MyCntyField.Name:
		src = my(spec.UsacaCountiesField.Name)
	default:
		return false
	}
	sub, ok := r.Get(src)
	if !ok || sub.Value == "" || !usCountyPattern.MatchString(sub.Value) {
		return false
	}
	r.Set(adif.Field{Name: name, Value: sub.Value})
	return true
}

func inferLatLon(r *adif.Record, name string) bool {
	my := func(s string) string { return s }
	if strings.HasPrefix(name, "MY_") {
		my = func(s string) string { return "MY_" + s }
	}
	f, ok := r.Get(my(spec.GridsquareField.Name))
	if !ok || f.Value == "" {
		return false
	}
	gs := f.Value
	if ext, ok := r.Get(my(spec.GridsquareExtField.Name)); ok {
		gs += ext.Value
	}
	lat, lon, err := parseMaidenhead(gs)
	if err != nil {
		return false
	}
	if name == spec.LatField.Name || name == spec.MyLatField.Name {
		s, err := formatLatitude(lat)
		if err != nil {
			return false
		}
		r.Set(adif.Field{Name: name, Value: s})
		return true
	}
	if name == spec.LonField.Name || name == spec.MyLonField.Name {
		s, err := formatLongitude(lon)
		if err != nil {
			return false
		}
		r.Set(adif.Field{Name: name, Value: s})
		return true
	}
	return false
}

func inferGridsquare(r *adif.Record, name string) bool {
	my := func(s string) string { return s }
	if strings.HasPrefix(name, "MY_") {
		my = func(s string) string { return "MY_" + s }
	}
	var latf, lonf string
	if f, ok := r.Get(my(spec.LatField.Name)); ok && f.Value != "" {
		latf = f.Value
	} else {
		return false
	}
	if f, ok := r.Get(my(spec.LonField.Name)); ok && f.Value != "" {
		lonf = f.Value
	} else {
		return false
	}
	lat, lon, err := parseADIFCoordinates(latf, lonf)
	if err != nil {
		fmt.Println(err)
		return false
	}
	// Maidenhead locator uses positive values from south pole and antiprime meridian
	lat += 90
	lon += 180
	lons := maidenheadSlice{rem: lon, scale: 360}
	lats := maidenheadSlice{rem: lat, scale: 180}
	var gs strings.Builder
	// first pair is divided into 18 letters, 20° longitude, 10º latitude
	gs.WriteRune('A' + rune(lons.split(18)))
	gs.WriteRune('A' + rune(lats.split(18)))
	// second pair is divided into 10 digits, 2º longitude, 1º latitude
	gs.WriteRune('0' + rune(lons.split(10)))
	gs.WriteRune('0' + rune(lats.split(10)))
	// third pair is divided into 24 letters, 5' longitude, 2.5' latitude
	gs.WriteRune('a' + rune(lons.split(24)))
	gs.WriteRune('a' + rune(lats.split(24)))
	// fourth pair is divided into 10 letters, 30" longitude, 15" latitude
	gs.WriteRune('0' + rune(lons.split(10)))
	gs.WriteRune('0' + rune(lats.split(10)))
	// fifth pair is divided into 24 letters, 1.2" longitude (≈36m), 0.625" latitude (≈19m)
	gs.WriteRune('a' + rune(lons.split(24)))
	gs.WriteRune('a' + rune(lats.split(24)))
	// fifth pair is divided into 10 digits, 0.12" longitude (≈3.6m), 0.0625" latitude (≈1.9m)
	gs.WriteRune('0' + rune(lons.split(10)))
	gs.WriteRune('0' + rune(lats.split(10)))
	if strings.HasSuffix(name, "_EXT") {
		r.Set(adif.Field{Name: name, Value: gs.String()[8:]})
	} else {
		r.Set(adif.Field{Name: name, Value: gs.String()[0:8]})
	}
	return true
}

type maidenheadSlice struct{ rem, scale float64 }

func (s *maidenheadSlice) split(num int) int {
	div := s.scale / float64(num)
	d := int(s.rem / div)
	s.rem = s.rem - float64(d)*div
	s.scale = s.scale / float64(num)
	return d
}

func parseADIFCoordinates(lat, lon string) (float64, float64, error) {
	if lat == "" || lon == "" {
		return 0, 0, fmt.Errorf("empty coordinate %q, %q", lat, lon)
	}
	pat := "%c%03d %6f"
	var latdir, londir rune
	var latdeg, londeg int
	var latmin, lonmin float64
	if n, err := fmt.Sscanf(lat, pat, &latdir, &latdeg, &latmin); err != nil || n != 3 {
		return 0, 0, fmt.Errorf("could not parse latitude %q: %w", lat, err)
	}
	if n, err := fmt.Sscanf(lon, pat, &londir, &londeg, &lonmin); err != nil || n != 3 {
		return 0, 0, fmt.Errorf("could not parse longitude %q: %w", lon, err)
	}
	var latsign, lonsign float64
	switch latdir {
	case 'N', 'n':
		latsign = 1
	case 'S', 's':
		latsign = -1
	default:
		return 0, 0, fmt.Errorf("invalid latitude direction %c", latdir)
	}
	switch londir {
	case 'E', 'e':
		lonsign = 1
	case 'W', 'w':
		lonsign = -1
	default:
		return 0, 0, fmt.Errorf("invalid longitude direction %c", londir)
	}
	if latmin > 60 || lonmin > 60 || latmin < 0 || lonmin < 0 {
		return 0, 0, fmt.Errorf("minutes out of range: %s, %s", lat, lon)
	}
	rlat := latsign * (float64(latdeg) + latmin/60.0)
	rlon := lonsign * (float64(londeg) + lonmin/60.0)
	if math.Abs(rlat) > 90 || math.Abs(rlon) > 180 {
		return 0, 0, fmt.Errorf("coordinates out of range %s, %s", lat, lon)
	}
	return rlat, rlon, nil
}

func parseMaidenhead(gs string) (lat float64, lon float64, err error) {
	if len(gs) < 2 {
		err = errors.New("empty string")
		return
	}
	invalid := fmt.Errorf("invalid format %q", gs)
	if len(gs)%2 != 0 {
		err = invalid
		return
	}
	gs = strings.ToUpper(gs)
	lonscale := 360.0
	latscale := 180.0
	sizes := []float64{18, 10, 24, 10, 24, 10}
	for i, size := range sizes {
		if len(gs) <= i*2 {
			break
		}
		lonr := rune(gs[i*2])
		latr := rune(gs[i*2+1])
		var lonval, latval int
		if size == 10 {
			if lonr < '0' || latr < '0' || lonr > '9' || latr > '9' {
				err = invalid
				return
			}
			lonval = int(lonr - '0')
			latval = int(latr - '0')
		} else {
			if lonr < 'A' || latr < 'A' || lonr > 'A'+rune(size) || latr > 'A'+rune(size) {
				err = invalid
				return
			}
			lonval = int(lonr - 'A')
			latval = int(latr - 'A')
		}
		lonscale /= size
		latscale /= size
		lon += lonscale * float64(lonval)
		lat += latscale * float64(latval)
	}
	// center result in remaining square, then shift to negative/positive coords
	lon += lonscale / 2
	lat += latscale / 2
	lon -= 180
	lat -= 90
	return
}

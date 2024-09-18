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

package adif

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// CabrilloIO configures the headers and QSO inference for conversion to and
// from the Cabrillo format.  Most fields configure the value of a header with
// the same name.  Categories maps CATEGORY header names to value, e.g. "TIME"
// to "6-HOURS".  See https://wwrof.org/cabrillo/cabrillo-v3-header/ for
// details about header values.
type CabrilloIO struct {
	Callsign, Contest, Club, CreatedBy, Email,
	GridLocator, Location, Name, Address, Soapbox string
	Operators                              []string
	LowPowerMax, QRPPowerMax, ClaimedScore int
	MinReportedOfftime                     time.Duration
	Categories                             map[string]string
	MyExchange, TheirExchange, ExtraFields CabrilloFieldList
	TabDelimiter                           bool
}

func NewCabrilloIO() *CabrilloIO {
	return &CabrilloIO{
		LowPowerMax:   100,
		QRPPowerMax:   5,
		Categories:    make(map[string]string),
		MyExchange:    make(CabrilloFieldList, 0),
		TheirExchange: make(CabrilloFieldList, 0),
		ExtraFields:   make(CabrilloFieldList, 0),
	}
}

func (_ *CabrilloIO) String() string { return "cabrillo" }

func (o *CabrilloIO) Read(in io.Reader) (*Logfile, error) {
	headers := make(map[string]string)
	s := bufio.NewScanner(in)
	readLine := func() (k, v string, err error) {
		if !s.Scan() {
			err = s.Err()
			if err == nil {
				err = io.EOF
			}
			return
		}
		line := s.Text()
		for strings.TrimSpace(line) == "" {
			if !s.Scan() {
				err = s.Err()
				if err == nil {
					err = io.EOF
				}
				return
			}
			line = s.Text()
		}
		k, v, ok := strings.Cut(line, ":")
		if !ok {
			err = fmt.Errorf("invalid Cabrillo line %q", line)
			return
		}
		v = strings.Trim(v, " ")
		return
	}
	start, version, err := readLine()
	if err != nil {
		return nil, err
	}
	if start != "START-OF-LOG" {
		return nil, fmt.Errorf("Cabrillo file does not start with START-OF-LOG: %q", s.Text())
	}
	if v, err := strconv.ParseFloat(version, 64); err != nil || v != 3.0 {
		return nil, fmt.Errorf("don't know how to parse Cabrillo version %q", version)
	}
	l := NewLogfile()
	conf := o.toConfig()
	for {
		k, v, err := readLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, errors.New("got EOF before Cabrillo END-OF-LOG: line")
			}
			return nil, err
		}
		if k == "END-OF-LOG" {
			break
		}
		if k == "QSO" || k == "X-QSO" {
			r, err := conf.toADIF(v)
			if err != nil {
				return nil, err
			}
			if k == "X-QSO" {
				r.Set(Field{Name: "APP_CABRILLO_XQSO", Value: "Y", Type: TypeBoolean})
			}
			l.AddRecord(r)
		} else if v != "" && k != "X-INSTRUCTIONS" && k != "X-Q" {
			if vv, ok := headers[k]; ok {
				headers[k] = vv + "\n" + v
			} else {
				headers[k] = v
			}
		}
	}
	if k, v, err := readLine(); err == nil {
		return nil, fmt.Errorf(`got data after END-OF-LOG: line "%s: %s"`, k, v)
	} else if !errors.Is(err, io.EOF) {
		return nil, err
	}
	fromHeaders := make([]Field, 0)
	if v := headers["OPERATORS"]; v != "" && !strings.ContainsAny(v, " ,&@") {
		fromHeaders = append(fromHeaders, Field{Name: "OPERATOR", Value: v, Type: TypeString})
	}
	if v := headers["CONTEST"]; v != "" {
		fromHeaders = append(fromHeaders, Field{Name: "CONTEST_ID", Value: v, Type: TypeString})
	}
	if v := headers["GRID-LOCATOR"]; v != "" && !strings.ContainsAny(v, " ,&") {
		fromHeaders = append(fromHeaders, Field{Name: "GRIDSQUARE", Value: v, Type: TypeString})
	}
	for _, f := range fromHeaders {
		for _, r := range l.Records {
			if v, ok := r.Get(f.Name); !ok || v.Value != "" {
				r.Set(f)
			}
		}
	}
	headorder := maps.Keys(headers)
	slices.Sort(headorder)
	for _, k := range headorder {
		v := headers[k]
		l.Header.Set(Field{Name: "APP_CABRILLO_" + strings.ReplaceAll(k, "-", "_"), Value: v, Type: TypeString})
	}
	return l, nil
}

func (o *CabrilloIO) Write(l *Logfile, out io.Writer) error {
	var qlines [][2]string
	var err error
	qlines, err = o.toConfig().toLines(l)
	if err != nil {
		return err
	}
	headers := make(map[string]string)
	for _, f := range l.Header.Fields() {
		if h, ok := cutPrefix(f.Name, "APP_CABRILLO_"); ok {
			headers[strings.Replace(h, "_", "-", -1)] = f.Value
		}
	}
	w := bufio.NewWriter(out)
	writeLine := func(k, v string) error {
		for _, line := range splitLines.Split(v, -1) {
			if _, err := w.WriteString(k); err != nil {
				return err
			}
			if _, err := w.WriteRune(':'); err != nil {
				return err
			}
			if line != "" {
				if _, err := w.WriteRune(' '); err != nil {
					return err
				}
				if _, err := w.WriteString(line); err != nil {
					return err
				}
			}
			if _, err := w.WriteRune('\n'); err != nil {
				return err
			}
		}
		return nil
	}
	setHeader := func(hname, val string) {
		if val != "" || headers[hname] == "" {
			headers[hname] = val
		}
	}
	setSummaryField := func(hname, fname, priority string) {
		val := priority
		if val == "" {
			m := fieldValues(l, fname)
			if len(m) == 1 {
				val = maps.Keys(m)[0]
			} else {
				vals := make([]string, 0, len(m))
				for s, c := range m {
					vals = append(vals, fmt.Sprintf("%s (%d records)", s, c))
				}
				sort.Strings(vals)
				val = strings.Join(vals, " ")
			}
		}
		setHeader(hname, val)
	}
	setHeader("SOAPBOX", o.Soapbox)
	setSummaryField("CONTEST", "CONTEST_ID", o.Contest)
	setSummaryField("CALLSIGN", "STATION_CALLSIGN", o.Callsign)
	setHeader("CLUB", o.Club)
	ops := o.Operators
	if len(ops) == 0 {
		ops = maps.Keys(fieldValues(l, "OPERATOR"))
	}
	setHeader("OPERATORS", strings.Join(ops, ", "))
	setSummaryField("NAME", "MY_NAME", o.Name)
	setHeader("EMAIL", o.Email)
	setHeader("ADDRESS", o.Address)
	setSummaryField("GRID-LOCATOR", "MY_GRIDSQUARE", o.GridLocator)
	if headers["LOCATION"] == "" {
		// some contests use other location values, but ARRL section is a good first guess
		setSummaryField("LOCATION", "MY_ARRL_SECT", o.Location)
	} else if o.Location != "" {
		headers["LOCATION"] = o.Location
	}
	if s, ok := headers["CLAIMED-SCORE"]; !ok || s == "" || s == "0" {
		setHeader("CLAIMED-SCORE", fmt.Sprintf("%d", o.ClaimedScore))
	}
	// TODO compute offtimes
	setHeader("OFFTIME", "")
	cats := o.getCategories(l)
	for k, v := range cats {
		if strings.HasPrefix(k, "X-") {
			setHeader(k, v)
		} else {
			setHeader("CATEGORY-"+k, v)
		}
	}
	headerOrder := []string{"SOAPBOX", "CONTEST", "CALLSIGN", "CLUB", "OPERATORS", "NAME", "EMAIL", "ADDRESS", "ADDRESS-CITY", "ADDRESS-STATE-PROVINCE", "ADDRESS-POSTALCODE", "ADDRESS-COUNTRY", "GRID-LOCATOR", "LOCATION", "CLAIMED-SCORE", "OFFTIME"}
	categoryOrder := maps.Keys(CabrilloCategoryValues)
	sort.Strings(categoryOrder)
	for _, c := range categoryOrder {
		headerOrder = append(headerOrder, "CATEGORY-"+c)
	}
	if err := writeLine("START-OF-LOG", "3.0"); err != nil {
		return err
	}
	for _, s := range cabrilloInstructions {
		if err := writeLine("X-INSTRUCTIONS", s); err != nil {
			return err
		}
	}
	if err := writeLine("CREATED-BY", o.CreatedBy); err != nil {
		return err
	}
	wrote := map[string]bool{"CREATED-BY": true}
	for _, h := range headerOrder {
		if err := writeLine(h, headers[h]); err != nil {
			return err
		}
		wrote[h] = true
	}
	for h, v := range headers {
		if !wrote[h] {
			if err := writeLine(h, v); err != nil {
				return err
			}
		}
	}
	if err := writeLine("X-INSTRUCTIONS", "See contest rules for expected category values"); err != nil {
		return err
	}
	for _, q := range qlines {
		if err := writeLine(q[0], q[1]); err != nil {
			return err
		}
	}
	if err := writeLine("END-OF-LOG", ""); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}
	return nil
}

func (o *CabrilloIO) getCategories(l *Logfile) map[string]string {
	cats := make(map[string]string)
	for k, v := range o.Categories {
		cats[strings.ToUpper(k)] = v
	}
	if cats["MODE"] == "" {
		modes := fieldValues(l, "MODE")
		var phone, fm, cw, rtty, digi int
		for m := range modes {
			switch m {
			case "SSB", "AM", "DIGITALVOICE":
				phone++
			case "FM":
				fm++
			case "CW":
				cw++
			case "RTTY":
				rtty++
			default:
				digi++
			}
		}
		var mode string
		if fm == len(modes) {
			mode = "FM"
		} else if phone+fm == len(modes) {
			mode = "SSB"
		} else if cw == len(modes) {
			mode = "CW"
		} else if rtty == len(modes) {
			mode = "RTTY"
		} else if digi+rtty == len(modes) {
			mode = "DIGI"
		} else {
			mode = "MIXED"
		}
		cats["MODE"] = mode
	}
	if cats["BAND"] == "" {
		bands := fieldValues(l, "BAND")
		var band string
		if len(bands) == 1 {
			b := maps.Keys(bands)[0]
			band = cabrilloBandCategories[strings.ToLower(b)]
			if band == "" {
				band = cabrilloBandsRev[strings.ToLower(b)]
			}
		} else if len(bands) > 1 {
			band = "ALL"
		}
		cats["BAND"] = band
	}
	if cats["POWER"] == "" {
		pwrs := fieldValues(l, "TX_PWR")
		var max float64
		for pwr := range pwrs {
			if p, err := strconv.ParseFloat(pwr, 64); err == nil {
				if p > max {
					max = p
				}
			}
		}
		if max > 0 {
			if max <= float64(o.QRPPowerMax) && o.QRPPowerMax > 0 {
				cats["POWER"] = "QRP"
			} else if max <= float64(o.LowPowerMax) && o.LowPowerMax > 0 {
				cats["POWER"] = "LOW"
			} else {
				cats["POWER"] = "HIGH"
			}
			cats["X-MAX-POWER"] = strconv.FormatFloat(max, 'f', -1, 64)
		}
	}
	return cats
}

func (o *CabrilloIO) toConfig() cabrilloConfig {
	myCall := cabFieldMyCall
	if o.Callsign != "" {
		myCall.Default = o.Callsign
	}
	c := cabrilloConfig{
		useTabs:  o.TabDelimiter,
		coreLen:  len(cabCoreFields),
		myLen:    len(o.MyExchange) + 1,
		theirLen: len(o.TheirExchange) + 1,
		extraLen: len(o.ExtraFields)}
	// TODO use slices.Concat in go1.22
	c.fields = append(c.fields, cabCoreFields...)
	c.fields = append(c.fields, myCall)
	c.fields = append(c.fields, o.MyExchange...)
	c.fields = append(c.fields, cabFieldTheirCall)
	c.fields = append(c.fields, o.TheirExchange...)
	c.fields = append(c.fields, o.ExtraFields...)
	return c
}

func mhzToKhz(mhz string) string {
	pieces := strings.Split(mhz, ".")
	if len(pieces) == 1 {
		return mhz + "000" // 14 MHz -> 14000 kHz
	}
	switch len(pieces[1]) {
	case 0:
		return pieces[0] + "000"
	case 1:
		return pieces[0] + pieces[1] + "00"
	case 2:
		return pieces[0] + pieces[1] + "0"
	case 3:
		return pieces[0] + pieces[1]
	default:
		return pieces[0] + pieces[1][0:3] + "." + pieces[1][3:]
	}
}

type bandRange struct {
	cabrilloName, adifName string
	lo, hi                 float64
}

func (b bandRange) compare(f float64) int {
	if b.lo <= f && b.hi >= f {
		return 0
	}
	if b.hi < f {
		return -1
	}
	return 1
}

func findBandFreq(freq float64) (bandRange, bool) {
	i, ok := slices.BinarySearchFunc(cabrilloRanges, freq, bandRange.compare)
	if !ok {
		return bandRange{}, false
	}
	return cabrilloRanges[i], true
}

var (
	cabrilloRanges = []bandRange{
		{cabrilloName: "1800", adifName: "160m", lo: 1.8, hi: 2.0},
		{cabrilloName: "3500", adifName: "80m", lo: 3.5, hi: 4.0},
		{cabrilloName: "7000", adifName: "40m", lo: 7.0, hi: 7.3},
		{cabrilloName: "14000", adifName: "20m", lo: 14.0, hi: 14.35},
		{cabrilloName: "21000", adifName: "15m", lo: 21.0, hi: 21.45},
		{cabrilloName: "28000", adifName: "10m", lo: 28.0, hi: 29.7},
		{cabrilloName: "50", adifName: "6m", lo: 50.0, hi: 54.0},
		{cabrilloName: "70", adifName: "4m", lo: 70.0, hi: 71.0},
		{cabrilloName: "144", adifName: "2m", lo: 144.0, hi: 148.0},
		{cabrilloName: "222", adifName: "1.25m", lo: 222.0, hi: 225.0},
		{cabrilloName: "432", adifName: "70cm", lo: 420.0, hi: 450.0},
		{cabrilloName: "902", adifName: "33cm", lo: 902.0, hi: 928.0},
		{cabrilloName: "1.2G", adifName: "23cm", lo: 1240.0, hi: 1300.0},
		{cabrilloName: "2.3G", adifName: "13cm", lo: 2300.0, hi: 2450.0},
		{cabrilloName: "3.4G", adifName: "9cm", lo: 3300.0, hi: 3500.0},
		{cabrilloName: "5.7G", adifName: "6cm", lo: 5650.0, hi: 5925.0},
		{cabrilloName: "10G", adifName: "3cm", lo: 10000.0, hi: 10500.0},
		{cabrilloName: "24G", adifName: "1.25cm", lo: 24000.0, hi: 24250.0},
		{cabrilloName: "47G", adifName: "6mm", lo: 47000.0, hi: 47200.0},
		{cabrilloName: "75G", adifName: "4mm", lo: 75500.0, hi: 81000.0},
		{cabrilloName: "122G", adifName: "2.5mm", lo: 119980.0, hi: 123000.0},
		{cabrilloName: "134G", adifName: "2mm", lo: 134000.0, hi: 149000.0},
		{cabrilloName: "241G", adifName: "1mm", lo: 241000.0, hi: 250000.0},
		{cabrilloName: "LIGHT", adifName: "submm", lo: 300000.0, hi: 7500000.0},
	}
	cabrilloBands    = map[string]string{}
	cabrilloBandsRev = map[string]string{}
	// CATEGORY-BAND uses wavelength in meters up to 2m, but QSOs use the lowest frequency on the band.
	// Both QSOs and the category use the same values for 222 MHz and up.
	cabrilloBandCategories = map[string]string{
		"160m": "160M",
		"80m":  "80M",
		"40m":  "40M",
		"20m":  "20M",
		"15m":  "15M",
		"10m":  "10M",
		"6m":   "6M",
		"4m":   "4M",
		"2m":   "2M",
	}

	cabrilloModes = map[string]string{
		"CW": "CW",
		"PH": "SSB", // most likely; could be AM or DIGITALVOICE
		"FM": "FM",
		"RY": "RTTY",
		"DG": "DIGITAL", // not an ADIF mode, but don't know the actual mode
	}

	cabrilloInstructions = []string{
		"Fill out headers following contest instructions",
		"Delete any unnecessary headers",
		"Double-check QSO lines, keeping columns in order",
		"Report bugs at https://github.com/flwyd/adif-multitool",
	}

	// Allowed values for Cabrillo categories.
	// See https://wwrof.org/cabrillo/cabrillo-v3-header/
	CabrilloCategoryValues = map[string][]string{
		"ASSISTED": {"ASSISTED", "NON-ASSISTED"},
		"BAND": {
			"ALL",
			"160M",
			"80M",
			"40M",
			"20M",
			"15M",
			"10M",
			"6M",
			"4M",
			"2M",
			"222",
			"432",
			"902",
			"1.2G",
			"2.3G",
			"3.4G",
			"5.7G",
			"10G",
			"24G",
			"47G",
			"75G",
			"122G",
			"134G",
			"241G",
			"LIGHT",
			"VHF-3-BAND",
			"VHF-FM-ONLY",
		},
		"MODE":     {"CW", "DIGI", "FM", "RTTY", "SSB", "MIXED"},
		"OPERATOR": {"SINGLE-OP", "MULTI-OP", "CHECKLOG"},
		"OVERLAY":  {"CLASSIC", "ROOKIE", "TB-WIRES", "YOUTH", "NOVICE-TECH", "YL"},
		"POWER":    {"HIGH", "LOW", "QRP"},
		"STATION": {
			"DISTRIBUTED",
			"FIXED",
			"MOBILE",
			"PORTABLE",
			"ROVER",
			"ROVER-LIMITED",
			"ROVER-UNLIMITED",
			"EXPEDITION",
			"HQ",
			"SCHOOL",
			"EXPLORER",
		},
		"TIME":        {"6-HOURS", "8-HOURS", "12-HOURS", "24-HOURS"},
		"TRANSMITTER": {"ONE", "TWO", "LIMITED", "UNLIMITED", "SWL"},
	}

	cabrilloCallPat = regexp.MustCompile("[A-Za-z0-9]{3,}(/[A-Za-z0-9]+)?")
	cabrilloRSTPat  = regexp.MustCompile("[1-9]{2,3}")
	splitLines      = regexp.MustCompile(`[\r\n]+`)
)

func init() {
	for _, b := range cabrilloRanges {
		cabrilloBands[b.cabrilloName] = b.adifName
		cabrilloBandsRev[b.adifName] = b.cabrilloName
	}
}

/*
	CabrilloField represents a field in a Cabrillo file which will appear in a column.

The data in a CabrilloField can come from one field, a series of ADIF fields,
or be set to a default value, e.g. an exchange used throughout a contest.  A
field optionally has a header, shown only for informational purposes.  Fields
are required by default, but may be made optional, in which case one or more
hyphens will be used in the output if the ADIF fields are all empty.
*/
type CabrilloField struct {
	TryFields  []string
	Default    string
	Header     string
	AllowEmpty bool
}

const CabrilloFieldExample = "header:field_a/field_b?=default"

var cabrilloFieldSyntax = regexp.MustCompile(`^([^:]*:)?(\w+(?:/\w+)*)?(\??)(=\S+)?$`)

func (c *CabrilloField) String() string {
	var s strings.Builder
	if c.Header != "" {
		s.WriteString(c.Header)
		s.WriteRune(':')
	}
	for i, x := range c.TryFields {
		if i > 0 {
			s.WriteRune('/')
		}
		s.WriteString(x)
	}
	if c.Default != "" {
		s.WriteRune('=')
		s.WriteString(c.Default)
	}
	if c.AllowEmpty {
		s.WriteRune('?')
	}
	return s.String()
}

func parseCabrilloField(v string) (CabrilloField, error) {
	g := cabrilloFieldSyntax.FindStringSubmatch(v)
	if g == nil {
		return CabrilloField{}, fmt.Errorf("%q did not match format %s", v, CabrilloFieldExample)
	}
	try := strings.Split(g[2], "/")
	def := strings.TrimPrefix(g[4], "=")
	if len(try) == 0 && def == "" {
		return CabrilloField{}, fmt.Errorf("%v did not specify field name(s) or =default value", v)
	}
	return CabrilloField{Header: strings.TrimSuffix(g[1], ":"), TryFields: try, AllowEmpty: g[3] == "?", Default: def}, nil
}

func (c *CabrilloField) fromADIF(r *Record) (string, error) {
	var v string
	for _, t := range c.TryFields {
		if v != "" {
			break
		}
		t = strings.ToUpper(t)
		switch t {
		default:
			if f, ok := r.Get(t); ok && f.Value != "" {
				v = strings.ReplaceAll(strings.TrimSpace(f.Value), " ", "_")
			}
		case "QSO_DATE", "QSO_DATE_OFF":
			d, err := r.ParseDate(t)
			if err != nil {
				return "", err
			}
			v = d.Format("2006-01-02")
		case "TIME_ON", "TIME_OFF":
			d, err := r.ParseTime(t)
			if err != nil {
				return "", err
			}
			v = d.Format("1504")
		case "FREQ":
			if n, err := r.ParseFloat(t); err != nil {
				continue
			} else if n >= 30 { // Cabrillo uses band names above 30 MHz
				b, ok := findBandFreq(n)
				if !ok {
					continue
				}
				v = b.cabrilloName
			} else {
				f, _ := r.Get(t) // ParseFloat already determined it's set
				// string-to-string to avoid floating point precision issues
				v = mhzToKhz(f.Value)
			}
		case "BAND":
			if f, ok := r.Get(t); ok && f.Value != "" {
				if b, ok := cabrilloBandsRev[strings.ToLower(f.Value)]; !ok {
					return "", fmt.Errorf("invalid band %q for Cabrillo", f.Value)
				} else {
					v = b
				}
			}
		case "MODE":
			if f, ok := r.Get(t); ok && f.Value != "" {
				switch strings.ToUpper(f.Value) {
				case "CW":
					v = "CW"
				case "RTTY", "RTTYM":
					v = "RY"
				case "SSB", "AM", "DIGITALVOICE":
					v = "PH" // TODO verify DIGITALVOICE counts as phone
				case "FM":
					v = "FM"
				default:
					v = "DG" // most modes are digital
				}
			}
		}
	}
	if v == "" { // no TryFields value was present and/or valid
		v = c.Default
	}
	if v == "" { // Default was empty
		if !c.AllowEmpty {
			return "", fmt.Errorf("missing %s in %s", strings.Join(c.TryFields, ","), r)
		}
		v = strings.Repeat("-", maxInt(len(c.Header), 1))
	}
	return v, nil
}

func (c *CabrilloField) toADIF(val string) (Field, error) {
	if val == "" {
		val = c.Default
	}
	var f Field
	var fieldErr error
	for _, t := range c.TryFields {
		switch strings.ToUpper(t) {
		case "QSO_DATE", "QSO_DATE_OFF":
			if d, err := time.ParseInLocation("2006-01-02", val, time.UTC); err != nil {
				fieldErr = fmt.Errorf("invalid Cabrillo date %q: %w", val, err)
			} else {
				f = Field{Name: t, Value: d.Format("20060102"), Type: TypeDate}
			}
		case "TIME_ON", "TIME_OFF":
			if _, err := time.ParseInLocation("1504", val, time.UTC); err != nil {
				fieldErr = fmt.Errorf("invalid hhmm time %q: %w", val, err)
			} else {
				f = Field{Name: t, Value: val, Type: TypeTime}
			}
		case "FREQ":
			if b, ok := cabrilloBands[strings.ToUpper(val)]; ok {
				f = Field{Name: "BAND", Value: b, Type: TypeEnumeration}
				break
			}
			if khz, err := strconv.ParseFloat(val, 64); err != nil {
				fieldErr = fmt.Errorf("invalid frequency %q kHz: %w", val, err)
			} else if khz/1000 < cabrilloRanges[0].lo {
				fieldErr = fmt.Errorf("frequency %s kHz too low", val)
			} else {
				f = Field{Name: t, Value: strconv.FormatFloat(khz/1000.0, 'f', -1, 64), Type: TypeNumber}
			}
		case "BAND":
			if b, ok := cabrilloBands[strings.ToUpper(val)]; !ok {
				fieldErr = fmt.Errorf("unknown Cabrillo band %q", val)
			} else {
				f = Field{Name: t, Value: b, Type: TypeEnumeration}
			}
		case "MODE":
			f = Field{Name: t, Value: cabrilloModes[strings.ToUpper(val)], Type: TypeEnumeration}
		case "SRX", "STX":
			if _, err := strconv.ParseInt(val, 10, 64); err != nil {
				fieldErr = fmt.Errorf("non-numeric serial number %q: %w", val, err)
			} else {
				f = Field{Name: t, Value: val, Type: TypeNumber}
			}
		case "APP_CABRILLO_TRANSMITTER_ID":
			if len(val) != 1 || !isAllDigits(val) {
				fieldErr = fmt.Errorf("invalid transmitter ID %q", val)
			} else {
				f = Field{Name: t, Value: val, Type: TypeNumber}
			}
		default:
			// TODO check spec type, validate format
			if strings.HasSuffix(val, "-") && val == strings.Repeat("-", len(val)) {
				val = ""
			}
			f = Field{Name: t, Value: val}
		}
		if f.Name != "" {
			break
		}
	}
	if f.Name == "" {
		return f, fieldErr
	}
	return f, nil
}

/*
CabrilloFieldList is a slice of CabrilloFields.  Repeated appearances of a flag
append to the list, or a single flag value can have multiple fields separated by
whitespace.  The syntax is:

  - header:field_a=default (single labeled field, default value)
  - field_a/field_b? (two possible fields, allow empty, no header)
  - header:=default (no lookup, same value for all QSOs)
  - header:field_a/field_b?=default (the works)

Examples: "rst:rst_sent=59" "srx_string/state?" "exch:=CT"
*/
type CabrilloFieldList []CabrilloField

func (l *CabrilloFieldList) String() string {
	s := make([]string, len(*l))
	for i, f := range *l {
		s[i] = f.String()
	}
	return strings.Join(s, " ")
}

func (l *CabrilloFieldList) Set(v string) error {
	sp := strings.Fields(v)
	for _, s := range sp {
		f, err := parseCabrilloField(s)
		if err != nil {
			return err
		}
		*l = append(*l, f)
	}
	return nil
}

func (l *CabrilloFieldList) Get() any { return *l }

type cabrilloConfig struct {
	// core, my, their, extra
	fields                             CabrilloFieldList
	coreLen, myLen, theirLen, extraLen int
	useTabs                            bool
}

func (c cabrilloConfig) toADIF(qso string) (*Record, error) {
	cols := strings.Fields(qso)
	if len(cols) != len(c.fields) {
		return nil, fmt.Errorf("got %d fields, expected %d %s in %q", len(cols), len(c.fields), &c.fields, qso)
	}
	r := NewRecord()
	for i, f := range c.fields {
		v, err := f.toADIF(cols[i])
		if err != nil {
			return nil, err
		}
		r.Set(v)
	}
	if _, ok := r.Get("BAND"); !ok {
		if f, err := r.ParseFloat("FREQ"); err == nil {
			if b, ok := findBandFreq(f); ok {
				r.Set(Field{Name: "BAND", Value: b.adifName, Type: TypeEnumeration})
			}
		}
	}
	return r, nil
}

func (c cabrilloConfig) toCabrillo(r *Record) (cabrilloRecord, error) {
	cr := cabrilloRecord{fields: make([]string, len(c.fields))}
	for i, f := range c.fields {
		v, err := f.fromADIF(r)
		if err != nil {
			return cabrilloRecord{}, err
		}
		cr.fields[i] = v
	}
	if x, err := r.ParseBool("APP_CABRILLO_XQSO"); err == nil {
		cr.xQSO = x
	}
	return cr, nil
}

func (c cabrilloConfig) toLines(l *Logfile) ([][2]string, error) {
	h := make([]string, len(c.fields))
	for i, f := range c.fields {
		h[i] = f.Header
	}
	crs := make([]cabrilloRecord, len(l.Records))
	for i, r := range l.Records {
		cr, err := c.toCabrillo(r)
		if err != nil {
			return nil, err
		}
		crs[i] = cr
	}
	formatLine := func(s []string) string { return strings.Join(s, "\t") }
	sentRcvd := strings.Repeat("\t", c.coreLen-1) + strings.Repeat("\tsent", c.myLen) + strings.Repeat("\trcvd", c.theirLen)
	if !c.useTabs {
		widths := make([]int, len(c.fields))
		for i, s := range h {
			widths[i] = len(s)
		}
		for _, r := range crs {
			for i, s := range r.fields {
				widths[i] = maxInt(widths[i], len(s))
			}
		}
		formatLine = func(s []string) string {
			var r strings.Builder
			for i, f := range s {
				if i != 0 {
					r.WriteRune(' ')
				}
				r.WriteString(f)
				if i < len(widths)-1 {
					for j := widths[i] - len(f); j > 0; j-- {
						r.WriteRune(' ')
					}
				}
			}
			return r.String()
		}
		var corelen, sentlen, rcvdlen int
		for i := 0; i < c.coreLen; i++ {
			corelen += widths[i]
			corelen++ // space between fields
		}
		for i := 0; i < c.myLen; i++ {
			sentlen += widths[c.coreLen+i]
			if i != 0 {
				sentlen++ // space between fields
			}
		}
		for i := 0; i < c.theirLen; i++ {
			rcvdlen += widths[c.coreLen+c.myLen+i]
			if i != 0 {
				rcvdlen++ // space between fields
			}
		}
		isent := "--info sent"
		if sentlen < len(isent) {
			isent = "--sent"
		}
		isent += strings.Repeat("-", maxInt(sentlen-len(isent), 0))
		ircvd := "--info rcvd"
		if rcvdlen < len(ircvd) {
			ircvd = "--rcvd"
		}
		ircvd += strings.Repeat("-", maxInt(rcvdlen-len(ircvd), 0))
		sentRcvd = strings.Repeat(" ", corelen) + isent + " " + ircvd
	}
	lines := make([][2]string, len(l.Records)+2)
	lines[0][0] = "X-Q"
	lines[0][1] = sentRcvd
	lines[1][0] = "X-Q"
	lines[1][1] = formatLine(h)
	for i, cr := range crs {
		if cr.xQSO {
			lines[i+2][0] = "X-QSO"
		} else {
			lines[i+2][0] = "QSO"
		}
		lines[i+2][1] = formatLine(cr.fields)
	}
	return lines, nil
}

type cabrilloRecord struct {
	fields []string
	xQSO   bool
}

var (
	cabFieldFreq      = CabrilloField{TryFields: []string{"FREQ", "BAND"}, Header: "freq"}
	cabFieldMode      = CabrilloField{TryFields: []string{"MODE"}, Header: "mo"}
	cabFieldDate      = CabrilloField{TryFields: []string{"QSO_DATE", "QSO_DATE_OFF"}, Header: "date"}
	cabFieldTime      = CabrilloField{TryFields: []string{"TIME_ON", "TIME_OFF"}, Header: "time"}
	cabFieldMyCall    = CabrilloField{TryFields: []string{"STATION_CALLSIGN", "OPERATOR", "STATION_OWNER"}, Header: "call"}
	cabFieldTheirCall = CabrilloField{TryFields: []string{"CALL"}, Header: "call"}
	cabCoreFields     = []CabrilloField{cabFieldFreq, cabFieldMode, cabFieldDate, cabFieldTime}
)

// Good list of Cabrillo templates: https://www.qrz.lt/ly1vp/ataskaitu_formatai/cabrillo/qso-template.html

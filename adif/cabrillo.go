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
)

// CabrilloIO configures the headers and QSO inference for conversion to and
// from the Cabrillo format.  Most fields configure the value of a header with
// the same name.  Categories maps CATEGORY header names to value, e.g. "TIME"
// to "6-HOURS".  See https://wwrof.org/cabrillo/cabrillo-v3-header/ for
// details about header values.
type CabrilloIO struct {
	Callsign, Contest, Club, CreatedBy, Email,
	GridLocator, Location, Name, Address, Soapbox, MyExchange string
	Operators                           []string
	LowPowerMax, QRPPowerMax            int
	MinReportedOfftime                  time.Duration
	Categories                          map[string]string
	TheirExchangeField, MyExchangeField string
}

func NewCabrilloIO() *CabrilloIO {
	return &CabrilloIO{LowPowerMax: 100, QRPPowerMax: 5, Categories: make(map[string]string)}
}

func (_ *CabrilloIO) String() string { return "cabrillo" }

func (o *CabrilloIO) Read(in io.Reader) (*Logfile, error) {
	headers := make(map[string]string) // this doesn't preserve header order in ADIF
	qsos := make([]cabrilloQSO, 0)
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
		v = strings.TrimPrefix(v, " ")
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
			q, err := parseCabrilloQSO(v)
			if err != nil {
				return nil, err
			}
			q.xQSO = k == "X-QSO"
			qsos = append(qsos, q)
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
	l := NewLogfile()
	for _, q := range qsos {
		r, err := o.toRecord(q)
		if err != nil {
			return nil, err
		}
		for _, f := range fromHeaders {
			r.Set(f)
		}
		l.AddRecord(r)
	}
	for k, v := range headers {
		l.Header.Set(Field{Name: "APP_CABRILLO_" + strings.ReplaceAll(k, "-", "_"), Value: v, Type: TypeString})
	}
	return l, nil
}

func (o *CabrilloIO) Write(l *Logfile, out io.Writer) error {
	qsos := make([]cabrilloQSO, len(l.Records))
	for i, r := range l.Records {
		if q, err := o.toCabrilloQSO(r); err != nil {
			return err
		} else {
			qsos[i] = q
		}
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
			if _, err := w.WriteString(fmt.Sprintf("%s: %s\n", k, line)); err != nil {
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
	setHeader("CLAIMED-SCORE", "")
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
	headerOrder := []string{"SOAPBOX", "CONTEST", "CALLSIGN", "CLUB", "OPERATORS", "NAME", "EMAIL", "ADDRESS", "GRID-LOCATOR", "LOCATION", "CLAIMED-SCORE", "OFFTIME"}
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
	wrote := make(map[string]bool)
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
	// space-align QSO fields
	var widths [11]int
	for i, s := range qsoHeader {
		widths[i] = len(s)
	}
	qstrs := make([][11]string, len(qsos))
	for j, q := range qsos {
		qstrs[j] = [11]string{q.freq, q.mode, q.date, q.time, q.myCall, q.myRST, q.myExch, q.theirCall, q.theirRST, q.theirExch, q.txmitID}
		for i, s := range qstrs[j] {
			widths[i] = maxInt(widths[i], len(s))
		}
	}
	align := func(q [11]string) string {
		var r strings.Builder
		for i, f := range q {
			if i != 0 {
				r.WriteRune(' ')
			}
			r.WriteString(f)
			for j := widths[i] - len(f); j > 0; j-- {
				r.WriteRune(' ')
			}
		}
		return r.String()
	}
	// write QSO field guide and space-aligned QSO records
	isoff := widths[0] + widths[1] + widths[2] + widths[3] + 4
	iroff := widths[4] + widths[5] + widths[6] + 2 - len("--info sent")
	trail := widths[7] + widths[8] + widths[9] + 2 - len("--info rcvd")
	info := strings.Repeat(" ", isoff) + "--info sent" + strings.Repeat("-", iroff) + " --info rcvd" + strings.Repeat("-", trail)
	if err := writeLine("X-Q", info); err != nil {
		return err
	}
	if err := writeLine("X-Q", align(qsoHeader)); err != nil {
		return err
	}
	for i, q := range qsos {
		k := "QSO"
		if q.xQSO {
			k = "X-QSO"
		}
		if err := writeLine(k, align(qstrs[i])); err != nil {
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
		var mode string
		if len(modes) == 1 {
			switch maps.Keys(modes)[0] {
			case "CW":
				mode = "CW"
			case "SSB", "AM", "DIGITALVOICE":
				mode = "SSB"
			case "FM":
				mode = "FM"
			case "RTTY":
				mode = "RTTY"
			default:
				mode = "DIGI"
			}
		} else if len(modes) > 1 {
			mode = "DIGI"
			for m := range modes {
				if m == "CW" || m == "SSB" || m == "FM" || m == "AM" {
					mode = "MIXED"
					break
				}
			}
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

// qsoHeader contains column names to print above space-aligned QSO columns.
// See https://wwrof.org/cabrillo/cabrillo-qso-data/
// Note that the site's example has this row starting with "QSO:" but that
// is incorrect, since these are not part of an actual QSO.
var qsoHeader = [11]string{
	"freq", "mo", "date", "time", "call", "rst", "exch", "call", "rst", "exch", "t"}

type cabrilloQSO struct {
	freq, mode, date, time, myCall, myRST, myExch, theirCall, theirRST, theirExch, txmitID string
	xQSO                                                                                   bool
}

func parseCabrilloQSO(line string) (cabrilloQSO, error) {
	qso := cabrilloQSO{}
	order := []*string{
		&qso.freq, &qso.mode, &qso.date, &qso.time,
		&qso.myCall, &qso.myRST, &qso.myExch,
		&qso.theirCall, &qso.theirRST, &qso.theirExch,
		&qso.txmitID, // transmitter ID is optional
	}
	chunks := strings.Fields(line)
	if len(order) != len(chunks) && len(order)-1 != len(chunks) {
		return qso, fmt.Errorf("want %d fields in Cabrillo QSO, got %d in %q", len(order), len(chunks), line)
	}
	for i, s := range chunks {
		*order[i] = s
	}
	return qso, nil
}

func firstFieldValue(r *Record, fnames ...string) (string, error) {
	for _, n := range fnames {
		if n != "" { // e.g. unset MyExchangeField
			if f, ok := r.Get(n); ok && f.Value != "" {
				return f.Value, nil
			}
		}
	}
	return "", fmt.Errorf("missing %s in %v", strings.Join(fnames, ", "), r)
}

func (o *CabrilloIO) toCabrilloQSO(r *Record) (cabrilloQSO, error) {
	q := cabrilloQSO{}
	if f, ok := r.Get("BAND"); ok && f.Value != "" {
		if q.freq, ok = cabrilloBandsRev[strings.ToLower(f.Value)]; !ok {
			return q, fmt.Errorf("invalid band %q for Cabrillo", f.Value)
		}
	}
	if mhz, err := r.ParseFloat("FREQ"); err == nil && mhz > 0 && mhz < 30 {
		mhzf, _ := r.Get("FREQ")
		q.freq = mhzToKhz(mhzf.Value)
	}
	if q.freq == "" {
		return q, fmt.Errorf("missing FREQ or BAND in %v", r)
	}
	if f, ok := r.Get("MODE"); ok && f.Value != "" {
		switch strings.ToUpper(f.Value) {
		case "CW":
			q.mode = "CW"
		case "RTTY", "RTTYM":
			q.mode = "RY"
		case "SSB", "AM", "DIGITALVOICE":
			q.mode = "PH" // TODO verify DIGITALVOICE counts as phone
		case "FM":
			q.mode = "FM"
		default:
			q.mode = "DG" // most modes are digital
		}
	} else {
		return q, fmt.Errorf("missing MODE in %v", r)
	}
	if d, err := r.ParseDate("QSO_DATE"); err != nil {
		return q, fmt.Errorf("invalid QSO_DATE: %w", err)
	} else {
		q.date = d.Format("2006-01-02")
	}
	var err error
	if q.time, err = firstFieldValue(r, "TIME_ON"); err != nil {
		return q, err
	}
	if q.myCall, err = firstFieldValue(r, "STATION_CALLSIGN", "OPERATOR"); err != nil {
		return q, err
	}
	if q.myRST, err = firstFieldValue(r, "RST_SENT"); err != nil {
		return q, err
	}
	if q.myExch, err = firstFieldValue(r, o.MyExchangeField, "STX_STRING", "STX"); err != nil {
		return q, err
	}
	if q.theirCall, err = firstFieldValue(r, "CALL"); err != nil {
		return q, err
	}
	if q.theirRST, err = firstFieldValue(r, "RST_RCVD"); err != nil {
		return q, err
	}
	if q.theirExch, err = firstFieldValue(r, o.TheirExchangeField, "SRX_STRING", "SRX"); err != nil {
		return q, err
	}
	// TODO is there a beter ADIF field for this?
	q.txmitID = "0"
	if f, ok := r.Get("APP_CABRILLO_TRANSMITTER_ID"); ok && (f.Value == "0" || f.Value == "1") {
		q.txmitID = f.Value
	}
	if x, err := r.ParseBool("APP_CABRILLO_XQSO"); err == nil {
		q.xQSO = x
	}
	return q, nil
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

func (o *CabrilloIO) toRecord(q cabrilloQSO) (*Record, error) {
	r := NewRecord()
	if band := cabrilloBands[strings.ToUpper(q.freq)]; band != "" {
		r.Set(Field{Name: "BAND", Value: band, Type: TypeString})
	} else {
		f, err := strconv.ParseFloat(q.freq, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid frequency %q: %w", q.freq, err)
		}
		// Cabrillo frequencies are in kHz, ADIF are in MHz
		r.Set(Field{Name: "FREQ", Value: strconv.FormatFloat(f/1000.0, 'f', -1, 64), Type: TypeNumber})
		var mb string
		if f > 1800 && f < 2000 {
			mb = "160m"
		} else if f > 3500 && f < 4000 {
			mb = "80m"
		} else if f > 7000 && f < 7300 {
			mb = "40m"
		} else if f > 14000 && f < 14350 {
			mb = "20m"
		} else if f > 21000 && f < 21450 {
			mb = "15m"
		} else if f > 28000 && f < 29700 {
			mb = "10m"
		}
		if mb != "" {
			r.Set(Field{Name: "BAND", Value: mb, Type: TypeString})
		}
	}
	m := cabrilloModes[q.mode]
	if m == "" {
		return nil, fmt.Errorf("invalid Cabrillo mode %q", q.mode)
	}
	r.Set(Field{Name: "MODE", Value: m, Type: TypeString})
	d, err := time.ParseInLocation("2006-01-02", q.date, time.UTC)
	if err != nil {
		return nil, fmt.Errorf("invalid date %q: %w", q.date, err)
	}
	r.Set(Field{Name: "QSO_DATE", Value: d.Format("20060102"), Type: TypeDate})
	t, err := time.ParseInLocation("1504", q.time, time.UTC)
	if err != nil {
		return nil, fmt.Errorf("invalid time %q: %w", q.time, err)
	}
	r.Set(Field{Name: "TIME_ON", Value: t.Format("1504"), Type: TypeTime})
	if !cabrilloCallPat.MatchString(q.myCall) {
		return nil, fmt.Errorf("invalid callsign %q", q.myCall)
	}
	r.Set(Field{Name: "STATION_CALLSIGN", Value: q.myCall, Type: TypeString})
	if !cabrilloRSTPat.MatchString(q.myRST) {
		return nil, fmt.Errorf("invalid RST %q", q.myRST)
	}
	r.Set(Field{Name: "RST_SENT", Value: q.myRST, Type: TypeString})
	if o.MyExchangeField == "" {
		r.Set(Field{Name: "STX_STRING", Value: q.myExch, Type: TypeString})
		if !isAllDigits(q.myExch) {
			r.Set(Field{Name: "STX", Value: q.myExch, Type: TypeNumber})
		}
	} else {
		r.Set(Field{Name: o.MyExchangeField, Value: q.myExch, Type: TypeString})
	}
	if !cabrilloCallPat.MatchString(q.theirCall) {
		return nil, fmt.Errorf("invalid callsign %q", q.theirCall)
	}
	r.Set(Field{Name: "CALL", Value: q.theirCall, Type: TypeString})
	if !cabrilloRSTPat.MatchString(q.theirRST) {
		return nil, fmt.Errorf("invalid RST %q", q.theirRST)
	}
	r.Set(Field{Name: "RST_RCVD", Value: q.theirRST, Type: TypeString})
	if o.TheirExchangeField == "" {
		r.Set(Field{Name: "SRX_STRING", Value: q.theirExch, Type: TypeString})
		if isAllDigits(q.theirExch) {
			r.Set(Field{Name: "SRX", Value: q.theirExch, Type: TypeNumber})
		}
	} else {
		r.Set(Field{Name: o.TheirExchangeField, Value: q.theirExch, Type: TypeString})
	}
	if q.txmitID != "" {
		if q.txmitID != "0" && q.txmitID != "1" {
			return nil, fmt.Errorf("invalid transmitter ID %q", q.txmitID)
		}
		r.Set(Field{Name: "APP_CABRILLO_TRANSMITTER_ID", Value: q.txmitID, Type: TypeNumber})
	}
	if q.xQSO {
		r.Set(Field{Name: "APP_CABRILLO_XQSO", Value: "Y", Type: TypeBoolean})
	}
	return r, nil
}

var (
	cabrilloBands = map[string]string{
		"1800":  "160m",
		"3500":  "80m",
		"7000":  "40m",
		"14000": "20m",
		"21000": "15m",
		"28000": "10m",
		"50":    "6m",
		"70":    "4m",
		"144":   "2m",
		"222":   "1.25m",
		"432":   "70cm",
		"902":   "33cm",
		"1.2G":  "23cm",
		"2.3G":  "13cm",
		"3.4G":  "9cm",
		"5.7G":  "6cm",
		"10G":   "3cm",
		"24G":   "1.25cm",
		"47G":   "6mm",
		"75G":   "4mm",
		"122G":  "2.5mm",
		"134G":  "2mm",
		"241G":  "1mm",
		"LIGHT": "submm",
	}
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
	for k, v := range cabrilloBands {
		cabrilloBandsRev[v] = k
	}
}

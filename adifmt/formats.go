// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/flwyd/adif-multitool/adif"
)

type formatConfig interface {
	Format() adif.Format
	IO() adif.ReadWriter
	AddFlags(fs *flag.FlagSet)
}

var formatConfigs = []formatConfig{
	adiConfig{adif.NewADIIO()},
	adxConfig{adif.NewADXIO()},
	cabrilloConfig{adif.NewCabrilloIO()},
	csvConfig{adif.NewCSVIO()},
	jsonConfig{adif.NewJSONIO()},
	tsvConfig{adif.NewTSVIO()},
}

func formatNamed(n string) formatConfig {
	for _, c := range formatConfigs {
		if strings.EqualFold(string(c.Format()), n) {
			return c
		}
	}
	return nil
}

type adiConfig struct{ io *adif.ADIIO }

func (c adiConfig) Format() adif.Format { return adif.FormatADI }

func (c adiConfig) IO() adif.ReadWriter { return c.io }

func (c adiConfig) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&c.io.ASCIIOnly, "adi-ascii-only", false,
		"ADI files: error on any non-ASCII characters, instead of writing UTF-8")
	fs.BoolVar(&c.io.LowerCase, "adi-lower-case", false,
		"ADI files: print tags in lower case instead of upper case")
	sepHelp := "options: " + strings.Join(adif.SeparatorNames(), ", ")
	fs.Var(&c.io.FieldSep, "adi-field-separator",
		"ADI files: field `separator`\n"+sepHelp)
	fs.Var(&c.io.RecordSep, "adi-record-separator",
		"ADI files: record `separator`\n"+sepHelp)
}

type adxConfig struct{ io *adif.ADXIO }

func (c adxConfig) Format() adif.Format { return adif.FormatADX }

func (c adxConfig) IO() adif.ReadWriter { return c.io }

func (c adxConfig) AddFlags(fs *flag.FlagSet) {
	fs.IntVar(&c.io.Indent, "adx-indent", 1, "ADX files: indent nested XML structures `n` spaces, 0 for no whitespace")
}

type cabrilloConfig struct{ io *adif.CabrilloIO }

func (c cabrilloConfig) Format() adif.Format { return adif.FormatCabrillo }

func (c cabrilloConfig) IO() adif.ReadWriter { return c.io }

func (c cabrilloConfig) AddFlags(fs *flag.FlagSet) {
	c.io.CreatedBy = "ADIF Multitool " + version
	fs.IntVar(&c.io.LowPowerMax, "cabrillo-max-power-low", c.io.LowPowerMax, "Higest allowed power in `watts` considered LOW power by the contest")
	fs.IntVar(&c.io.QRPPowerMax, "cabrillo-max-power-qrp", c.io.QRPPowerMax, "Higest alqrped power in `watts` considered QRP power by the contest")
	fs.StringVar(&c.io.Callsign, "cabrillo-callsign", "", "Cabrillo files: CALLSIGN header `value`")
	fs.StringVar(&c.io.Club, "cabrillo-club", "", "Cabrillo files: CLUB header `value`")
	// TODO Operators (string slice)
	fs.StringVar(&c.io.Contest, "cabrillo-contest", "", "Cabrillo files: CONTEST header `value`")
	fs.StringVar(&c.io.Email, "cabrillo-email", "", "Cabrillo files: EMAIL address header `value`")
	fs.StringVar(&c.io.GridLocator, "cabrillo-grid-locator", "", "Cabrillo files: GRID-LOCATOR header `value`")
	fs.StringVar(&c.io.Location, "cabrillo-location", "", "Cabrillo files: LOCATION header `value` (e.g. ARRL section)")
	fs.StringVar(&c.io.Name, "cabrillo-name", "", "Cabrillo files: NAME header `value` (your name or club name)")
	fs.StringVar(&c.io.Address, "cabrillo-address", "", "Cabrillo files: ADDRESS header `value` (include newlines)")
	fs.StringVar(&c.io.Soapbox, "cabrillo-soapbox", "", "Cabrillo files: SOAPBOX header `value` (free-form comment)")
	// TODO MinReportedOfftime (duration)
	fs.StringVar(&c.io.MyExchange, "cabrillo-my-exchange", "", "Cabrillo files: `value` sent as exchange in QSOs")
	fs.StringVar(&c.io.MyExchangeField, "cabrillo-my-exchange-field", "", "Cabrillo files: ADIF `field` used for contest exchange sent")
	fs.StringVar(&c.io.TheirExchangeField, "cabrillo-their-exchange-field", "", "Cabrillo files: ADIF `field` used for contest exchange sent")
	fs.StringVar(&c.io.TheirExchangeAlt, "cabrillo-their-exchange-field-alt", "", "Cabrillo files: ADIF `field` used as exchange if --cabrillo-their-exchange-field, SRX_STRING, and SRX are not set")
	for v, a := range adif.CabrilloCategoryValues {
		fs.Var(&mapValue{c.io.Categories, v, a},
			"cabrillo-category-"+strings.ToLower(v),
			fmt.Sprintf("Cabrillo files: CATEGORY-%s header `value` (%s)", v, strings.Join(a, ", ")))
	}
}

type csvConfig struct{ io *adif.CSVIO }

func (c csvConfig) Format() adif.Format { return adif.FormatCSV }

func (c csvConfig) IO() adif.ReadWriter { return c.io }

func (c csvConfig) AddFlags(fs *flag.FlagSet) {
	// TODO csv-lower-case
	// TODO separate comma values for input and output?
	fs.Var(&runeValue{&c.io.Comma}, "csv-field-separator", "CSV files: field separator `character` if not comma")
	fs.Var(&runeValue{&c.io.Comment}, "csv-comment", "CSV files: ignore lines beginnig with `character`")
	fs.BoolVar(&c.io.LazyQuotes, "csv-lazy-quotes", false, "CSV files: be relaxed about quoting rules")
	fs.BoolVar(&c.io.RequireFullRecord, "csv-require-all-fields", false, "CSV files: error if fewer fields in a record than in header")
	fs.BoolVar(&c.io.TrimLeadingSpace, "csv-trim-space", false, "CSV files: ignore leading space in fields")
	fs.BoolVar(&c.io.CRLF, "csv-crlf", false, "CSV files: output MS Windows line endings")
}

type jsonConfig struct{ io *adif.JSONIO }

func (c jsonConfig) Format() adif.Format { return adif.FormatJSON }

func (c jsonConfig) IO() adif.ReadWriter { return c.io }

func (c jsonConfig) AddFlags(fs *flag.FlagSet) {
	// TODO json-lower-case
	fs.BoolVar(&c.io.HTMLSafe, "json-html-safe", false, "JSON files: escape characters including < > & for use in HTML")
	fs.IntVar(&c.io.Indent, "json-indent", 1, "JSON files: indent nested JSON structures `n` spaces, 0 for no whitespace")
	fs.BoolVar(&c.io.TypedOutput, "json-typed-output", false, "JSON files: output numbers and booleans instead of strings")
}

type tsvConfig struct{ io *adif.TSVIO }

func (c tsvConfig) Format() adif.Format { return adif.FormatTSV }

func (c tsvConfig) IO() adif.ReadWriter { return c.io }

func (c tsvConfig) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&c.io.CRLF, "tsv-crlf", false, "TSV files: output MS Windows line endings")
	fs.BoolVar(&c.io.EscapeSpecial, "tsv-escape-special", false, "TSV files: accept and produce \\t \\r \\n and \\\\ escapes in fields")
	fs.BoolVar(&c.io.IgnoreEmptyHeaders, "tsv-ignore-empty-headers", false, "TSV files: do not return error if a TSV file has an empty header field")
}

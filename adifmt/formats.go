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
	"github.com/flwyd/adif-multitool/adif/spec"
)

type formatConfig interface {
	Format() adif.Format
	IO() adif.ReadWriter
	AddFlags(fs *flag.FlagSet)
	Help() string
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
	fs.BoolVar(&c.io.AllowUnknownTag, "adi-allow-unknown-tags", false,
		"Convert <tag> to [tag] in comments instead of error")
	fs.BoolVar(&c.io.ASCIIOnly, "adi-ascii-only", false,
		"Error on any non-ASCII characters, instead of writing UTF-8")
	fs.BoolVar(&c.io.LowerCase, "adi-lower-case", false,
		"Print tags in lower case instead of upper case")
	sepHelp := fmt.Sprintf("(%s)", strings.Join(adif.SeparatorNames(), ", "))
	fs.Var(&c.io.FieldSep, "adi-field-separator",
		"Field `separator` after values\n"+sepHelp)
	fs.Var(&c.io.RecordSep, "adi-record-separator",
		"Record `separator` after <EOR>\n"+sepHelp)
}

func (c adiConfig) Help() string {
	return `ADI is an ASCII text file format defined for amateur radio data interchange.
The specification is at ` + spec.ADIFSpecURL + `#ADI_File_Format
This program accepts and outputs Unicode data in UTF-8 encoding unless the
--adi-ascii-only option is given.  Fields are not required to be part of the
ADIF specification.
`
}

type adxConfig struct{ io *adif.ADXIO }

func (c adxConfig) Format() adif.Format { return adif.FormatADX }

func (c adxConfig) IO() adif.ReadWriter { return c.io }

func (c adxConfig) AddFlags(fs *flag.FlagSet) {
	fs.IntVar(&c.io.Indent, "adx-indent", 1,
		"Indent nested XML structures `n` spaces, 0 for no whitespace")
}

func (c adxConfig) Help() string {
	return `ADX is an XML file format defined for amateur radio data interchange.
The specification is at ` + spec.ADIFSpecURL + `#ADX_File_Format
Fields are not required to be part of the ADIF specification.
`
}

type cabrilloConfig struct{ io *adif.CabrilloIO }

func (c cabrilloConfig) Format() adif.Format { return adif.FormatCabrillo }

func (c cabrilloConfig) IO() adif.ReadWriter { return c.io }

func (c cabrilloConfig) AddFlags(fs *flag.FlagSet) {
	c.io.CreatedBy = "ADIF Multitool " + version
	fs.BoolVar(&c.io.TabDelimiter, "cabrillo-delimiter-tab", false,
		"Use tabs rather than space-aligned columns")
	fs.IntVar(&c.io.LowPowerMax, "cabrillo-max-power-low", c.io.LowPowerMax,
		"Highest allowed power in `watts` considered LOW power by the contest")
	fs.IntVar(&c.io.QRPPowerMax, "cabrillo-max-power-qrp", c.io.QRPPowerMax,
		"Highest allowed power in `watts` considered QRP power by the contest")
	fs.StringVar(&c.io.Callsign, "cabrillo-callsign", "",
		"CALLSIGN Cabrillo header `value`")
	fs.IntVar(&c.io.ClaimedScore, "cabrillo-claimed-score", 0,
		"CLAIMED-SCORE Cabrillo header `value`")
	fs.StringVar(&c.io.Club, "cabrillo-club", "",
		"CLUB Cabrillo header `value`")
	// TODO Operators (string slice)
	fs.StringVar(&c.io.Contest, "cabrillo-contest", "",
		"CONTEST Cabrillo header `value`")
	fs.StringVar(&c.io.Email, "cabrillo-email", "",
		"EMAIL address Cabrillo header `value`")
	fs.StringVar(&c.io.GridLocator, "cabrillo-grid-locator", "",
		"GRID-LOCATOR Cabrillo header `value`")
	fs.StringVar(&c.io.Location, "cabrillo-location", "",
		"LOCATION Cabrillo header `value`, e.g. ARRL section")
	fs.StringVar(&c.io.Name, "cabrillo-name", "",
		"NAME Cabrillo header `value` (your name or club name)")
	fs.StringVar(&c.io.Address, "cabrillo-address", "",
		"ADDRESS Cabrillo header `value` (include newlines)")
	fs.StringVar(&c.io.Soapbox, "cabrillo-soapbox", "",
		"SOAPBOX Cabrillo header `value` (free-form comment)")
	// TODO MinReportedOfftime (duration)
	fs.Var(&c.io.MyExchange, "cabrillo-my-exchange",
		"Field `configuration` ("+adif.CabrilloFieldExample+") of\nlogging station's exchange (repeatable)")
	fs.Var(&c.io.TheirExchange, "cabrillo-their-exchange",
		"Field `configuration` ("+adif.CabrilloFieldExample+") of\ncontacted station's exchange (repeatable)")
	fs.Var(&c.io.ExtraFields, "cabrillo-extra-field",
		"Name of `field` added at the end of QSO lines (repeatable)\ne.g. APP_CABRILLO_TRANSMITTER_ID")
	for v, a := range adif.CabrilloCategoryValues {
		fs.Var(&mapValue{c.io.Categories, v, a},
			"cabrillo-category-"+strings.ToLower(v),
			fmt.Sprintf("CATEGORY-%s Cabrillo header `value`\n(%s)", v, strings.Join(a, ", ")))
	}
}

func (c cabrilloConfig) Help() string {
	return `Cabrillo is a text file format designed for amateur radio contest logs.
The format is described at https://wwrof.org/cabrillo/cabrillo-v3-header/ and
https://wwrof.org/cabrillo/cabrillo-qso-data/
The ADIF fields used in the contest exchange must be given in the
--cabrillo-my-exchange and --cabrillo-their-exchange options.  The format of
these values is '` + adif.CabrilloFieldExample + `'
Fields are printed in QSO: lines separated by space.  The field header is
printed in a comment above the field column.  If multiple ADIF fields are
separated by / the Cabrillo QSO will include the first non-blank field value.
At least one field must have a value unless the ? suffix is given or a default
value is specified after an = character.  For contest exchange examples, see
` + helpUrl + `#cabrillo
`
}

type csvConfig struct{ io *adif.CSVIO }

func (c csvConfig) Format() adif.Format { return adif.FormatCSV }

func (c csvConfig) IO() adif.ReadWriter { return c.io }

func (c csvConfig) AddFlags(fs *flag.FlagSet) {
	// TODO csv-lower-case
	// TODO separate comma values for input and output?
	fs.Var(&runeValue{&c.io.Comma}, "csv-field-separator",
		"Field separator `character` if not comma")
	fs.Var(&runeValue{&c.io.Comment}, "csv-comment",
		"Ignore lines beginning with `character`")
	fs.BoolVar(&c.io.LazyQuotes, "csv-lazy-quotes", false,
		"Be relaxed about quoting rules")
	fs.BoolVar(&c.io.RequireFullRecord, "csv-require-all-fields", false,
		"Error if fewer fields in a record than in header")
	fs.BoolVar(&c.io.TrimLeadingSpace, "csv-trim-space", false,
		"Ignore leading space in fields")
	fs.BoolVar(&c.io.CRLF, "csv-crlf", false,
		"Output MS Windows line endings")
	fs.BoolVar(&c.io.OmitHeader, "csv-omit-header", false,
		"Don't output the header line")
}

func (c csvConfig) Help() string {
	return `CSV (comma-separated values) is a widely-used format for sharing tabular data
defined at https://datatracker.ietf.org/doc/html/rfc4180
ADIF field names (case-insensitive) must appear in the first line.  Values
must be surrounded by double quotes (") if they contain commas, line breaks,
or double quotes, which are repeated ("") as an escape.  Input is not required
to have the same number of values in each line; output will have one value for
each ADIF field which appears in the log, even if a record does not have a
value for that field.  CSV is a convenient intermediate format for data logged
in a spreadsheet program or transcribed from paper logs.
`
}

type jsonConfig struct{ io *adif.JSONIO }

func (c jsonConfig) Format() adif.Format { return adif.FormatJSON }

func (c jsonConfig) IO() adif.ReadWriter { return c.io }

func (c jsonConfig) AddFlags(fs *flag.FlagSet) {
	// TODO json-lower-case
	fs.BoolVar(&c.io.HTMLSafe, "json-html-safe", false,
		"Escape characters including < > & for use in HTML")
	fs.IntVar(&c.io.Indent, "json-indent", 1,
		"Indent nested JSON structures `n` spaces, 0 for no whitespace")
	fs.BoolVar(&c.io.TypedOutput, "json-typed-output", false,
		"Output numbers and booleans instead of strings")
}

func (c jsonConfig) Help() string {
	return `JSON is a popular format for exchanging arbitrary structured data between
computer programs.  The ADIF specification does not define a JSON mapping, so
this program uses a simple structure: a log is an object with a "HEADER"
object property and a "RECORDS" array property.  All-caps ADIF fields are
properties of the header and record objects and all values are strings unless
--json-typed-output is set.
`
}

type tsvConfig struct{ io *adif.TSVIO }

func (c tsvConfig) Format() adif.Format { return adif.FormatTSV }

func (c tsvConfig) IO() adif.ReadWriter { return c.io }

func (c tsvConfig) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&c.io.CRLF, "tsv-crlf", false,
		"Output MS Windows line endings")
	fs.BoolVar(&c.io.EscapeSpecial, "tsv-escape-special", false,
		`Accept and produce \t \r \n and \\ escapes in fields`)
	fs.BoolVar(&c.io.IgnoreEmptyHeaders, "tsv-ignore-empty-headers", false,
		"Don't return error if a TSV file has an empty header field")
	fs.BoolVar(&c.io.OmitHeader, "tsv-omit-header", false,
		"Don't output the header line")
}

func (c tsvConfig) Help() string {
	return `TSV (tab-separated values) is a widely-used format for sharing tabular data,
see https://www.iana.org/assignments/media-types/text/tab-separated-values
ADIF field names (case-insensitive) must appear in the first line.
TSV cannot handle multi-line values like mailing addresses unless the
--tsv-escape-special option is given and the receiving application supports
escape sequences.  Input is not required to have the same number of values in
each line; output will have one value for each ADIF field which appears in the
log, even if a record does not have a value for that field.  TSV is a
convenient intermediate format for data logged in a spreadsheet program or
transcribed from paper logs.
`
}

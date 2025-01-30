# ADIF Multitool

Validate, modify, and convert ham radio log files with a handy command-line
tool. 📻🌳🪓

`adifmt` provides a suite of commands for working with
[ADIF](https://adif.org/) logs from ham radio software.  It is run from a
shell, via [Terminal](https://en.wikipedia.org/wiki/Terminal_(macOS)) on macOS
and [PowerShell](https://en.wikipedia.org/wiki/PowerShell),
[cmd.exe](https://en.wikipedia.org/wiki/Cmd.exe), or
[Windows Terminal](https://en.wikipedia.org/wiki/Windows_Terminal) on Windows.
Each `adifmt` invocation reads log files from the command line or standard
input and prints an ADIF log to standard output, allowing multiple commands to
be chained together in a pipeline.  For example, to add a `BAND` field based on
the `FREQ` (radio frequency) field, add your station's maidenhead locator
(`MY_GRIDSQURE`) to all entries, automatically fix some incorrectly formatted
fields, validate that all fields are properly formatted, and save a log file
containing only SSB voice contacts, a pipeline might look like

```sh
adifmt infer --fields band my_original_log.adi \
  | adifmt edit --add my_gridsquare=FN31pr \
  | adifmt fix \
  | adifmt validate \
  | adifmt find --if mode=SSB \
  | adifmt save my_ssb_log.adx
```

On Windows, PowerShell uses the backtick character (`` ` ``) and Command Prompt
uses caret (`^`)  instead of backslash (`\`) for multi-line pipelines.  You
can also put the whole pipeline on a single line; they are presented as
multiple lines here for readability.

*Note*: `adifmt` is pronounced “ADIF M T” or “ADIF multitool”, not “adi fmt” nor
”addy format”.

## Quick start

Binaries for each ADIF Multitool version are available on the
[releases page](https://github.com/flwyd/adif-multitool/releases).  You can also
build it from source code with a [Go compiler](https://go.dev/dl/).  Run
`go install github.com/flwyd/adif-multitool/adifmt@latest` to make the `adifmt`
command available.  (You may need to add the `$GOBIN` environment variable to
your path.)  To see if it works, run `adifmt help`.  If the command is not
found, try `go run github.com/flwyd/adif-multitool/adifmt help`

To do something useful with ADIF Multitool, the syntax is

```
adifmt command [files...] [options] [files...]
```

For example, the `cat` command concatenates all input files and outputs ADIF
data to standard output:

```sh
adifmt cat log1.adi log2.adi > combined.adi
```

prints all of the records in the two `logX.adi` files to the `combined.adi`
file.

Options (also known as flags) control input and output options.  For example,
to print records with a UNIX newline between fields, two newlines between
records, use lower case for all field names, and add user defined fields
`gain_db` (range ±100) and `radio_color` (values black, white, or gray):

```sh
adifmt cat --adi-field-separator=newline \
  --adi-record-separator=2newline \
  --adi-lower-case \
  --userdef='GAIN_DB,{-100:100}' \
  --userdef='radio_color,{black,white,gray}' \
  log1.csv
```

File names to process can come before or after the options list, but cannot be
intermixed.  A `--` by itself ends option parsing, which allows working with a
log file whose name starts with `-` (hyphen/minus/dash).

Multiple input and output formats are supported (currently ADI and ADX per the
ADIF spec, Cabrillo according to the WWROF spec, CSV and TSV with field names
matching the ADIF list, and JSON with a similar format to ADX).

```sh
adifmt cat --input=adi --output=csv log1.adi > log1.csv
adifmt cat --input=csv --output=adi log2.csv > log2.adi
```

`--input` need not be specified if it’s implied by the file name or can be
inferred from the structure of the data.  `--ouput=adi` is the default for
output format.  `adifmt save` infers the output format from the file’s
extension.  Input files can be in different formats:

```sh
adifmt cat log1.adi log2.adx log3.csv log4.json log5.tsv log6.cbr > combined.adi
```

If no file names are given, input is read from standard input:

```
gunzip --stdout mylog.csv.gz | adifmt cat --output=adx | gzip > mylog.adx.gz
```

This is useful in composing several `adifmt` invocations together.  Commands
can be combined in a Unix-style pipeline.  The `fix` command automatically
changes some values to match the expected ADIF format such as changing a time
field from `12:34:56` to `123456` and a date from `2012-03-04` to `20120304`.
The `select` command prints only a subset of fields.  The `save` command writes
the input data to a file.  These can be combined:

```sh
adifmt fix log1.adi \
  | adifmt select --fields qso_date,time_on,call \
  | adifmt save minimal.csv
```

creates a file named `minimal.csv` with just the date, time, and callsign from
each record in the input file `log1.adi`.

## Practical examples

The following examples transform data into a format expected by a particular
program or ham radio activity.  For details on how they work, see documentation
for individual commands or formats below.  For contest log examples, see the
[Cabrillo](#cabrillo) section.  Contributions of useful pipelines are welcome.

### Add station and location to a log

This example uses `edit` to add several fields to a POTA log saved in CSV
format.  It then uses `fix` to remove the `:` from the time, transform the
decimal (GPS) latitude and longitude to ADIF sexagesimal format and transforms
`USA` and `CAN` country abbreviations to DXCC entity numbers.  `flatten` makes
two copies of each record, one for park `US-0791` and one for park `US-4567`.
`infer` then sets the band from the frequency, grid square (Maidenhead locator)
based on the latitude and longitude, `STATION_CALLSIGN` field to the `OPERATOR`
field, and `SIG` and `SIG_INFO` from the `POTA_REF` field.  (POTA doesn't
require the country or latitude/longitude fields; they're included for
illustration.) The input log might look like this:

```csv
TIME_ON,FREQ,MODE,CALL,STATE,COUNTRY
12:34,7.012,CW,W1AW,CT,USA
12:56,14.234,SSB,VA1XYZ,NS,CAN
```

```sh
adiifmt edit mylog.csv \
  --add qso_date=20240704 \
  --add operator=WT0RJ \
  --add my_pota_ref=US-0791,US-4567 \
  --add my_state=DC --add my_country=USA \
  --add my_lat=38.899736 --add my_lon=-77.063331 \
| adifmt fix \
| adifmt flatten --fields pota_ref,my_pota_ref \
| adifmt infer --fields band,my_gridsquare,station_callsign
```

### ADIF to SOTA CSV

This example uses `find` to filter out any records which don’t have `SOTA_REF`
or `MY_SOTA_REF` fields, `edit` to add a `V2` field to each record (required by
the SOTA uploader), `select` to output only the fields expected by the SOTA
uploader and in the right order, `validate` to ensure fields are present and
correctly formatted, and `save --csv-omit-header` to create a file with just
the records, no file header.  If your log lacks frequencies, replace the `freq`
field with `band`.  (Note that the SOTA uploader now accepts ADIF files, so you
could just use the `find` command and upload directly.  This example may be
useful if the data need to be further transformed or imported by a SOTA data
analysis program.)

```sh
SOTA_CSV_ORDER=version,station_callsign,my_sota_ref,qso_date,time_on,freq,mode,call,sota_ref,comment
adifmt find mylog.adi --if-not 'my_sota_ref=' --or-if-not 'sota_ref=' \
| adifmt edit --set version=V2 \
| adifmt select --fields $SOTA_CSV_ORDER \
| adifmt validate --required-fields station_callsign,qso_date,time_on,freq,mode,call \
| adifmt save --csv-omit-header --field-order $SOTA_CSV_ORDER sotalog.csv
```

The variable assignment syntax for `SOTA_CSV_ORDER` works on Mac and Linux.  On
Windows PowerShell assign the variable as `$SOTA_CSV_ORDER =
version,station_callsign,...`.  On Windows cmd.exe, assign it as
`set SOTA_CSV_ORDER = version,station_callsign,...` and reference it as
`%SOTA_CSV_ORDER%` rather than the `$` prefix.

### Filter WARC bands

This example uses `infer` to set the band from the frequency if the former
isn’t set.  It then  uses `find` to filter out any contacts made on the
[WARC bands](https://en.wikipedia.org/wiki/WARC_bands): 12, 17, 30, and 60
meters.  Contesting is not allowed on those bands, so this is useful when
preparing a contest submission from a general station log where contacts may
have been made on bands not part of the contest.

```sh
adifmt infer --fields band mylog.adi \
| adifmt find --if-not 'band=60m|30m|17m|12m'
```

### Set mode from frequency (U.S. band plan)

This example sets the mode and submode based on the frequency, according to the
U.S. band plan.  It assumes that CW and SSB are the only modes in use (no FM,
AM, or digital), but can be extended if there are frequency ranges that you use
exclusively for one mode.  `edit --add` will not overwrite the mode if it
already has a value (`edit --set` would force the new value).

```sh
# US HF SSB (overlaps SSTV & AM) 80m: 3.6:4, 40m: 7.125:7.3, 20m:14.15:14.35,
# 17m:18.11:18.168, 15m:21.2:21.45, 12m:24.93 to 24.99, 10m: 28.3:29
# 6m 50.1:50.3 is CW/SSB, SSB calling 60.125, assume 60.120+ is SSB
adifmt edit --if 'freq>3.6' --if 'freq<4' \
  --or-if 'freq>7.125' --if 'freq<=7.3' \
  --or-if 'freq>=14.15' --if 'freq<14.35' \
  --or-if 'freq>=18.11' --if 'freq<18.168' \
  --or-if 'freq>=21.2' --if 'freq<21.45' \
  --or-if 'freq>=24.93' --if 'freq<24.99' \
  --or-if 'freq>=28.3' --if 'freq<29' \
  --or-if 'freq>=50.12' --if 'freq<50.3' \
  --add mode=SSB |\

# US HF CW (some digital could occur) low end of the band, below FT8 and friends
adifmt edit --if 'freq>3.5' --if 'freq<3.7' \
  --or-if 'freq>7' --if 'freq<7.07' \
  --or-if 'freq>10.1' --if 'freq<10.3' \
  --or-if 'freq>14' --if 'freq<14.07' \
  --or-if 'freq>18.068' --if 'freq<18.1' \
  --or-if 'freq>21' --if 'freq<21.07' \
  --or-if 'freq>24.89' --if 'freq<14.91' \
  --or-if 'freq>28' --if 'freq<28.07' \
  --or-if 'freq>50.1' --if 'freq<50.12' \
  --add mode=CW |\

# SSB is usually LSB on 40m and below except 60m, USB on 20m and above
adifmt edit --if mode=SSB --if 'freq<8' --if-not band=60m --add submode=LSB |\
adifmt edit --if mode=SSB --if 'freq>=14' --or-if band=60m --add submode=USB
```

## Features

### Input/Output formats

`adifmt` can read from and write to the following formats.  ADI (tag-based) and
ADX (XML-based) formats are [specified by ADIF](https://adif.org.uk/adiif).
The Cabrillo V3 contest log format is
[specified by WWROF](https://wwrof.org/cabrillo/).
Others use standard formats for arbitrary key-value data.  Format-specific
options are configured with option flags.  Formats are inferred from file names
or can be set explicitly via `--input` and `--output` options.

Name     | Extension                   | Notes
-------- | --------------------------- | -----
ADI      | `.adi`                      | Outputs `IntlString` (Unicode fields) in UTF-8
ADX      | `.adx`                      |
Cabrillo | `.cbr`, `.log`, `.cabrillo` | See [Cabrillo](#cabrillo) section
CSV      | `.csv`                      | Comma-separated values; other delimiters supported via the `--csv-field-separator` option
JSON     | `.json`                     | Can parse number and boolean typed data, to write these set the `--json-typed-output` option
TSV      | `.tsv`                      | Tab-separated values, tabs and line breaks escaped if `--tsv-escape-special` is set

Input files can have fields with any names, even if they’re not part of the
ADIF spec.  The `--userdef` option will add user-defined field metadata to ADI
and ADX output specifying type, range, or valid enumeration values.  ADX XML
tags must be upper case; other formats accept any case field names in input
files and use `UPPER_SNAKE_CASE` for output by default.  Application-defined
fields in CSV, TSV, and JSON should use the `APP_PROGRAMNAME_FIELD_NAME` syntax
used in ADI files. JSON input files should be structured as follows; `HEADER` is
optional.

```json
{
 "HEADER": {
  "ADIF_VER": "3.1.5",
  "more": "header fields"
 },
 "RECORDS": [
  {
   "CALL": "W1AW",
   "more_fields": "record fields"
  },
  {
   "CALL": "NA1SS",
   "more_fields": "additional record fields"
  }
 ]
}
```

Some (but not all) comments found in ADI and ADX files are preserved from input
to output.  Details of comment handling are subject to change and should not be
depended upon.

#### Cabrillo

**Note: app-specific fields for Cabrillo are currently experimental and may be
replaced by official ADIF fields in a future version, pending proposals to
update the ADIF specification.**

The [Cabrillo](https://wwrof.org/cabrillo/) format is commonly used to submit
logs for ham radio contests.  ADIF Multitool can convert to and from Carbillo,
but this is a lossy process: many ADIF fields are not included in Cabrillo and
some Cabrillo values don't perfectly map to ADIF like the `DIGI` mode and the
transmitter ID field.  The latter is imported as an app-specific field,
`APP_CABRILLO_TRANSMITTER_ID`.

Cabrillo contacts starting with `X-QSO:` rather than `QSO:` are imported with an
`APP_CABRILLO_XQSO` boolean field set; if this field is set and true (`Y`) then
`X-QSO:` will be used for export.  These contacts are used by contest organizers
to confirm contacts without granting credit, e.g. if they were made with too
high a power for the submitting station’s category.

When converting to or from Cabrillo, the relevant fields for the contest must be
provided as options.  The “core” fields—frequency/band, mode, date, time,
your station’s callsign, and the contacted station’s callsign—are always
included.  You need to provide the field(s) or values for your exchange as
`--cabrillo-my-exchange` and the contacted station’s exchange fields as
`--cabrillo-their-exchange`.  If additional fields are needed at the end of the
QSO (typically transmitter ID) use `--cabrillo-extra-field`.  Each of these
options can be given multiple times for multiple fields, or given once with a
quoted space-delimited list of fields.  The syntax for each field is
`header:field_a/field_b?=default`, with each portion optional, but either at
least one field or a default value must be specified.  See the next section for
full examples.

* `header:` provides a short name for the field which appears in a comment
  (`X-Q:` line) above the QSO list to make it easy to see what the fields mean.
  Examples: `rst:`, `exch:`, `name:`.
* `field1` or `field1/field2` specifies which ADIF field should provide the
  value for the Cabrillo field.  If there is no value in `field1`, `field2` is
  used.  There is no limit on the number of fields.  Field names are
  case-insensitive.  Examples: `RST_SENT`, `STX/STX_STRING`,
  `srx_string/arrl_section/state`.
* `?` makes a field optional.  If there is no data in any of the fields and no
  default set, one or more `-` characters will be used in the Cabrillo file to
  mark the absence of a value, since Cabrillo is space-separated.  If `?` is not
  present in a field definition, conversion to Cabrillo will fail if no ADIF
  field has a value for a record.  Note that only one `?` at the end of the
  field list is required.  Examples: `IOTA?`, `srx/srx_string?`.
* `=default` provides a default value if no ADIF field is set.  A default value
  may be provided without any fields if an exchange is the same for the whole
  contest, e.g. your name or location, though this will result in the field
  being discarded during import from Cabrillo.  Examples: `=Hiram` (Contest
  exchange is name), `MY_GRIDSQUARE=FN31PR` (use grid locator if present, use
  `FN31PR` if not set), `arrl_section/state=DX` (use `ARRL_SECTION` or `STATE`
  field if set, otherwise log exchange as `DX`).

When converting from Cabrillo, header fields like `CLUB` and `CATEGORY-OVERLAY`
are preserved as ADIF headers with `APP_CABRILLO_` prefixes, e.g.
`APP_CABRILLO_CLUB` and `APP_CABRILLO_CATEGORY_OVERLAY` (hyphens are replaced
by underscores).  (ADIF does not technically support app-defined fields in the
header.  The `--suppress-app-headers` flag will disable this output.)  When
converting from ADIF to Cabrillo, header fields can be set by the same app
headers or command-line flags like `adifmt cat --output=cabrillo
--cabrillo-club="Springfield ARC" --cabrillo-category-overlay=YOUTH log.adi`.
ADIF Multitool will infer `CONTEST`, `CALLSIGN`, `OPERATORS`, `GRID-LOCATOR`,
`LOCATION`, `CATEGORY-BAND`, `CATEGORY-MODE`, and `CATEGORY-POWER` headers from
values in the log's records, but make sure to double-check the output.  Power
levels for LOW and QRP are set with `--cabrillo-max-power-low` and
`--cabrillo-max-power-qrp`.  Other headers are included in the output file with
no value; fill these lines in based on contest instructions or delete them if
not needed by the contest sponsor.  ADIF Multitool does not attempt to
calculate scores for any contests.

Since the mapping between ADIF and Cabrillo is not a perfect match, double-check
your log file carefully and
[report any conversion bugs](https://github.com/flwyd/adif-multitool/issues).
Cabrillo 3.0 is currently the only supported format for import or export;
Cabrillo 2.0 support could be added if there is demand.

##### Common Cabrillo exchanges

Rather than exhaustively list every contest and a Cabrillo template, the
following are options for exchanges formats which are in common use.  Find the
exchange used by your contest.  Contest exchanges can typically be found on
[WA7BNM Contest Calendar](https://contestcalendar.com/); [LY1VP maintains many
Cabrillo file examples](https://www.qrz.lt/ly1vp/ataskaitu_formatai/cabrillo/qso-template.html).

The headers (text before `:`) in field definitions are printed in an `X-Q:`
comment for easy visual field alignment.  Example output for an ARRL Field Day
log might look like

```
X-Q:                          --info sent--- --info rcvd----
X-Q: freq  mo date       time call class sec call  class sec
QSO: 21345 PH 2024-06-22 2345 VE3Z 1B    ONN W1AW  10A   CT
QSO: 7015  CW 2024-06-23 0015 VE3Z 1B    ONN K1USN 2A    MA
```

###### WWROF example with transmitter ID

This set of options outputs the
[example in the Cabrillo specification](https://wwrof.org/cabrillo/cabrillo-qso-data/)
which has the logging station sending a numeric value, the contacted station
sending a two-letter state, and a transmitter ID (`0` or `1`) as an extra field.
RST signal reports are also in the exchange.

`--cabrillo-my-exchange 'rst:RST_SENT exch:STX/STX_STRING'
--cabrillo-their-exchange 'rst:RST_RCVD exch:STATE'
--cabrillo-extra-field 't:APP_CABRILLO_TRANSMITTER_ID'`

###### Signal report and location

This example prefers `SRX_STRING`, e.g. county abbreviations in a state QSO
party, then `STATE`.  A different location field like `ARRL_SECTION` could be
used.  Use `STX_STRING/MY_STATE=NY` to provide `NY` as your default exchange.

`--cabrillo-my-exchange 'rst:RST_SENT exch:STX_STRING/MY_STATE'
--cabrillo-their-exchange 'rst:RST_RCVD exch:SRX_STRING/STATE'`

###### Name and location

Operator name (rather than signal report) is used by North America QSO Party,
some other state QSO parties, and a variety of international contests.  ADIF
Multitool will insert an underscore in any multi-word exchange value like
`PEGGY_SUE`.  Make sure to use your own name as default.

`--cabrillo-my-exchange 'name:MY_NAME=Napoleon exch:STX_STRING/MY_STATE'
--cabrillo-their-exchange 'name:NAME exch:SRX_STRING/STATE'`

###### ARRL Field Day

ARRL Field Day, Winter Field Day, and some other contests use a station
classifier and ARRL section.  There is no `MY_CLASS` field in ADIF, so it’s
provided with just a default (and therefore would not be imported to ADIF).
Change `3A` and `PAC` to your station’s class and section.

`--cabrillo-my-exchange 'class:=3A sec:MY_ARRL_SECTION=PAC'
--cabrillo-their-exchange 'class:CLASS sec:ARRL_SECTION'`

###### Just grid square

VHF+ contests often only exchange callsign and grid locator.

`--cabrillo-my-exchange 'exch:MY_GRIDSQUARE'
--cabrillo-their-exchange 'exch:GRIDSQUARE'`

###### Signal report and serial number

For numeric exchanges, typically auto-incrementing.  Used by many CQ magazine
and RSGB contests.  Uses `=59` as a default so perfunctory signal reports don’t
need to be logged.

`--cabrillo-my-exchange 'rst:RST_SENT=59 nr:STX'
--cabrillo-their-exchange 'rst:RST_RCVD=59 nr:SRX'`

###### Signal report and exchange text

When no more-specific field fits, use `STX_STRING` (your exchange) and
`SRX_STRING` (their exchange).  This makes contest logging simple.

`--cabrillo-my-exchange 'rst:RST_SENT exch:STX_STRING'
--cabrillo-their-exchange 'rst:RST_RCVD exch:SRX_STRING'`

If data might be a serial number or might be a string, use both:

`--cabrillo-my-exchange 'rst:RST_SENT exch:STX/STX_STRING'
--cabrillo-their-exchange 'rst:RST_RCVD exch:SRX/SRX_STRING'`

###### Serial number and IOTA island ID

This exchange is used by the RSGB IOTA contest, and is an example of an optional
field (for stations which aren’t on an island).

`--cabrillo-my-exchange 'rst:RST_SENT ex1:STX ex2:MY_IOTA?'
--cabrillo-their-exchange 'rst:RST_RCVD ex1:SRX ex2:IOTA?'`

###### Serial number, precedence, check, ARRL section

This four-part exchange is used in ARRL Sweepstakes contests.  Both precedence
and check have ADIF fields, `MY_PRECEDENCE` and `MY_CHECK` are not part of the
ADIF 3.1.5 specification.

`--cabrillo-my-exchange 'nr:stx p:my_precedence ck:my_check sec:my_arrl_section'
--cabrillo-their-exchange 'nr:srx p:precedence ck:check sec:arrl_section'`

###### User-defined field

The RSGB Club Calls contest exchange includes signal report, serial number,
club status (broadcasting from club headquarters, associated with a club, or
not associated with a club), and the club the station is associated with.
Since neither the status nor club are ADIF fields, only one can use `SRX_STRING`
so a user-defined enumeration field `RSGB_STATUS` is added.  Note that status
and club aren’t part of the logging station’s side of the exchange, per
[RSGB Cabrillo example](https://www.rsgbcc.org/hf/rules/2021/Cabrillo-Information.shtml).
(Alternatively, the `CLASS` and `CHECK` fields could perhaps be used.)

`--cabrillo-my-exchange 'rst:RST_SENT nr:STX'
--cabrillo-their-exchange 'rst:RST_RCVD nr:SRX status:RSGB_STATUS club:SRX_STRING?'
--userdef 'RSGB_STATUS,{HQ,CM,NC}'`

###### DARC WAE QTC

QTC Traffic, as used in the DARC Worked All Europe contest, is not yet supported
since there are two separate sets of fields encoded in the same log.  If you are
interested in support, please open a GitHub issue and describe how this data is
stored in an ADIF file.

#### International text and Unicode

`adifmt` currently assumes all input files are encoded in
[UTF-8](https://en.wikipedia.org/wiki/UTF-8), which includes ASCII-only files.

For backwards-compatibility with ASCII-only software, the
[ADIF specification](https://adif.org.uk/314/ADIF_314.htm#Data_Types_Enumerations_and_Fields)
defines `Character` and `String` types as
[ASCII](https://en.wikipedia.org/wiki/ASCII)-only, with `IntlCharacter` and
`IntlString` as allowing any Unicode character (except line breaks unless in a
`IntlMultilineString` field).  Additionally, as of ADIF version 3.1.5, ADI
files are supposed to be ASCII-only and may not have `Intl*` fields.  ADIF
Multitool deviates from the spec by passing through Intl fields in ADI files
and writing Unicode characters in UTF-8.  (This allows ADI to be the default
output format in a pipeline of several commands, then save to a format which
allows Unicode.)  Unicode characters can be rejected in ADI files with the
`--adi-ascii-only` option, though if used with
`adifmt save --overwrite-existing` the file may be deleted before the program
aborts with an error; this will still output Intl fields if they contain only
ASCII characters.  `adifmt validate` ensures that “non-intl” fields are
ASCII-only; other commands pass through Unicode strings untouched.

### Conditions and Comparisons

Several `adifmt` commands can produce output only if a record matches one or
more conditions.  For example, `adifmt find` can be used to filter a larger log
file to a subset of records, like only CW contacts, or only QSOs on the 20
meter band.  Conditions are specified with one or more flag options, e.g. `--if
mode=CW` to match records where the `MODE` field is set to `CW`.  Conditions
can be negated, e.g. `--if-not mode=CW` to match all records **except** CW.  If
more than one condition is given, a record must match **all** of the `--if` and
`--if-not` conditions (boolean AND logic).  The `--or-if` and `--or-if-not`
flags introduce a boolean OR, matching if either all conditions before the flag
*or* all of the conditions after the flag are met.  A contrived example:
`adifmt find --if mode=CW --if-not band=20m --or-if tx_pwr=5 --if-not band=20m --or-if call=W1AW`
will filter a logfile, producing only records which are _either_ (a) CW contacts
_not_ on the 20 meter band, (b) 5 watt contacts (any mode) _not_ on the 20 meter
band, or (c) contacts with W1AW (any band, any mode).

In addition to equality checks, greater-than and less-than comparisons can be
used in a condition.  Comparisons use the type of the field, so numeric fields
like frequency and power sort numerically while digits in string fields sort
alphabetically.  For example, `FREQ>21` will match the frequency `146.52` but
`ADDRESS>21` will _not_ match someone whose address is `146 Main St` since `1`
comes before `2` in a string field.  The `--locale` flag indicates the language
to use for string comparisons, using the
[BCP-47 format](https://en.wikipedia.org/wiki/IETF_language_tag).  Available
comparisons are

* `field = value`: Case-insensitive equality, e.g. `contest_id=ARRL-field-day`
* `field < value`: Less than, `freq<29.701`
* `field <= value`: Less than or equal, `band<=10m`
* `field > value`: Greater than, `tx_pwr>100`
* `field >= value`: Greater than or equal, `qso_date>=20200101`

Fields can be compared to other fields by enclosing in `{` and `}`:

* `gridsquare={my_gridsquare}`: Contact with a station the same maidenhead grid
* `freq<{freq_rx}`: Operating split, with transmit below other station.

Conditions can match multiple values separated by `|` characters:

* `mode=SSB|FM|AM|DIGITALVOICE`: Any phone mode was used
* `arrl_sect={my_arrl_sect}|ENY|NLI|NNY|WNY` : Contact in the same ARRL section,
  or in New York

Conditions match fields with a list type if any value in the list matches.
If the `POTA_REF` field has value `K-0034,K-4556` then the record will match the
condition `--if pota_ref=K-4556` even though it doesn’t specify all the parks.

Empty or absent fields can be matched by omitting value:

* `operator=`: `OPERATOR` field not set
* `my_sig_info>`: `MY_SIG_INFO` field is set ("greater than empty")

Make sure to use quotes around conditions so that operators are not treated as
special shell characters:
  `adifmt find --if 'freq>=7' --if-not 'state={my_state}' --or-if 'tx_pwr<=5'`

The `--if`, `--if-not`, `--or-if`, and `--or-if-not` options are used by the
`edit` and `find` commands.  Field comparison rules are also used by `sort`.

“International” fields like `NAME_INTL` use Unicode sorting rules with a
language given by the `--locale` option, e.g. `--locale=da` for Danish or
`--locale=fr-CA` for Canadian French.  Non-international String fields like
`NAME` and `CALL` use basic ASCII sorting, regardless of locale.

Boolean fields sort false before true.  Integer and number fields compare by
numeric order.

Date and time fields are compared in chronological order.  In particular, the
time `123456` (4 seconds before 12:35 pm) is less than time `2030` (8:30 pm),
which would not be true if they were compared as numbers.

Latitude and longitude location fields are sorted west-to-east and
south-to-north so that string sorting by gridsquare has the same results as
sorting by latitude and then longitude.

Most enumeration fields use string sorting, but the `BAND` enum sorts
numerically by frequency ranges (so `40m`, `10m`, `70cm` are in order) and the
`DXCC Entity Code` enum sorts numerically, so DXCC code `7` (Albania) sorts
before `63` (French Guiana), which in turn sorts before `305` (Bangladesh).  To
sort alphabetically by country name, use the `COUNTRY` or `MY_COUNTRY` string
fields.

Several comparisons, including date, time, and location, are strict about field
format, so consider using `adifmt fix` and/or `adifmt validate` before
`adifmt find`, `adifmt edit`, or `adifmt sort`.  Missing or empty fields compare
as less than non-empty fields, and incorrectly formatted fields generally compare
as less than correctly formatted fields.

### Commands

ADIF Multitool behavior is organized into _commands_; each `adifmt` invocation
runs one command.  Commands are the first program argument, before any options
or file names: `adifmt command --some-option --other=value file1.adi file2.csv`

Name       | Description |
---------- | ----------- |
`cat`      | Concatenate all input files to standard output |
`count`    | Count records or unique field combinations |
`edit`     | Add, change, remove, or adjust field values |
`find`     | Include only records matching a condition |
`fix`      | Correct field formats to match the ADIF specification |
`flatten`  | Flatten multi-instance fields to multiple records |
`help`     | Print program, command, or format usage information |
`infer`    | Add missing fields based on present fields |
`save`     | Save standard input to file with format inferred by extension |
`select`   | Print only specific fields from the input |
`sort`     | Sort records by a list of fields |
`validate` | Validate field values; non-zero exit and no stdout if invalid |
`version`  | Print program version information |

`adifmt help` will also show this list.

#### help

`adifmt help` prints usage information, a list of supported formats, available
commands, and options which apply to any command.  `adifmt help cmd` prints
usage information about and options for command `cmd`.  `adifmt help fmt`
prints options for input/output format `fmt`.  There are a lot of options, so
consider running `adifmt help | less`.

#### cat

`adifmt cat` reads all input records and prints them to standard output.  Given
several input files (perhaps one per day, callsign, or location) `cat` will
combine them into a single file.  `cat` can also be used to convert from one
format to another, e.g. `adifmt cat --output=csv mylog.adi` to convert from ADI
format to CSV.  (If `--input` is not specified the file type is inferred from
the file name; if `--output` is not specified ADI is used.)

#### count

`adifmt cat` groups equal field values and adds a field with the number of times
the input had each group.  If no `--fields` are given, it outputs a single
record with a single field with the number of records in the input, not
including headers.  If one or more `--fields` are given, outputs one record for
each unique combination of those fields, with a count of the occurrences.
Fields which are unset in a record contribute an empty string to the
combination; "all empty" is a valid group.

The result of `count` is valid ADIF in the format specified by `--output`; using
TSV or CSV may make the results easy to consume.  The name of the added field
can be set like `--count-field=NUM`, it defaults to `APP_ADIFMT_COUNT`.  The
output is sorted in the order given by `--fields` and in the natural order for
each field (strings are case-insensitive, numeric fields sort by value; band
sorts numerically, longitude and latitude sort in the same order as gridsquare.)

This can be combined with `find` to discover duplicate QSOs with a given
uniqueness criterion:

```
adifmt count --fields qso_date,call,band,mode --count-field num mylog.adi | \
  adifmt find --if 'num>1'
```

#### edit

`adifmt edit` adds, changes, or removes fields in each input record.
Options can be specified multiple times, e.g.
`adifmt edit --add my_gridsquare=FN31pr --add "my_name=Hiram Percy Maxim" log.adi`

The `--set` option (`name=value`) changes the value of the given field on all
records, adding it if it is not present.  The `--add` option (`name=value`)
only adds the field if it is not already present in the record.  The `--rename`
option (`old=new` field names) changes an old field name to a new one.  The
`--remove` option (field names, optionally comma-separated) deletes the field
from all records.  The `--remove-blank` removes all blank fields (string
representation is empty).

The `--time-zone-from` and `--time-zone-to` options will shift the `TIME_ON` and
`TIME_OFF` fields (along with `QSO_DATE` and `QSO_DATE_OFF` if applicable) from
one time zone to another, defaulting to UTC.  For example, if you have a CSV
file with contact times in your local QTH in New South Wales you can convert it
to UTC (Zulu time) with `adifmt edit --time-zone-from Australia/Sydney file.csv`.

Edits can be applied to only records matching a condition, using the
[Conditions and Comparisons](#conditions-and-comparisons) options.  Records
which do not match the conditions will be output unchanged.  If different edits
should be applied based on different conditions, multiple edit commands should
be chained together in a pipeline.  For example, to set the `SUBMODE` for SSB
contacts to upper sideband on the 20 meter and higher bands and to lower
sideband for the 40, 80, and 160 meter bands, express each edit as a condition
and a change:

```sh
adifmt cat mylog.adi \
  | adifmt edit --if 'mode=SSB' --if 'band>=20m' --or-if 'band=60m' --add 'submode=USB' \
  | adifmt edit --if 'mode=SSB' --if 'band=40m|80m|160m' --add 'submode=LSB' \
  | adifmt save fixed_sideband.adi
```

#### find

`adifmt find` filters the input, outputting only records which match one or more
conditions.  For details on condition syntax, see
[Conditions and Comparisons](#conditions-and-comparisons) above.  An example
which finds all records where the contest ID is set to ARRL Field Day but
ignoring records on the WARC bands (60, 30, 17, and 12 meters) is
`adifmt find --if 'contest_id=ARRL-FIELD-DAY' --if-not 'band=60m|30m|17m|12m'`

#### fix

`adifmt fix` coerces some fields into the format dictated by the ADIF
specification.  The rule of thumb for default fixes is that they should be
unsurprising to almost anyone, like converting `3:45 PM` to `1545` for a time
field.  Currently only date, time, and location fields are coerced.  Dates must
already be in year, month, day order.  Location fields can be converted from
decimal (GPS) coordinates to degrees/minutes.

`fix` also changes [ISO 3166-1 alpha-2 and alpha-3](https://en.wikipedia.org/wiki/ISO_3166-1)
codes in the `COUNTRY` and `MY_COUNTRY` to
[DXCC entity names](https://adif.org.uk/314/ADIF_314.htm#DXCC_Entity_Code_Enumeration)
if a match is found.  This can save a lot of typing for `BA` -> `BOSNIA-HERZEGOVINA`
or `USA` → `UNITED STATES OF AMERICA`  Note that some DXCC entities like
Alaska, Hawaii, Crete, Corsica, Sardinia, many other remote islands, and
international organizations do not have ISO 3166 codes.  A few countries do not
have a single DXCC entity for “the mainland”, including the United Kingdom
(separated into England, Wales, Scotland, and Northern Ireland), Russia
(European Russia, Asiatic Russia, and Kaliningrad), Kiribati (separated into
island chains), and a few dependent island territories.  Country code
translations will not be applied for those since it’s not obvious which DXCC
entity was contacted.

In the future, other formats may be fixable, including varieties of the Boolean
data types, forcing some string fields to upper case, and perhaps correcting
some other common variations on enum fields as is done with countries.  A
future update will also provide options like date formats so that
day/month/year or month/day/year input data can be unambiguously fixed.

#### flatten

`adifmt flatten` converts single records with a multi-instance field into
multiple records with a single value for that field.  Non-flattened fields are
included unchanged in each record.  This can be useful when processing the
output with tools which don’t expect a list of values in a field, e.g. counting
the number of contacts you’ve made with each grid square while treating
contacts on the border of a square as separate:
`adifmt flatten --fields VUCC_GRIDS | adifmt count --fields VUCC_GRIDS --output tsv`

The `flatten` command will turn

```
CALL	VUCC_GRIDS
W1AW	EN98,FM08,EM97,FM07
AH1Z	FM07,FM08
```

into

```
CALL	VUCC_GRIDS
W1AW	EN98
W1AW	FM08
W1AW	EM97
W1AW	FM07
AH1Z	FM07
AH1Z	FM08
```

and the rest of the pipeline will produce grid counts like

```
NUM VUCC_GRIDS
1   EM97
1   EN98
2   FM07
2   FM08
```

If multiple fields are flattened and each has multiple instances, a Cartesian
combination will be output.  For example, if `MY_POTA_REF` has two POTA
references and `POTA_REF` has three POTA references on one record, six records
will be output, one for each pair.  As of June 2024, POTA uploads don’t handle
multi-instance `POTA_REF` fields, so

```sh
adifmt flatten --fields POTA_REF,MY_POTA_REF \
  | adifmt infer --fields SIG_INFO,MY_SIG_INFO \
  | adifmt save '{station_callsign}@{my_sig_info}-{qso_date}.adi'
```

is needed to get full credit for park-to-park 2-fers.

The delimiter (usually a comma, except SecondarySubdivisionList which uses a
colon) is implied by the field’s data type in the ADIF spec.  You may specify
the delimiter for a field with the `--delimiter field=delim` flag, make sure to
quote any special shell characters, e.g.
`adifmt flatten --fields STX_STRING --delimiter 'STX_STRING=;'`  Escape
sequences can be surrounded with single or double quotes, e.g.
`adifmt flatten --fields COMMENT --delimiter 'COMMENT="\r\n"' (the `'`s are
seen by the shell so the `"`s are passed to the program) or
`--delimiter "COMMENT='\u3000'" (the `"`s are seen by the shell, the `'`s by
adifmt, and can only contain a single character).  Double-quoted delimiters are
interpreted as a [Go string literal](https://go.dev/ref/spec#String_literals)
and single-quoted as a [rune literal](https://go.dev/ref/spec#Rune_literals).

#### infer

`adifmt infer` guesses the value for fields which are not present in a record.
Field names to infer are given by the `--fields` option, which can be repeated
multiple times and/or comma-separated.  Fields in the list will not be changed
if they are present in a record with a non-empty value.

`SIG_INFO` and `MY_SIG_INFO` are handled specially.  If `SIG`/`MY_SIG` is is
present, that value determines which field to use for `SIG_INFO`/`MY_SIG_INFO`.
For example, if `SIG` is `SOTA`, `SIG_INFO` will be set to the value of
`SOTA_REF` even if `POTA_REF` is also present.  If `SIG`/`MY_SIG` is absent,
all special activity fields (IOTA, POTA, SOTA, WWFF) will be checked.  If
exactly one of them is present, that value will be used for
`SIG_INFO`/`MY_SIG_INFO` and `SIG`/`MY_SIG` will be set to the activity name.
If (`MY_`)`SIG` is set to a special interest activity or event that does not
have a dedicated ADIF field (e.g. 13 Colonies, Volunteers on the Air),
(`MY_`)`SIG_INFO` will not be inferred.

Inferable fields:

* `BAND` from `FREQ`
* `BAND_RX` from `FREQ_RX`
* `MODE` from `SUBMODE`
* `COUNTRY` from `DXCC`
* `MY_COUNTRY` from `MY_DXCC`
* `DXCC` from `COUNTRY`
* `MY_DXCC` from `MY_COUNTRY`
* `CQZ` from `COUNTRY`/`DXCC`
* `MY_CQ_ZONE` from `MY_COUNTRY`/`MY_DXCC`
* `ITUZ` from `COUNTRY`/`DXCC`
* `MY_ITU_ZONE` from `MY_COUNTRY`/`MY_DXCC`
* `CONT` from `COUNTRY`/`DXCC`
* `CNTY` from `USACA_COUNTIES` (unless multiple county-line counties)
* `MY_CNTY` from `MY_USACA_COUNTIES` (unless multiple county-line counties)
* `USACA_COUNTIES` from `CNTY` (if a USA DXCC entity)
* `MY_USACA_COUNTIES` from `MY_CNTY` (if a USA DXCC entity)
* `GRIDSQUARE` and `GRIDSQUARE_EXT` from `LAT`/`LON`
* `MY_GRIDSQUARE` and `MY_GRIDSQUARE_EXT` from `MY_LAT`/`MY_LON`
* `OPERATOR` from `GUEST_OP`
* `STATION_CALLSIGN` from `OPERATOR` or `GUEST_OP`
* `OWNER_CALLSIGN` from `STATION_CALLSIGN`, `OPERATOR`, or `GUEST_OP`
* `SIG_INFO` from one of `IOTA`, `POTA_REF`, `SOTA_REF`, or `WWFF_REF` based on
  `SIG` (sets `SIG` if unset and only one of the others is set)
* `MY_SIG_INFO` from one of `MY_IOTA`, `MY_POTA_REF`, `MY_SOTA_REF`, or
  `MY_WWFF_REF` based on `MY_SIG` (sets `MY_SIG` if unset and only one of the
  others is set)
* `IOTA`, `POTA_REF`, `SOTA_REF`, and `WWFF_REF` from `SIG_INFO` if `SIG` is
  set to the appropriate program.
* `MY_IOTA`, `MY_POTA_REF`, `MY_SOTA_REF`, and `MY_WWFF_REF` from `MY_SIG_INFO`
  if `MY_SIG` is set to the appropriate program.

#### save

`adifmt save` writes ADIF records from standard input to a file.  The output
format is inferred from the file name or can be given explicitly with
`--output`.  Existing files will not be overwritten unless the
`--overwrite-existing` option is given.  The output file will not be written
(and will exit with a non-zero code) if there are no records in the input; this
allows a chain like `adifmt fix log.adi | adifmt validate | adifmt save
--overwrite-existing log.adi` which will attempt to fix any errors in `log.adi`
and save back to the same file, but which won’t clobber it if validation still
fails.  Writing a zero-record file can be forced with `--write-if-empty`.

`save` can split the input into multiple files based on a filename template.
The template uses field names in curly braces: `{FIELD_NAME}`, which is not
case-sensitive.  Enclose the template in quotes to avoid shell metacharacters.
For example, `adifmt cat log.adi | adifmt save '{BAND}-{MODE}.adi'` writes each
band/mode pair to a separate file, perhaps producing `10M-SSB.adi 10M-FM.adi
20M-CW.adi 20M-DIGITAL.adi 20M-SSB.adi 40M-CW.adi 80M-SSB.adi`.  Another example
using the [Parks on the Air filename format](https://docs.pota.app/docs/activator_reference/submitting_logs.html)
is `adifmt save '{station_callsign}@{my_sig_info}-{qso_date}.adi'`.  All field
values will be converted to upper case and special file system characters are
replaced by `-` (so `{CALL}.csv` with `w1aw/2` becomes `W1AW-2.csv`).  Fields
without a value are replaced with `FIELD_NAME-EMPTY`.  Special characters in the
template itself are not replaced, and can be used to split a log into separate
directories: `adifmt save --create-dirs '{operator}/{band}.adx`.

#### select

`adifmt select` outputs only the specified fields.  Currently each field must
be specified by name, either in a comma-separated list or by specifying the
`--field` option multiple times.  The following uses are equivalent:

```sh
adifmt select --fields call,qso_date,time_on,time_off mylog.adi
adifmt select --fields call --fields qso_date --fields time_on,time_off mylog.adi
```

`select` can be effectively combined with other standard Unix utilities.  To
find duplicate QSOs by date, band, and mode, use
[sort](https://man7.org/linux/man-pages/man1/sort.1.html) and
[uniq](https://man7.org/linux/man-pages/man1/uniq.1.html).
(See [`count`](#count)` for another approach.)

```sh
adifmt select --fields call,qso_date,band,mode --output tsv --tsv-omit-header mylog.adi \
  | sort | uniq -d
```

This is similar to a SQL `SELECT` clause, except it cannot (yet?) transform the
values it selects.

#### sort

`adifmt sort` sorts records by one or more fields, specified by the `--fields`
option.  A field name can be prefixed with a minus sign (`-`) to sort that field
in descending order.  See the [Conditions and Comparisons](#conditions-and-comparisons)
section for details about data type ordering.  For example, to sort a log by
callsign of the contacted station (ascending) in reverse chronological order:

```sh
adifmt sort --fields call,-qso_date,-time_on mylog.adi
```

The `--locale` option will use language-specific rules for sorting international
strings, e.g. `adifmt sort --locale=da --fields QTH_INTL` will use the
alphabetic order for Danish and Norwegian, producing
`Arendal, Bergen, Oslo, Trondheim, Ænes, Østfold, Ålgård` while using
`--locale=en` will use an English sort order which treats Æ, Ø, and Å as
accented letters, sorted as AE, O, and A respectively.

#### validate

`adifmt validate` checks that field values match the format and enumeration
values in [the ADIF specification](https://adif.org.uk/adif).  Errors and
warnings are printed to standard error.  If any field has an error, nothing is
printed to standard output and exit status is `1`; if no errors are present (or
only warnings), the input will be printed to standard output as in
[`cat`](#cat) and exit status is `0`.  If the output format is ADI or ADX,
warnings will be included as record-level comments in the output.

Validations include field type syntax (e.g. number and date formats);
enumeration values (e.g. modes and bands), and number ranges.  The ADIF
specification allows some fields to have values which do not match the
enumerated options, for example the `SUBMODE` field says “use enumeration values
for interoperability” but the type is string, allowing any value.  These
warnings will be printed to standard error with `adifmt validate` but will not
block the logfile from being printed to standard output.  Dates and times in the
future (based on the computer’s current wall clock) will print a warning; there
is not currently a way to override the current time.

The `--required-fields` option provides a list of fields which must be present
in a valid record.  Multiple fields may be comma-separated or the option given
several times.  For example, checking a contest log might use `adifmt validate
--reqiured-fields qso_date,time_on,call,band,mode,srx_string` The
`--if --if-not --or-if --or-if-not` options will only check `--required-fields`
on records where the condition matches.  For example, to ensure SUBMODE is set
for MSFK and SSB contacts without worrying about CW, RTTY, or FM contacts:
`adifmt --required-fields=submode --if mode=MFSK --or-if mode=SSB`.  Data type
validity checks will occur even if the condition does not match.

Some but not all validation errors can be corrected with [`adifmt fix`](#fix).

#### version

`adifmt version` prints the version number of the installed program, the ADIF
specification version, and URLs to learn more.

### Future features (under construction)

ADIF Multitool was created because I was recording
[Parks on the Air](https://parksontheair.com/) logs on paper and then typing
them into a spreadsheet. I needed a way to convert exported CSV files into ADIF
format for upload to [the POTA website](https://pota.app/) while fixing
incompatibilities between the spreadsheet data format and the expected ADIF
structure.  I decided to solve this problem with a “Swiss Army knife for ADIF
files” following the
[Unix pipeline philosophy](https://en.wikipedia.org/wiki/Pipeline_\(Unix\)) of
simple tools that do one thing and can be easily composed together to build more
powerful expressions.

There are a lot of things that a ham radio log file program could do, and I
would like `adifmt` to do many of them. The program has most core features I've
planned to add.  If you've got a use case for working with ADIF files
that `adifmt` can’t do yet, please create a GitHub issue to discuss how it
might work.

Further features I plan to add:

*   Validate more fields.
*   Identify duplicate records using flexible criteria, e.g., two contacts with
    the same callsign on the same band with the same mode on the same Zulu day
    and the same `MY_SIG_INFO` value.  (`count` plus `select` can do this, but
    does not print the full records.)  I would also like a way to combine
    duplicate records into one, e.g. reversing the `flatten` operation.
*   Option for `save` to append records to an existing ADIF file.
*   [FLE (fast log entry)](https://df3cb.com/fle/documentation/) format support.
*   Support for Cabrillo 2.0 format if needed.

See the [issues page](https://github.com/flwyd/adif-multitool/issues) for more
ideas or to suggest your own.

### Non-goals

I don't expect ADIF Multitool to support the following use cases. A different
piece of software will be needed.

*   Upload logs to any service like QRZ, eQSL, or LotW.  `adifmt` is a useful
    tool in preparing logs for upload, though.
*   Log-editing GUI. `adifmt` is a command-line tool; a GUI could be built which
    uses it to make edits, but that would be a separate program and project. I
    am open to the idea of an interactive console mode, though.
*   Live logging. `adifmt` is meant for processing logs that have already been
    created, not for logging contacts as they happen over the air. There are
    many fine amateur radio logging programs, most of which can export ADIF
    files that `adifmt` can process.  You could also keep logs in a text file,
    massage it to a CSV or TSV, and then process it with `adifmt`.

## Scripting and compatibility

ADIF Multitool is designed to be easy to include in scripts.  If you have a
workflow for dealing with ham radio logs, such as converting from CSV, adding
fields, and validating field syntax before uploading to the POTA or SOTA
websites, consider automating that process with `adifmt`.

ADIF Multitool is still “version zero” and the command line interface should be
considered unstable.  If you use v0 `adifmt` in a script or other program, be
prepared to update your code if commands or options change.  In particular, use
GNU-style double dashes for options (`--input`) rather than Go-style single
dashes (`-input`); the program may change to a GNU/POSIX-style flag-parsing
library which requires double dashes.

The `adif` and `cmd` packages should be considered “less stable” than the CLI
during the v0 phase and may undergo significant change.  Use of those packages
in your own program should only be done with significant tolerance to churn.

The v1 and future releases will follow
[Semantic Versioning](https:///semver.org/) and any breaking changes to the CLI
or public Go APIs will need to wait for v2.  ADIF spec updates and new features
will lead to a new minor version and bug fixes will increment the patch number.
If you find this useful as a library, please let me know.

I have not yet tested this on Windows; please
[report issues](https://github.com/flwyd/adif-multitool/issues) if anything does
not work, or is particularly awkward.

## Contributions welcome

ADIF Multitool is open source, using the Apache 2.0 license.  It is written in
the [Go programming language](https://go.dev/).  Bug fixes, new features, and
other contributions are welcome; please read the [contributing](CONTRIBUTING.md)
and [code of conduct](CODE_OF_CONDUCT.md) pages.  The primary author is Trevor
Stone, WT0RJ, @flwyd.

### Source Code Headers

Every file containing source code must include copyright and license
information.  Use the [`addlicense` tool](https://github.com/google/addlicense)
to ensure it’s present when adding files: `addlicense .`

Apache header:

```
Copyright 2025 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

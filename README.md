# ADIF Multitool

Validate, modify, and convert ham radio log files with a handy command-line
tool. 📻🌳🪓

`adifmt` provides a suite of commands for working with [ADIF](https://adif.org/)
logs from ham radio software. Each `adifmt` invocation reads log files from the
command line or standard input and prints an ADIF log to standard output,
allowing multiple commands to be chained together in a pipeline. For example, to
add a `BAND` field based on the `FREQ` (radio frequency) field, add your
station's maidenhead locator (`MY_GRIDSQURE`) to all entries, and save a log
file containing only SSB voice contacts, a pipeline might look like

```sh
adifmt infer -field band my_original_log.adi \
  | adifmt edit -add my_gridsquare=FN31pr \
  | adifmt filter -field mode=SSB \
  > my_ssb_log.adi
```

*Note*: `adifmt` is pronounced "ADIF M T" or "ADIF multitool", not "adi fmt" nor
"addy format".

## Quick start

ADIF Multitool is not yet available as a binary distribution, so you will need
a Go compiler on your system.  To check, run `go version` at the command line.
If the `go` program is not found, [download and install it](https://go.dev/dl/).
Then run `go install github.com/flwyd/adif-multitool` to make the `adifmt`
command available.  (You may need to add the `$GOBIN` environment variable to
your path.)  To see if it works, run `adifmt -help`.  If the command is not
found, try `go run github.com/flwyd/adif-multitool/adifmt -help`

To do something useful with ADIF Multitool, the syntax is

```
adifmt command [flags] files...
```

For example, the `cat` command concatenates all input files and outputs ADIF
data to standard output:

```sh
adifmt cat log1.adi log2.adi > combined.adi
```

prints all of the records in the two `logX.adi` files to the `combined.adi`
file.

Flags control input and output options.  For example, to print records with
a UNIX newline between fields, two newlines between records, and use lower
case for all field names:

```sh
adifmt cat -adi-field-separator=newline \
  -adi-record-separator=2newline \
  -adi-lower-case \
  log1.adi
```

Multiple input and output formats are supported (currently ADI per the ADIF
spec and CSV with field names matching the ADIF list).

```sh
adifmt cat -input=adi -output=csv log1.adi > log1.csv
adifmt cat -input=csv -output=adi log2.adi > log2.adi
```

`-input` need not be specified if it’s implied by the file name, and
`-ouput=adi` is the default.

If no file names are given, input is read from standard input:

```
gunzip --stdout mylog.csv.gz | adifmt cat -input=csv | gzip > mylog.adi.gz
```

This will be useful in composing several `adifmt` invocations together, once
more commands than `cat` are supported.

## Features (under construction)

ADIF Multitool was created because I was recording
[Parks on the Air](https://parksontheair.com/) logs on paper and then typing
them into a spreadsheet. I needed a way to convert exported CSV files into ADIF
format for upload to [the POTA app](https://pota.app/) while fixing
incompatibilities between the spreadsheet data format and the expected ADIF
structure. I decided to solve this problem with a "Swiss Army knife for ADIF
files" following the
[Unix pipeline philosophy](https://en.wikipedia.org/wiki/Pipeline_\(Unix\)) of
simple tools that do one thing and can be easily composed together to build more
powerful expressions.

There are a lot of things that a ham radio log file could do, and I would like
`adifmt` to do many of them. Development is in the early stages, so it doesn't
do very much *yet*. If you've got a use case for working with ADIF files, please
create a GitHub issue to discuss how it might work.

Features I plan to add:

*   Support several input and output formats: ADI and ADX from ADIF along with
    CSV and JSON. Maybe convert to and from Cabrillo format for contests.
*   Validate (and fix where possible) fields which don't meet ADIF expectations.
    For example, ensure that date and time strings match the ADIF format, check
    that enumeration fields values are in the enum list, etc.
*   Filter a log to only records matching some criteria, similar to a SQL
    `WHERE` clause.
*   Print only certain fields for each record, similar to a SQL `SELECT` clause.
*   Add, edit, or remove specific fields for each record, e.g., setting your
    station's location.
*   Infer missing fields based on the values of other fields. For example, the
    `BAND` field can be inferred from the frequency; `MODE` can be inferred from
    `SUBMODE`; `OPERATOR`, `STATION_CALLSIGN`, and `OWNER_CALLSIGN` can stand in
    for each other; `GRIDSQUARE` can be determined from `LAT` and `LON`; and a
    missing `DXCC` or `COUNTRY` field can be determined from the value of the
    other one.
*   Identify duplicate records using flexible criteria, e.g., two contacts with
    the same callsign on the same band with the same mode on the same Zulu day
    and the same `MY_SIG_INFO` value.
*   Count the total number of records or the number of distinct values of a
    field.

### Non-goals

I don't expect `adifmt` to support the following use cases. A different piece of
software will be needed.

*   Upload logs to any service like QRZ, eQSL, or LotW.
*   Log-editing GUI. `adifmt` is a command-line tool; a GUI could be built which
    uses it to make edits, but that would be a separate program and project. I
    am open to the idea of an interactive console mode, though.
*   Live logging. `adifmt` is meant for processing logs that have already been
    created, not for logging contacts as they happen over the air. There are
    many fine amateur radio logging programs, most of which can export ADIF
    files that `adifmt` can process.

## Contributions welcome

ADIF Multitool is open source, using the Apache 2.0 license. It is written in
the [Go programming language](https://go.dev/). Bug fixes, new features, and
other contributions are welcome; please read the [contributing](CONTRIBUTING.md)
and [code of conduct](CODE_OF_CONDUCT.md) pages.

### Source Code Headers

Every file containing source code must include copyright and license
information.  Use the [`addlicense` tool](https://github.com/google/addlicense)
to ensure it’s present when adding files:
`addlicense -c “Google LLC” -l apache .`

Apache header:

```
Copyright 2022 Google LLC

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

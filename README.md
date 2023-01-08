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

Commands can be combined in a Unix-style pipeline.  The `fix` command
automatically changes some values to match the expected ADIF format such as
changing a time field from `12:34:56` to `123456` and a date from `2012-03-04`
to `20120304`.  The `select` command prints only a subset of fields.  These can
be combined:

```sh
adifmt fix log1.adi | adifmt select -fields qso_date,time_on,call > minimal.adi
```

creates a file named `minimal.adi` with just the date, time, and callsign from
each record in the input file `log1.adi`.

## Features

### Commands

Name     | Description |
-------- | ----------- |
`cat`    | Concatenate all input files to standard output |
`fix`    | Correct field formats to match the ADIF specification |
`select` | Print only specific fields from the input; skip records with no matching fields |

#### cat

`adifmt cat` reads all input records and prints them to standard output.  Given
several input files (perhaps one per day, callsign, or location) `cat` will
combine them into a single file.  `cat` can also be used to convert from one
format to another, e.g. `adifmt cat -output=csv mylog.adi` to convert from ADI
format to CSV.  (If `-input` is not specified the file type is inferred from the
file name; if `-output` is not specified ADI is used.)

#### edit

`adifmt edit` adds, changes, or removes fields in each input record.  Flags can
be specified multiple times, e.g.
`adifmt edit -add my_gridsquare=FN31pr -add “my_name=Hiram Percy Maxim” log.adi`

The `-set` flag (`name=value`) changes the value of the given field on all
records, adding it if it is not present.  The `-add` flag (`name=value`) only
adds the field if it is not already present in the record.  The `-remove` flag
(field names, optionally comma-separated) deletes the field from all records.
The `-remove-blank` removes all blank fields (string representation is empty).

#### fix

`adifmt fix` coerces some fields into the format dictated by the ADIF
specification.  The rule of thumb for default fixes is that they should be
unsurprising to almost anyone, like converting `3:45 PM` to `1545` for a time
field.  Currently only date and time fields are coerced, and dates must already
be in year, month, day order.  In the future, other formats may be fixable,
including varieties of the Boolean data types, decimal coordinates to
degrees/minutes/seconds, forcing some string fields to upper case, and perhaps
correcting some common variations on enum fields, e.g. `USA` →
`UNITED STATES OF AMERICA`.  A future update will also provide flags like date
formats so that day/month/year or month/day/year input data can be unambiguously
fixed.

#### select

`adifmt select` outputs only the specified fields.  Currently each field must
be specified by name, either in a comma-separated list or by specifying the
`-field` flag multiple times.  The following uses are equivalent:

```sh
adifmt select -fields call,qso_date,time_on,time_off mylog.adi
adifmt select -fields call -fields qso_date -fields time_on,time_off mylog.adi
```

`select` can be effectively combined with other standard Unix utilities.  To
find duplicate QSOs by date, band, and mode, use
[sort](https://man7.org/linux/man-pages/man1/sort.1.html) and
[uniq](https://man7.org/linux/man-pages/man1/uniq.1.html):

```sh
adifmt select -fields call,qso_date,band,mode -output csv mylog.adi \
  | sort | uniq -d
```

This is similar to a SQL `SELECT` clause, except it cannot (yet?) transform the
values it selects.

### Future features (under construction)

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
    field.  (The total number of records can currently be counted with
    `-output=csv`, piping the output to `wc -l`, and subtracting 1 for the
    header row.)

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

If you would like to use the `adif` package in your own Go programs, be aware
that this code is still in a “version zero” state and APIs may change without
notice.  If you find this useful as a library, please let me know.

### Source Code Headers

Every file containing source code must include copyright and license
information.  Use the [`addlicense` tool](https://github.com/google/addlicense)
to ensure it’s present when adding files:
`addlicense -c “Google LLC” -l apache .`

Apache header:

```
Copyright 2023 Google LLC

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

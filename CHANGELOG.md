# Changelog

All notable changes to ADIF Multitool will be documented in this file.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

ADIF Multitool is still in the “v0” phase, and some user-facing functionality
may change, though so far the interface has been stable.  Users calling `adifmt`
from scripts are encouraged to use the GNU/POSIX two-dash `--option` style
rather than the Go one-dash `-option` style, since support for the latter could
be dropped.  During v0 there is no commitment to API stability of the `adif`,
`adif/spec`, or `cmd` Go packages.  If you use this as a library in your own Go
program I would appreciate hearing about it so the API can be evolved with your
input.

## [Unreleased]

### Added

* `count` command prints the number of records of each combination of a set of
  fields.  For example, `adifmt count --fields BAND,MODE` prints the number of
  contacts made for each band plus mode combination.  If a field is blank or not
  set on a particular record it contributes an empty string to the combination.
  `adifmt count` without the `--fields` option prints the total number of
  records in the file.  CSV or TSV output makes `count` output easy to pass to
  another program like a spreadsheet for further analysis.
* `CONT` (continent) works with `adifmt infer` and `adifmt validate` based on
  the `DXCC` or `COUNTRY` field.
* `validate --required-fields` supports conditionals (`--if`, `--if-not`,
  `--or-if`, `--or-if-not`) to determine which records must have certain fields
  set.  All records are checked for data type validity, e.g. number formatting.

[CQ Zones](https://mapability.com/ei8ic/maps/cqzone.php) and
[ITU Zones](https://mapability.com/ei8ic/maps/ituzone.php):

* `infer` will set the relevant zone field if the `DXCC` or `COUNTRY` field is
  set and the entity is in a single zone.  It will also infer the zone if the
  `STATE` field is set and the state is only in one zone and the ADIF primary
  administrative division enumeration indicates the zone of the subdivision.
  For example, `COUNTRY=India` will set `CQZ=22 ITUZ=41` because India does not
  overlap CQ or ITU zone boundaries.  `COUNTRY=Kazakhstan` will set `CQZ=17` but
  not set `ITUZ` at all, since Kazakhstan straddles the ITU Zone 29/30 boundary
  and ADIF does not define any subdivisions for the country.  `DXCC=291` (USA)
  would infer nothing because it’s in both multiple CQ Zones and multiple ITU
  Zones.  `DXCC=291 STATE=CA` would infer `CQZ=3 ITUZ=6` since California is in
  just one zone of each type, but `DXCC=291 STATE=MT` would just infer `CQZ=4`
  since Montana is in two ITU Zones.
* `validate` will consider a record valid if the CQ and ITU Zones (if present)
  are consistent with the DXCC and primary subdivision.  No additional location
  checks, e.g. comparing gridsquare to zone boundaries, is done.
* `MY_CQ_ZONE` and `MY_ITU_ZONE` work as well.
* CQ and ITU Zones also work for DXCC entities which have been removed from the
  active list, e.g. Zanzibar.
* Some multi-zone countries are missing zone data in the primary administrative
  subdivision enumeration.
  [An ADIF proposal](https://groups.io/g/adifdev/topic/cq_zone_in_subdivision_and/110819473)
  to expand and correct this data has been proposed, but it may be some time
  before the next ADIF version.  ADIF Multitool may add its own subdivision
  associations before then.
* There is some variance between amateur radio organizations on exactly where
  ITU boundaries lie.  This program uses the
  [CQWW WAZ list](https://cqww.com/cq_waz_list.htm) for CQ-zone-to-entity
  associations, the
  [AARL DXCC list](https://www.arrl.org/files/file/DXCC/2022_Current_Deleted.txt)
  for ITU-to-entity associations, and [zone-check.eu](https://zone-check.eu/) to
  determine which ITU Zones a subdivision crosses.  See also the mapability
  links above for some explanation of data challenges.  Maritime boundaries
  between zones are not well defined and this program does not attempt to
  map geographic points to zones (or to geopolitical entities).

### Changed

Started a [changelog](CHANGELOG.md) file so it’s easier to learn what’s new in
a release.

### Fixed

* Consistently trim leading and trailing space from ADI and ADX comments.  If a
  comment is nothing but space, don’t preserve it in the output.
  **Note:** record comments and file comments are preserved from input through
  output, but header comments, such as the text at the top of an ADI file, are
  ignored.  This could change in the future, but I don’t want to preserve all
  the “Generated” lines created as part of a pipeline.
  [See discussion](https://github.com/flwyd/adif-multitool/discussions/13).
* Franz Josef Land DXCC entity is part of Russia, Arkhangelsk Oblast.

### Removed

Nothing yet


## [v0.1.18] - 2024-09-26

### Added

Upgrade data specification to [ADIF 3.1.5](https://adif.org/adif).

### Security

Upgrade golang.org/x/crypto package.


## [v0.1.17] - 2024-09-26

### Added

`--field-order` option to control the order of fields in the output.  Any fields
not in the list will still be included in the output.

Add some practical command examples to the README.

## [v0.1.16] - 2024-09-18

### Added

* Cabrillo: flexible ADIF field mapping to Cabrillo exchange.  The WWROF
  Cabrillo documentation gives an _example_ exchange format (RST and SRX/STX),
  not the only valid exchange format.  See [README](README.md) or
  `adifmt help cabrillo` for examples of `--cabrillo-their-exchange` and
  `--cabrillo-my-exchange` options.
* CSV and TSV: `--csv-omit-header` and `--tsv-omit-header` to only print
  records, no header row.  This makes it easy to process the output with other
  POSIX utilities like `sort | uniq` without needing to pipe through `tail +2`
  first.
* Filenames can now come before or after `--` options, e.g.
  `adifmt command file1.adi --option1=foo --option2=bar file2.adi`.  Previously
  filenames needed to be after all `--` options, which made it challenging to
  replace a filename from the end of the first step of a long pipeline.

### Fixed

* Fix crash in `sort` when sorting fields with a list of strings.
* Cabrillo: treat multiple phone modes (e.g. SSB, AM, and FM) as SSB and treat
  RTTY plus other digital modes as DIGI, rather than MIXED.
* Cabrillo: support tab field delimiter, don’t leave dangling spaces, and other
  minor fixes.
* ADX: output a newline at the end of the file.
* Fix user-defined field format in `help` output.

## [v0.1.15] - 2024-08-18

### Added

`validate`: warn if a date or time is in the future.

### Fixed

Support whitespace delimiters in `flatten`, e.g.
`adifmt flatten --fields SRX_STRING --delimiter 'SRX_STRING= '` to turn a QSO
party county-line record with exchange `ABC DEF` into two QSOs.  Delimiters may
also be quoted, e.g. `--delimiter 'foo="\t"'`.

Internal improvements to make command-line options easier to test.


## [v0.1.14] - 2024-08-12

### Added

`adifmt infer --fields USACA_COUNTIES,MY_USACA_COUNTIES` from the `CNTY` and
`MY_CNTY` in the US.  Also infers `CNTY` from USACA if there’s only one in the
list.


## [v0.1.13] - 2024-07-05

### Added

* `flatten` command to convert multi-value fields into multiple records.  This
  is particularly useful along with `adifmt infer` for turning multi-park
  contacts (`POTA_REF=US-4572,US-4576`) into multiple QSOs
  (`SIG=POTA SIG_INFO=US-4572` and `SIG=POTA SIG_INFO=US-4576`), as expected by
  POTA uploads.
* Cabrillo: `--their-exchange-field-alt` option so contest exchange can come
  from two fields.  This is useful for example in a QSO party where `SRX_STRING`
  was only filled out for in-state county contacts and out-of-state values were
  recorded in the `STATE` field.

### Changed

* Make `help` output less overhelming; only show format-specific options when
  asked explicitly like `adifmt help csv`.
* Associate more DXCC entities with their parent ISO 3166-1 codes.  `adifmt fix`
  can now, in most cases, pick the right DXCC entity if `STATE` is set to a
  known primary subdivision code and `COUNTRY` is set to an ISO code, e.g.
  `COUNTRY=USA STATE=AK` sets `COUNTRY=ALASKA DXCC=6`.
* Improve validation of primary subdivision validation.  If ADIF does not define
  subdivisions for a country, `validate` logs a warning if there’s a value in
  the `STATE` field.  If the country has subdivisions defined and the `STATE`
  doesn’t match, it reports an error.
* Release process improvements.

### Fixed

* Cabrillo: preserve frequency precision when converting between MHz and kHz.

### Security

Upgrade golang.org/x/crypto package.


## [v0.1.12-rc1] - 2024-02-03

### Changed

Release process improvements.  `adifmt version` shows git revision.


## [v0.1.11] - 2024-01-30

### Added

* `adifmt validate --required-fields` option to ensure an expected set of fields
  are set on each record.
* `adifmt edit --rename old=new` to change a field’s name.  Doesn’t permit
  overwriting existing data, but cyclical names like
  `--rename name=my_name --rename my_name=name` are supported.


## [v0.1.10] - 2024-01-28

### Fixed

Fix build error on Go 1.18.


## [v0.1.9] - 2024-01-27

### Added

Automated GitHub releases.  Thanks @pcunning!

Support for [Cabrillo 3.0](https://wwrof.org/cabrillo/) contest log format.
Converts Cabrillo to and from ADIF.  Cabrillo header fields can be set with
command-line options or `APP_CABRILLO_` app-specific fields in the ADIF header.
Headers converted from Cabrillo are preserved in ADIF.  Some header fields are
inferred if all QSOs are equal or fit a category.  Make sure to check the values
in a text editor before submitting your log, though.


## [v0.1.8] - 2023-05-15

### Fixed

Fix build error on Go versions before 1.20.


## [v0.1.7] - 2023-05-15

### Added

Conditions, comparisons, and the `find` and `sort` commands.  `find` filters a
log to only records matching a condition.  `sort` orders a log by the specified
`--fields` list.

Commands which take conditionals accept `--if`, `--if-not`, `--or-if`,
`--or-if-not` options to define a boolean condition.  The `--or` options split
groups of AND conditions.  Conditions support string and numeric equality, and
less-than, greater-than, and less-equal/greater-equal comparisons.  Equality can
check several values at once like `--if 'band=40m|20m'`.  “Not equals” is
expressed as `--if-not mode=CW`.  Fields can be compared with other fields using
braces, e.g. `--if 'TIME_ON<={TIME_OFF}'`.


## [v0.1.6] - 2023-04-10

### Added

* Use `adif edit --if field=value` to only apply edits to records matching a
  particular condition.
* Use `{FIELD}` as a template in a filename, e.g.
  `adifmt save '{STATION_CALLSIGN}_{QSO_DATE}'`.

### Changed

Use 1-based record indexing in error and warning messages.


## [v0.1.5] - 2023-03-12

### Added

Allow “ragged” CSV files with different field counts in each line.

### Fixed

Improved error handling.


## [v0.1.4] - 2023-03-12

### Fixed

Ignore case in scoped enum validation.


## [v0.1.3] - 2023-03-11

### Fixed

Properly validate scoped enumerations.  The `STATE` field needs to be a primary
administrative subdivision of the `COUNTRY` in the record, not just a valid
abbreviation in _any_ country.


## [v0.1.2] - 2023-03-07

### Added

Map [ISO 3166-1 alpha codes](https://en.wikipedia.org/wiki/ISO_3166-1) to DXCC
entity names in `adifmt fix`.  This allows you to write `USA` and `ZA` in the
`COUNTRY` field and transform them to `UNITED STATES OF AMERICA` and `REPUBLIC
OF SOUTH AFRICA` which `adifmt infer` can use for the DXCC number.


## [v0.1.1] - 2023-03-04

### Fixed

Adhere to ADIF spec with CRLF line breaks and limit to ASCII if `--adi-ascii-only` option is set.  (Default behavior allows UTF-8 data.)


## [v0.1.0] - 2023-02-27

First numbered version.  Supports `cat`, `edit`, `fix`, `infer`, `select`,
`save`, `help`, and `version` commands.  Supports `ADI`, `ADX`, `CSV`, `TSV`,
and `JSON` formats.

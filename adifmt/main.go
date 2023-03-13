// Copyright 2022 Google LLC
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

// adifmt provides a variety of subcommands for manipulating ADIF log files.
// Run adifmt -help or see https://github.com/flwyd/adif-multitool for
// detailed documentation.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/flwyd/adif-multitool/adif"
	"github.com/flwyd/adif-multitool/adif/spec"
	"github.com/flwyd/adif-multitool/cmd"
)

const (
	helpUrl = "https://github.com/flwyd/adif-multitool"
)

var (
	programName = "adifmt"
	version     = "v0.0.0"
)

func init() {
	if build, ok := debug.ReadBuildInfo(); ok {
		programName = filepath.Base(build.Path)
		version = build.Main.Version
	}
}

func main() {
	ctx := &cmd.Context{}
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.SetOutput(os.Stderr)
	configureContext(ctx, fs)
	if len(os.Args) < 2 {
		fs.Usage = usage(fs, "")
		fs.Usage()
		os.Exit(2)
	}
	name := os.Args[1]
	fs.Usage = usage(fs, name)

	// special case commands, can also be specified as flags
	switch strings.TrimLeft(name, "-") {
	case "help", "h":
		name = ""
		if len(os.Args) > 2 {
			if _, ok := commandNamed(os.Args[2]); ok {
				name = os.Args[2]
			}
		}
		fs.Usage = usage(fs, name)
		// help explicitly requested, so print to stdout and exit without error
		fs.SetOutput(os.Stdout)
		fs.Usage()
		os.Exit(0)
	case "version":
		name = "version"
	}

	c, ok := commandNamed(name)
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown command %q\n", name)
		fmt.Fprintf(os.Stderr, "Usage: %s command [options] [file ...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands are %s\n", strings.Join(commandNames(), ", "))
		fmt.Fprintf(os.Stderr, "Run %s help for more details\n", os.Args[0])
		os.Exit(2)
	}
	if c.Configure != nil {
		c.Configure(ctx, fs)
	}
	fs.Parse(os.Args[2:])
	err := c.Run(ctx, fs.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running %s: %v\n", name, err)
		os.Exit(1)
	}
}

type runeValue struct {
	r *rune
}

func (v runeValue) String() string {
	if v.r == nil {
		return ""
	}
	return fmt.Sprintf("%q", *v.r)
}

func (v runeValue) Set(s string) error {
	switch utf8.RuneCountInString(s) {
	case 0:
		return nil
	case 1:
		r, _ := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError {
			return fmt.Errorf("invalid UTF-8 encoding %q", s)
		}
		*v.r = r
		return nil
	default:
		return fmt.Errorf("expecting one character, not %q", s)
	}
}

func (v runeValue) Get() rune { return *v.r }

func configureContext(ctx *cmd.Context, fs *flag.FlagSet) {
	adiio := adif.NewADIIO()
	adxio := adif.NewADXIO()
	csvio := adif.NewCSVIO()
	jsonio := adif.NewJSONIO()
	tsvio := adif.NewTSVIO()
	ctx.Readers = map[adif.Format]adif.Reader{
		adif.FormatADI: adiio, adif.FormatADX: adxio, adif.FormatCSV: csvio, adif.FormatJSON: jsonio, adif.FormatTSV: tsvio,
	}
	ctx.Writers = map[adif.Format]adif.Writer{
		adif.FormatADI: adiio, adif.FormatADX: adxio, adif.FormatCSV: csvio, adif.FormatJSON: jsonio, adif.FormatTSV: tsvio,
	}
	ctx.Out = os.Stdout
	ctx.Prepare = func(l *adif.Logfile) {
		t := time.Now()
		l.Header.SetComment(fmt.Sprintf("Generated at %s with %d records by %s", t.Format(time.RFC1123Z), len(l.Records), helpUrl))
		l.Header.Set(adif.Field{Name: spec.AdifVerField.Name, Value: spec.ADIFVersion})
		l.Header.Set(adif.Field{Name: spec.CreatedTimestampField.Name, Value: t.Format("20060102 150405")})
		l.Header.Set(adif.Field{Name: spec.ProgramidField.Name, Value: programName})
		l.Header.Set(adif.Field{Name: spec.ProgramversionField.Name, Value: version})
	}

	// General flags
	fmtopts := "options: " + strings.Join(adif.FormatNames(), ", ")
	fs.Var(&ctx.InputFormat, "input",
		"input `format` when it cannot be inferred from file extension\n"+fmtopts)
	fs.Var(&ctx.OutputFormat, "output",
		"output `format` written to stdout\n"+fmtopts)
	fs.Var(&ctx.UserdefFields, "userdef",
		fmt.Sprintf("define a USERDEF `field` name and optional type, range, or enum (multi)\nfield formats: STRING_F:S NUMBER_F{0:360} ENUM_F:{A,B,C}\ntype indicators: %s#Data_Types", spec.ADIFSpecURL))

	// ADI flags
	fs.BoolVar(&adiio.ASCIIOnly, "adi-ascii-only", false,
		"ADI files: error on any non-ASCII characters, instead of writing UTF-8")
	fs.BoolVar(&adiio.LowerCase, "adi-lower-case", false,
		"ADI files: print tags in lower case instead of upper case")
	sepHelp := "options: " + strings.Join(adif.SeparatorNames(), ", ")
	fs.Var(&adiio.FieldSep, "adi-field-separator",
		"ADI files: field `separator`\n"+sepHelp)
	fs.Var(&adiio.RecordSep, "adi-record-separator",
		"ADI files: record `separator`\n"+sepHelp)

	// ADX flags
	fs.IntVar(&adxio.Indent, "adx-indent", 1, "ADX files: indent nested XML structures `n` spaces, 0 for no whitespace")

	// CSV flags
	// TODO csv-lower-case
	// TODO separate comma values for input and output?
	fs.Var(&runeValue{&csvio.Comma}, "csv-field-separator", "CSV files: field separator `character` if not comma")
	fs.Var(&runeValue{&csvio.Comment}, "csv-comment", "CSV files: ignore lines beginnig with `character`")
	fs.BoolVar(&csvio.LazyQuotes, "csv-lazy-quotes", false, "CSV files: be relaxed about quoting rules")
	fs.BoolVar(&csvio.RequireFullRecord, "csv-require-all-fields", false, "CSV files: error if fewer fields in a record than in header")
	fs.BoolVar(&csvio.TrimLeadingSpace, "csv-trim-space", false, "CSV files: ignore leading space in fields")
	fs.BoolVar(&csvio.CRLF, "csv-crlf", false, "CSV files: output MS Windows line endings")

	// JSON flags
	// TODO json-lower-case
	fs.BoolVar(&jsonio.HTMLSafe, "json-html-safe", false, "JSON files: escape characters including < > & for use in HTML")
	fs.IntVar(&jsonio.Indent, "json-indent", 1, "JSON files: indent nested JSON structures `n` spaces, 0 for no whitespace")
	fs.BoolVar(&jsonio.TypedOutput, "json-typed-output", false, "JSON files: output numbers and booleans instead of strings")

	// TSV flags
	fs.BoolVar(&tsvio.CRLF, "tsv-crlf", false, "TSV files: output MS Windows line endings")
	fs.BoolVar(&tsvio.EscapeSpecial, "tsv-escape-special", false, "TSV files: accept and produce \\t \\r \\n and \\\\ escapes in fields")
	fs.BoolVar(&tsvio.IgnoreEmptyHeaders, "tsv-ignore-empty-headers", false, "TSV files: do not return error if a TSV file has an empty header field")
}

func usage(fs *flag.FlagSet, command string) func() {
	return func() {
		out := fs.Output()
		fmt.Fprintln(out, "ADIF Multitool: read and transform ADIF amateur radio logs")
		fmt.Fprintln(out, "See examples at https://github.com/flwyd/adif-multitool")
		fmt.Fprintln(out)
		fmt.Fprintf(out, "Usage: %s command [options] [file ...]\n", fs.Name())
		fmt.Fprintln(out, "Process ADIF files (or - for standard input), write ADIF to standard output.")
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Global options:")
		fs.PrintDefaults()
		fmt.Fprintln(out)
		if c, ok := commandNamed(command); ok {
			fmt.Fprintf(out, "%s: %s\n", c.Name, c.Description)
			cfs := flag.NewFlagSet(command, flag.ContinueOnError)
			if c.Configure != nil {
				c.Configure(&cmd.Context{}, cfs)
			}
			cfs.SetOutput(out)
			cfs.PrintDefaults()
			if c.Help != nil {
				fmt.Fprint(out, c.Help())
			}
		} else {
			fmt.Fprintln(out, "Commands:")
			for _, c := range cmds {
				fmt.Fprintf(out, "  %s: %s\n", c.Name, c.Description)
			}
			fmt.Fprintf(out, "To see options specific to a particular command, run\n%s help command\n", fs.Name())
		}
	}
}

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
	"unicode/utf8"

	"github.com/flwyd/adif-multitool/adif"
	"github.com/flwyd/adif-multitool/adif/spec"
	"github.com/flwyd/adif-multitool/cmd"
)

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

var (
	cmds  = []cmd.Command{cmd.Cat, cmd.Edit, cmd.Fix, cmd.Select, cmd.Validate}
	adiio = adif.NewADIIO()
	csvio = adif.NewCSVIO()
	ctx   = &cmd.Context{
		ADIFVersion: spec.ADIFVersion,
		ProgramName: filepath.Base(os.Args[0]),
		Readers:     map[adif.Format]adif.Reader{adif.FormatADI: adiio, adif.FormatCSV: csvio},
		Writers:     map[adif.Format]adif.Writer{adif.FormatADI: adiio, adif.FormatCSV: csvio},
		Out:         os.Stdout,
	}
	global = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
)

func init() {
	if build, ok := debug.ReadBuildInfo(); !ok {
		ctx.ProgramVersion = "v0.0.0"
	} else {
		ctx.ProgramVersion = build.Main.Version
	}

	// General flags
	fmtopts := "options: " + strings.Join(adif.FormatNames(), ", ")
	global.Var(&ctx.InputFormat, "input",
		"input `format` when it cannot be inferred from file extension\n"+fmtopts)
	global.Var(&ctx.OutputFormat, "output",
		"output `format` written to stdout\n"+fmtopts)

	// ADI flags
	global.BoolVar(&adiio.LowerCase, "adi-lower-case", false,
		"ADI files: print tags in lower case instead of upper case")
	sepHelp := "options: " + strings.Join(adif.SeparatorNames(), ", ")
	global.Var(&adiio.FieldSep, "adi-field-separator",
		"ADI files: field `separator`\n"+sepHelp)
	global.Var(&adiio.RecordSep, "adi-record-separator",
		"ADI files: record `separator`\n"+sepHelp)

	// CSV flags
	// TODO separate comma values for input and output?
	global.Var(&runeValue{&csvio.Comma}, "csv-field-separator", "CSV files: field separator `character` if not comma")
	global.Var(&runeValue{&csvio.Comment}, "csv-comment", "CSV files: ignore lines beginnig with `character`")
	global.BoolVar(&csvio.LazyQuotes, "csv-lazy-quotes", false, "CSV files: be relaxed about quoting rules")
	global.BoolVar(&csvio.TrimLeadingSpace, "csv-trim-space", false, "CSV files: ignore leading space in fields")
	global.BoolVar(&csvio.UseCRLF, "csv-crlf", false, "CSV files: output MS Windows line endings")

	global.Usage = func() {
		out := global.Output()
		fmt.Fprintf(out, "Usage: %s command [flags] files...\n", os.Args[0])
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Commands:")
		for _, c := range cmds {
			fmt.Fprintf(out, "%s: %s\n", c.Name, c.Description)
		}
		fmt.Fprintln(out)
		// TODO this includes command flags, not just global; it would be nice to
		// highlight those
		fmt.Fprintln(out, "Global flags:")
		global.PrintDefaults()
		fmt.Fprintln(out)
		fmt.Fprintf(out, "To see flags specific to a particular command, run\n%s command -help\n", os.Args[0])
		fmt.Fprintln(out, "ADIF Multitool: read and transform ADIF radio logs, output to stdout")
		fmt.Fprintln(out, "See examples at https://github.com/flwyd/adif-multitool")
	}
}

func main() {
	if len(os.Args) < 2 {
		global.Usage()
		os.Exit(2)
	}
	cmd := os.Args[1]
	if cmd == "-help" || cmd == "help" {
		global.Usage()
		os.Exit(2)
	}
	for _, c := range cmds {
		if c.Name == cmd {
			if c.AddFlags != nil {
				c.AddFlags(ctx, global)
			}
			global.Parse(os.Args[2:])
			err := c.Run(ctx, global.Args())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error running %s: %v\n", cmd, err)
				os.Exit(1)
			}
			return
		}
	}
	fmt.Fprintf(global.Output(), "Unknown command %q\n", cmd)
	cmdNames := make([]string, 0, len(cmds))
	for _, cmd := range cmds {
		cmdNames = append(cmdNames, cmd.Name)
	}
	fmt.Fprintf(global.Output(), "Commands are %s\n", strings.Join(cmdNames, ", "))
	fmt.Fprintf(global.Output(), "Run %s -help for more details\n", os.Args[0])
	os.Exit(2)
}

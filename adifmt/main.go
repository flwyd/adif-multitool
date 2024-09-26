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

	"github.com/flwyd/adif-multitool/adif"
	"github.com/flwyd/adif-multitool/adif/spec"
	"github.com/flwyd/adif-multitool/cmd"
	"golang.org/x/exp/slices"
)

const (
	helpUrl = "https://github.com/flwyd/adif-multitool"
)

var (
	programName = "adifmt"
	version     = ""
	vcsRevision = "(unknown)"
)

func init() {
	if b, ok := debug.ReadBuildInfo(); ok {
		programName = filepath.Base(b.Path)
		if version == "" {
			version = b.Main.Version
		}
		for _, s := range b.Settings {
			if s.Key == "vcs.revision" {
				vcsRevision = s.Value
			}
		}
	}
}

func main() {
	os.Exit(runMain(defaultPrepare))
}

func runMain(prepare func(l *adif.Logfile)) int {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.SetOutput(os.Stderr)
	ctx := buildContext(fs, prepare)

	if len(os.Args) < 2 {
		fs.Usage = usage(fs, "")
		fs.Usage()
		return 2
	}
	name := os.Args[1]
	fs.Usage = usage(fs, name)

	// special case commands, can also be specified as flags
	switch strings.TrimLeft(name, "-") {
	case "help", "h":
		name = ""
		if len(os.Args) > 2 && !strings.HasPrefix(os.Args[2], "-") {
			name = os.Args[2]
		}
		fs.Usage = usage(fs, name)
		// help explicitly requested, so print to stdout and exit without error
		fs.SetOutput(os.Stdout)
		fs.Usage()
		return 0
	case "version":
		name = "version"
	}

	// Add format flags late so help is less overwhelming
	for _, f := range formatConfigs {
		f.AddFlags(fs)
	}

	c, ok := commandNamed(name)
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown command %q\n", name)
		fmt.Fprintf(os.Stderr, "Usage: %s command [options] [file ...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands are %s\n", strings.Join(commandNames(), ", "))
		fmt.Fprintf(os.Stderr, "Run %s help for more details\n", os.Args[0])
		return 2
	}
	if c.Configure != nil {
		c.Configure(ctx, fs)
	}
	// filenames can come before or after flags, but not interspersed
	args := os.Args[2:]
	firstflag := slices.IndexFunc(args, func(s string) bool { return strings.HasPrefix(s, "-") })
	if firstflag < 0 {
		firstflag = len(args)
	}
	nonflags := args[0:firstflag]
	args = args[firstflag:]
	fs.Parse(args)
	nonflags = append(nonflags, fs.Args()...)
	err := c.Run(ctx, nonflags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running %s: %v\n", name, err)
		return 1
	}
	return 0
}

func defaultPrepare(l *adif.Logfile) {
	t := time.Now()
	l.Header.SetComment(fmt.Sprintf("Generated at %s with %d records by %s", t.Format(time.RFC1123Z), len(l.Records), helpUrl))
	l.Header.Set(adif.Field{Name: spec.AdifVerField.Name, Value: spec.ADIFVersion})
	l.Header.Set(adif.Field{Name: spec.CreatedTimestampField.Name, Value: t.Format("20060102 150405")})
	l.Header.Set(adif.Field{Name: spec.ProgramidField.Name, Value: programName})
	l.Header.Set(adif.Field{Name: spec.ProgramversionField.Name, Value: version})
}

func buildContext(fs *flag.FlagSet, prepare func(l *adif.Logfile)) *cmd.Context {
	ctx := &cmd.Context{
		Out:     os.Stdout,
		Readers: make(map[adif.Format]adif.Reader),
		Writers: make(map[adif.Format]adif.Writer),
		Prepare: prepare,
	}
	for _, f := range formatConfigs {
		ctx.Readers[f.Format()] = f.IO()
		ctx.Writers[f.Format()] = f.IO()
	}

	// General flags
	fmtopts := "options: " + strings.Join(adif.FormatNames(), ", ")
	fs.Var(&ctx.FieldOrder, "field-order", "Comma-separated `field` order for output (repeatable)")
	fs.Var(&ctx.InputFormat, "input",
		"input `format` when it cannot be inferred from file extension\n"+fmtopts)
	fs.Var(&ctx.OutputFormat, "output",
		"output `format` written to stdout\n"+fmtopts)
	fs.Var(&languageValue{Tag: &ctx.Locale}, "locale",
		"BCP-47 `language` code for IntlString comparisons e.g. da, pt-BR, zh-Hant")
	fs.BoolVar(&ctx.SuppressAppHeaders, "suppress-app-headers", false,
		"Don't output app-defined headers, to comply with ADIF 3.1.4 spec")
	fs.Var(&ctx.UserdefFields, "userdef",
		fmt.Sprintf("define a USERDEF `field` name and optional type, range, or enum (multi)\nfield formats: STRING_F:S NUMBER_F,{0:360} ENUM_F,{A,B,C}\ntype indicators: %s#Data_Types", spec.ADIFSpecURL))
	return ctx
}

func usage(fs *flag.FlagSet, term string) func() {
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
		if c, ok := commandNamed(term); ok {
			fmt.Fprintf(out, "%s: %s\n", c.Name, c.Description)
			cfs := flag.NewFlagSet(term, flag.ContinueOnError)
			if c.Configure != nil {
				c.Configure(&cmd.Context{}, cfs)
			}
			cfs.SetOutput(out)
			cfs.PrintDefaults()
			if c.Help != nil {
				fmt.Fprint(out, c.Help())
			}
		} else if c := formatNamed(term); c != nil {
			fmt.Fprintf(out, "%s format:\n", c.Format())
			cfs := flag.NewFlagSet(term, flag.ContinueOnError)
			c.AddFlags(cfs)
			cfs.SetOutput(out)
			cfs.PrintDefaults()
			fmt.Fprint(out, c.Help())
		} else {
			fmt.Fprintln(out, "Formats:", strings.Join(adif.FormatNames(), ", "))
			fmt.Fprintf(out, "To see options specific to a format, run\n%s help formatname\n", fs.Name())
			fmt.Fprintln(out)
			fmt.Fprintln(out, "Commands:")
			for _, c := range cmds {
				fmt.Fprintf(out, "  %s: %s\n", c.Name, c.Description)
			}
			fmt.Fprintf(out, "To see options specific to a particular command, run\n%s help command\n", fs.Name())
		}
	}
}

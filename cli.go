package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		s string
		t string

		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.StringVar(&s, "s", "", "source language. ex. en, ja ...")
	flags.StringVar(&t, "t", "", "target language. ex. ja, en ...")
	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	apiParams, err := cli.parseFlags(s, t)
	if err != nil {
		fmt.Fprintln(cli.errStream, err.Error())
		return ExitCodeError
	}

	// Run
	err = cli.run(apiParams)
	if err != nil {
		fmt.Fprintln(cli.errStream, err.Error())
		return ExitCodeError
	}

	return ExitCodeOK
}

func (cli *CLI) parseFlags(s, t string) (*APIParams, error) {
	var (
		source = LanguageEn
		target = LanguageJa
	)

	if s != "" {
		if langMap[s] {
			source = s
		} else {
			return nil, fmt.Errorf("invalid source: %s", s)
		}
	}
	if t != "" {
		if langMap[t] {
			target = t
		} else {
			return nil, fmt.Errorf("invalid target: %s", t)
		}
	}

	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}

	if len(b) == 0 {
		return nil, fmt.Errorf("require text")
	}

	return &APIParams{
		Text:   string(b),
		Source: source,
		Target: target,
	}, nil
}

func (cli *CLI) run(p *APIParams) error {
	result, err := NewAPICaller(p).Call()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(cli.outStream, result)
	if err != nil {
		return err
	}

	return nil
}

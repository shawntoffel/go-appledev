package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"time"
)

type Flags struct {
	// Version subcommand
	Version bool
	// Token subcommand
	Token bool
	// Config subcommand
	GenerateConfig bool

	PrivateKeyFilePath string
	KeyID              string
	TeamID             string
	ServiceID          string
	Duration           time.Duration

	// ConfigFilePath is the path to a JSON file from which config may be parsed.
	ConfigFilePath string

	// OutputFile is the destination file path where a config file may be written.
	OutputFile string
}

func ParseFlags(args []string) (*Flags, error) {
	flags := &Flags{}

	subcommand := args[1]
	switch subcommand {
	case "token":
		flags.Token = true
		err := flags.parseTokenFlags(args)
		if err != nil {
			return nil, err
		}
	case "config":
		flags.GenerateConfig = true
		err := flags.parseConfigGenFlags(args)
		if err != nil {
			return nil, err
		}
	case "version":
		flags.Version = true
	}

	return flags, nil
}

func (f *Flags) parseTokenFlags(args []string) error {
	outputBuffer := bytes.NewBuffer([]byte{})

	set := flag.NewFlagSet("token", flag.ContinueOnError)
	set.SetOutput(outputBuffer)
	set.StringVar(&f.ConfigFilePath, "c", f.ConfigFilePath, "Path to a json config file containing args.")
	f.addCommonArgs(set)

	err := set.Parse(args[2:])
	if err != nil {
		return fmt.Errorf(outputBuffer.String())
	}

	if len(f.PrivateKeyFilePath) < 1 && len(f.ConfigFilePath) < 1 {
		fmt.Fprintf(set.Output(), "%s token: A Private key file path is required when a config file is not provided.\n\n", args[0])
		fmt.Fprintf(set.Output(), "Usage of %s token:\n", args[0])
		set.PrintDefaults()
	}

	errContent := outputBuffer.String()
	if len(errContent) > 0 {
		return fmt.Errorf(errContent)
	}

	return nil
}

func (f *Flags) parseConfigGenFlags(args []string) error {
	outputBuffer := bytes.NewBuffer([]byte{})

	set := flag.NewFlagSet("config", flag.ContinueOnError)
	set.SetOutput(outputBuffer)
	set.StringVar(&f.OutputFile, "o", f.OutputFile, "(required) The path where the config file will be written.")
	f.addCommonArgs(set)

	err := set.Parse(args[2:])
	if err != nil {
		return fmt.Errorf(outputBuffer.String())
	}

	errContent := outputBuffer.String()
	if len(errContent) > 0 {
		return fmt.Errorf(errContent)
	}

	return nil
}

func (f *Flags) addCommonArgs(flatSet *flag.FlagSet) {
	flatSet.StringVar(&f.PrivateKeyFilePath, "pk", f.PrivateKeyFilePath, "(required) The path to a file containing your PEM encoded private key.")
	flatSet.StringVar(&f.KeyID, "kid", f.KeyID, "The Key ID associated with your private key.")
	flatSet.StringVar(&f.TeamID, "tid", f.TeamID, "The Team ID from your developer account.")
	flatSet.StringVar(&f.ServiceID, "sid", f.ServiceID, "The Service ID from your developer account.")
	flatSet.DurationVar(&f.Duration, "d", f.Duration, "How long the token will be valid for.")
}

func exitHelp() {
	fmt.Printf("Usage of %s:\n", os.Args[0])
	fmt.Println("  token    \tCreate an apple developer token.")
	fmt.Println("  config    \tGenerate a config file.")
	fmt.Println("  version    \tPrint the version of this application.")
	os.Exit(1)
}

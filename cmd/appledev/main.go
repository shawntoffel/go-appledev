package main

import (
	"fmt"
	"os"
)

const app = "appledev"

var version = "dev"

func init() {
	if len(os.Args) < 2 {
		exitHelp()
	}
}

func main() {
	flags, err := ParseFlags(os.Args)
	if err != nil {
		exit(err.Error(), 1)
	}

	if flags.Version {
		exit(version, 0)
	}

	config, err := NewConfigFromFlags(flags)
	if err != nil {
		exit(err.Error(), 1)
	}

	if flags.Token {
		token, err := config.CreateToken()
		if err != nil {
			exit(err.Error(), 1)
		}

		exit(token, 0)
	}

	if flags.GenerateConfig {
		err = config.WriteToFile(flags.OutputFile)
		if err != nil {
			exit(err.Error(), 1)
		}

		os.Exit(0)
	}

	exitHelp()
}

func exit(content string, code int) {
	if code == 1 {
		fmt.Fprintln(os.Stderr, app+":", content)
	} else {
		fmt.Println(content)
	}
	os.Exit(code)
}

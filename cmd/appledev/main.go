package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/shawntoffel/go-appledev"
)

var version = "dev"

var (
	flagVersion = false
	flagToken   = false
)

var (
	flagTokenPKFile    = ""
	flagTokenKeyID     = ""
	flagTokenTeamID    = ""
	flagTokenServiceID = ""
	flagTokenDuration  = time.Minute * 30
)

func init() {
	if len(os.Args) < 2 {
		exitHelp()
	}

	subcommand := os.Args[1]
	flagVersion = subcommand == "version"
	if subcommand == "token" {
		flagToken = true
		parseTokenFlags()
	}
}

func main() {
	if flagVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if flagToken {
		token, err := createToken()
		if err != nil {
			fmt.Println("appledev:", err)
			os.Exit(1)
		}

		fmt.Println(token)
		os.Exit(0)
	}
}

func exitHelp() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Println("  token    \tCreate an apple developer token.")
	fmt.Println("  version    \tPrint the version of this application.")
	os.Exit(1)
}

func parseTokenFlags() {
	tokenCmd := flag.NewFlagSet("token", flag.ExitOnError)
	tokenCmd.StringVar(&flagTokenPKFile, "pk", flagTokenPKFile, "(required) The path to a file containing your PEM encoded private key.")
	tokenCmd.StringVar(&flagTokenKeyID, "kid", flagTokenKeyID, "(required) The Key ID associated with your private key.")
	tokenCmd.StringVar(&flagTokenTeamID, "tid", flagTokenTeamID, "(required) The Team ID from your developer account.")
	tokenCmd.StringVar(&flagTokenServiceID, "sid", flagTokenServiceID, "(required) The Service ID from your developer account.")
	tokenCmd.DurationVar(&flagTokenDuration, "d", flagTokenDuration, "How long the token will be valid for.")
	tokenCmd.Parse(os.Args[2:])

	if len(flagTokenPKFile) < 1 {
		fmt.Fprintf(os.Stderr, "%s token: Private key file path is required.\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage of %s token:\n", os.Args[0])
		tokenCmd.PrintDefaults()
		os.Exit(1)
	}
}

func createToken() (string, error) {
	tokenProvider := appledev.ApiTokenProvider{
		KeyID:     flagTokenKeyID,
		TeamID:    flagTokenTeamID,
		ServiceID: flagTokenServiceID,
		Duration:  flagTokenDuration,
	}

	bytes, err := os.ReadFile(flagTokenPKFile)
	if err != nil {
		return "", err
	}

	token, err := tokenProvider.SignedJWT(bytes)
	if err != nil {
		return "", err
	}

	return token, nil
}

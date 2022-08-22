# go-appledev

[![Go Reference](https://pkg.go.dev/badge/github.com/shawntoffel/go-appledev.svg)](https://pkg.go.dev/github.com/shawntoffel/go-appledev) 
 [![Go Report Card](https://goreportcard.com/badge/github.com/shawntoffel/go-appledev)](https://goreportcard.com/report/github.com/shawntoffel/go-appledev) [![Build status](https://github.com/shawntoffel/go-appledev/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/shawntoffel/go-appledev/actions/workflows/go.yml)

A library & command line application for generating the signed developer tokens used to authenticate against Apple REST APIs. A valid signed JWT token is constructed using your private key, key ID, team ID, and service ID. 

*go-appledev* is an open source project not affiliated with Apple Inc.

### Locating your identifiers:
* **Key ID**: An identifier associated with your private key. It can be found on the [Certificates, Identifiers & Profiles](https://developer.apple.com/account/resources/authkeys/list) page under Keys. Click on the appropriate key to view the ID. 
* **Team ID**: Found on the [membership](https://developer.apple.com/account/#!/membership) page. 
* **Service ID**: Found on the [Certificates, Identifiers & Profiles](https://developer.apple.com/account/resources/identifiers/list/serviceId) page under Identifiers. Make sure "Services IDs" is selected from the dropdown. 

### Use as a library:
```go

import (
  "github.com/shawntoffel/go-appledev"
 )

// Initialize the API token provider.
tokenProvider := appledev.ApiTokenProvider{
  KeyID:     "keyId",
  TeamID:    "teamId",
  ServiceID: "serviceId",
  Duration:  time.Minute * 30,
}

// Fetch your private key bytes. 
bytes, err := os.ReadFile(privateKeyFilePath)

// Generate a signed JWT string.
token, err := tokenProvider.SignedJWT(bytes)

```

### Use as a binary application:
Precompiled binaries are available on the [Releases](https://github.com/shawntoffel/go-appledev/releases) page. 
```sh

Usage of ./appledev:
  token         Create an apple developer token.
  version       Print the version of this application.

Usage of ./appledev token:
  -d duration
        How long the token will be valid for. (default 30m0s)
  -kid string
        (required) The Key ID associated with your private key.
  -pk string
        (required) The path to a file containing your PEM encoded private key.
  -sid string
        (required) The Service ID from your developer account.
  -tid string
        (required) The Team ID from your developer account.
```

```sh
./appledev token -pk AuthKey_ABCDE12345.p8 -kid keyId -sid serviceId -tid teamId
```

## Troubleshooting
Please use the GitHub [Discussions](https://github.com/shawntoffel/go-appledev/discussions) tab for questions regarding this client library.

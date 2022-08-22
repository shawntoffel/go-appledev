package appledev

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TokenProvider generates authentication tokens for apple REST APIs.
type ApiTokenProvider struct {
	// KeyID is the Key identifier from your developer accoung.
	KeyID string

	// TeamID is the Team ID from your developer accoung.
	TeamID string

	// ServiceID is the Service ID from your developer accoung.
	ServiceID string

	// Duration is how long the token will be valid for.
	Duration time.Duration
}

// GenerateSignedJWT generates a valid JWT signed with your PEM private key.
func (g *ApiTokenProvider) SignedJWT(pemBytes []byte) (string, error) {
	token, err := g.create()
	if err != nil {
		return "", err
	}

	return g.sign(token, pemBytes)
}

func (g *ApiTokenProvider) create() (*jwt.Token, error) {
	err := g.validate()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	return &jwt.Token{
		Header: map[string]interface{}{
			"alg": jwt.SigningMethodES256.Alg(),
			"kid": g.KeyID,
			"id":  g.TeamID + "." + g.ServiceID,
		},
		Claims: jwt.MapClaims{
			"iss": g.TeamID,
			"sub": g.ServiceID,
			"iat": now.Unix(),
			"exp": now.Add(g.Duration).Unix(),
		},
		Method: jwt.SigningMethodES256,
	}, nil
}

func (g *ApiTokenProvider) sign(token *jwt.Token, pemBytes []byte) (string, error) {
	privateKey, err := jwt.ParseECPrivateKeyFromPEM(pemBytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key. %s", err)
	}

	signedJwt, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to create signed JWT. %s", err)
	}

	return signedJwt, nil
}

func (g *ApiTokenProvider) validate() error {
	messages := []string{}

	if len(g.KeyID) < 1 {
		messages = append(messages, "key identifier may not be empty")
	}

	if len(g.TeamID) < 1 {
		messages = append(messages, "team ID may not be empty")
	}

	if len(g.ServiceID) < 1 {
		messages = append(messages, "service ID may not be empty")
	}

	if g.Duration <= 0 {
		messages = append(messages, "token expiration must be in the future (duration must be greater than 0)")
	}

	if len(messages) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(messages, ", "))
	}

	return nil
}

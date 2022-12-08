package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/shawntoffel/go-appledev"
)

const DefaultDuration string = "30m"

type Config struct {
	KeyID      string `json:"kid"`
	TeamID     string `json:"tid"`
	ServiceID  string `json:"sid"`
	Duration   string `json:"d"`
	PrivateKey string `json:"pk"`
}

func NewConfigFromFlags(flags *Flags) (*Config, error) {
	config, err := NewConfig(flags.ConfigFilePath)
	if err != nil {
		return nil, err
	}

	err = config.applyFlags(flags)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func NewConfig(filePath string) (*Config, error) {
	if len(filePath) < 1 {
		return &Config{
			Duration: DefaultDuration,
		}, nil
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	config := &Config{}

	err = json.Unmarshal(bytes, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) applyFlags(flags *Flags) error {
	if flags.KeyID != "" {
		c.KeyID = flags.KeyID
	}

	if flags.ServiceID != "" {
		c.ServiceID = flags.ServiceID
	}

	if flags.TeamID != "" {
		c.TeamID = flags.TeamID
	}

	if flags.Duration > 0 {
		c.Duration = flags.Duration.String()
	}

	if flags.PrivateKeyFilePath != "" {
		bytes, err := os.ReadFile(flags.PrivateKeyFilePath)
		if err != nil {
			return err
		}

		c.PrivateKey = string(bytes)
	}

	return nil
}

func (c *Config) CreateToken() (string, error) {
	d := c.Duration
	if len(d) < 1 {
		d = DefaultDuration
	}

	duration, err := time.ParseDuration(d)
	if err != nil {
		return "", err
	}

	tokenProvider := &appledev.ApiTokenProvider{
		KeyID:     c.KeyID,
		TeamID:    c.TeamID,
		ServiceID: c.ServiceID,
		Duration:  duration,
	}

	token, err := tokenProvider.SignedJWT([]byte(c.PrivateKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (c *Config) WriteToFile(filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer f.Close()

	content, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(content)
	return err
}

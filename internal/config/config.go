package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	Active string            `json:"active"`
	Tokens map[string]string `json:"tokens"`
	Models map[string]string `json:"models"`
}

func New() *Config {
	return &Config{
		Tokens: map[string]string{},
		Models: map[string]string{},
	}
}

func DefaultPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "lexicon", "config.json"), nil
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return New(), nil
	}
	if err != nil {
		return nil, err
	}
	cfg := New()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	if cfg.Tokens == nil {
		cfg.Tokens = map[string]string{}
	}
	if cfg.Models == nil {
		cfg.Models = map[string]string{}
	}
	return cfg, nil
}

func (c *Config) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return err
	}
	return os.Chmod(path, 0o600)
}

func (c *Config) Token(providerID string) string {
	return c.Tokens[providerID]
}

func (c *Config) SetToken(providerID, token string) {
	c.Tokens[providerID] = token
}

func (c *Config) Model(providerID, fallback string) string {
	if m := c.Models[providerID]; m != "" {
		return m
	}
	return fallback
}

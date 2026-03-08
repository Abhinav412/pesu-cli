package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Token  string `json:"token"`
	ApiURL string `json:"api_url"`
}

const configFileName = ".pesu-cli.json"

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configFileName), nil
}

func LoadConfig() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		// Return default config if file doesn't exist
		return &Config{ApiURL: "http://localhost:8080/api"}, nil
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.ApiURL == "" {
		cfg.ApiURL = "http://localhost:8080/api"
	}

	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

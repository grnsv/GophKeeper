package config

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	ServerAddress string `toml:"server_address"`
}

func getConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appConfigDir := filepath.Join(configDir, "GophKeeper")
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(appConfigDir, "config.toml"), nil
}

func saveConfig(path string, cfg *Config) error {
	data, err := toml.Marshal(*cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func Parse() (*Config, error) {
	configFile, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		ServerAddress: "http://localhost:8080",
	}

	data, err := os.ReadFile(configFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if len(data) > 0 {
		if err = toml.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}
	if err = saveConfig(configFile, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

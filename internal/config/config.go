package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	EnvTokenKey    = "PRODUCTBOARD_API_TOKEN"
	ConfigDir      = ".config/pboard"
	ConfigFileName = "config"
	ConfigFileType = "yaml"
	DefaultBaseURL = "https://api.productboard.com"
)

// Config holds application configuration.
type Config struct {
	APIToken string
	BaseURL  string
}

// Load reads configuration from env var and config file.
// Env var PRODUCTBOARD_API_TOKEN takes precedence over config file.
func Load() (*Config, error) {
	cfg := &Config{
		BaseURL: DefaultBaseURL,
	}

	// Check environment variable first (highest precedence)
	if token := os.Getenv(EnvTokenKey); token != "" {
		cfg.APIToken = token
		return cfg, nil
	}

	// Try config file
	home, err := os.UserHomeDir()
	if err != nil {
		return cfg, nil
	}

	configPath := filepath.Join(home, ConfigDir)
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileType)
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		// Config file not found is not an error — token may be set later
		return cfg, nil
	}

	cfg.APIToken = viper.GetString("api_token")
	if url := viper.GetString("api_url"); url != "" {
		cfg.BaseURL = url
	}

	return cfg, nil
}

// ConfigFilePath returns the full path to the config file.
func ConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ConfigDir, ConfigFileName+".yaml"), nil
}

// WriteToken writes the API token to the config file with mode 600.
func WriteToken(token string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot determine home directory: %w", err)
	}

	configPath := filepath.Join(home, ConfigDir)
	if err := os.MkdirAll(configPath, 0700); err != nil {
		return fmt.Errorf("cannot create config directory: %w", err)
	}

	filePath := filepath.Join(configPath, ConfigFileName+".yaml")
	content := fmt.Sprintf("api_token: %q\napi_url: %q\n", token, DefaultBaseURL)

	if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
		return fmt.Errorf("cannot write config file: %w", err)
	}

	return nil
}

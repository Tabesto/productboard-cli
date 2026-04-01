package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	EnvTokenKey       = "PRODUCTBOARD_API_TOKEN"
	EnvAPIVersionKey  = "PRODUCTBOARD_API_VERSION"
	ConfigDir         = ".config/pboard"
	ConfigFileName    = "config"
	ConfigFileType    = "yaml"
	DefaultBaseURL    = "https://api.productboard.com"
	DefaultAPIVersion = "2"
)

// Config holds application configuration.
type Config struct {
	APIToken   string
	BaseURL    string
	APIVersion string
}

// Load reads configuration from config file and env vars.
// Env vars take precedence over config file values.
func Load() (*Config, error) {
	cfg := &Config{
		BaseURL:    DefaultBaseURL,
		APIVersion: DefaultAPIVersion,
	}

	// Try config file first (lowest precedence)
	home, err := os.UserHomeDir()
	if err == nil {
		configPath := filepath.Join(home, ConfigDir)
		viper.SetConfigName(ConfigFileName)
		viper.SetConfigType(ConfigFileType)
		viper.AddConfigPath(configPath)

		if err := viper.ReadInConfig(); err == nil {
			cfg.APIToken = viper.GetString("api_token")
			if url := viper.GetString("api_url"); url != "" {
				cfg.BaseURL = url
			}
			if ver := viper.GetString("api_version"); ver != "" {
				cfg.APIVersion = ver
			}
		}
	}

	// Env vars override config file
	if token := os.Getenv(EnvTokenKey); token != "" {
		cfg.APIToken = token
	}
	if ver := os.Getenv(EnvAPIVersionKey); ver != "" {
		cfg.APIVersion = ver
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

// WriteConfig writes the API token and version to the config file with mode 600.
func WriteConfig(token, apiVersion string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot determine home directory: %w", err)
	}

	configPath := filepath.Join(home, ConfigDir)
	if err := os.MkdirAll(configPath, 0700); err != nil {
		return fmt.Errorf("cannot create config directory: %w", err)
	}

	filePath := filepath.Join(configPath, ConfigFileName+".yaml")
	if apiVersion == "" {
		apiVersion = DefaultAPIVersion
	}
	content := fmt.Sprintf("api_token: %q\napi_url: %q\napi_version: %q\n", token, DefaultBaseURL, apiVersion)

	if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
		return fmt.Errorf("cannot write config file: %w", err)
	}

	return nil
}

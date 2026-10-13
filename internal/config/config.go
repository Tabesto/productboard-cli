package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	EnvTokenKey        = "PRODUCTBOARD_API_TOKEN"
	EnvAPIVersionKey   = "PRODUCTBOARD_API_VERSION"
	EnvWriteEnabledKey = "PRODUCTBOARD_WRITE_ENABLED"
	ConfigDir          = ".config/pboard"
	ConfigFileName     = "config"
	ConfigFileType     = "yaml"
	DefaultBaseURL     = "https://api.productboard.com"
	DefaultAPIVersion  = "2"
)

// Config holds application configuration.
type Config struct {
	APIToken   string
	BaseURL    string
	APIVersion string
	// WriteEnabled is the write gate: false (default) means create/update calls are refused.
	WriteEnabled bool
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
			cfg.WriteEnabled = viper.GetBool("write_enabled")
		}
	}

	// Env vars override config file
	if token := os.Getenv(EnvTokenKey); token != "" {
		cfg.APIToken = token
	}
	if ver := os.Getenv(EnvAPIVersionKey); ver != "" {
		cfg.APIVersion = ver
	}
	// Only "1"/"true" enable; any other non-empty value forces writes off.
	if w := os.Getenv(EnvWriteEnabledKey); w != "" {
		cfg.WriteEnabled = w == "1" || strings.EqualFold(w, "true")
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

// Update merges values into the config file (created with mode 600 if absent),
// preserving every key not listed. api_url and api_version get defaults only when missing.
func Update(values map[string]any) error {
	filePath, err := ConfigFilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(filePath), 0700); err != nil {
		return fmt.Errorf("cannot create config directory: %w", err)
	}

	v := viper.New() // fresh instance: don't leak values into the global one used by Load
	v.SetConfigFile(filePath)
	v.SetConfigType(ConfigFileType)
	v.SetConfigPermissions(0600)
	if err := v.ReadInConfig(); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("cannot read config file: %w", err)
	}
	v.SetDefault("api_url", DefaultBaseURL)
	v.SetDefault("api_version", DefaultAPIVersion)
	for k, val := range values {
		v.Set(k, val)
	}

	if err := v.WriteConfigAs(filePath); err != nil {
		return fmt.Errorf("cannot write config file: %w", err)
	}
	return os.Chmod(filePath, 0600)
}

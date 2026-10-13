package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func setupHome(t *testing.T, yaml string) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv(EnvWriteEnabledKey, "")
	t.Setenv(EnvTokenKey, "")
	viper.Reset()
	t.Cleanup(viper.Reset)
	if yaml != "" {
		dir := filepath.Join(home, ConfigDir)
		if err := os.MkdirAll(dir, 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(yaml), 0600); err != nil {
			t.Fatal(err)
		}
	}
	return home
}

func TestLoadWriteEnabledEnv(t *testing.T) {
	cases := []struct {
		env, file string
		want      bool
	}{
		{"", "", false},
		{"1", "", true},
		{"true", "", true},
		{"TRUE", "", true},
		{"0", "", false},
		{"false", "", false},
		{"yes", "", false},
		{"", "write_enabled: true\n", true},
		{"0", "write_enabled: true\n", false}, // env overrides config
		{"1", "write_enabled: false\n", true},
	}
	for _, tc := range cases {
		setupHome(t, tc.file)
		t.Setenv(EnvWriteEnabledKey, tc.env)
		cfg, err := Load()
		if err != nil {
			t.Fatal(err)
		}
		if cfg.WriteEnabled != tc.want {
			t.Errorf("env=%q file=%q: WriteEnabled=%v, want %v", tc.env, tc.file, cfg.WriteEnabled, tc.want)
		}
	}
}

func TestUpdatePreservesKeys(t *testing.T) {
	setupHome(t, "api_token: old\napi_url: https://example.test\napi_version: \"1\"\nwrite_enabled: true\n")

	if err := Update(map[string]any{"api_token": "new"}); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.APIToken != "new" || cfg.BaseURL != "https://example.test" || cfg.APIVersion != "1" || !cfg.WriteEnabled {
		t.Errorf("settings not preserved: %+v", cfg)
	}
	path, _ := ConfigFilePath()
	if fi, _ := os.Stat(path); fi.Mode().Perm() != 0600 {
		t.Errorf("mode = %v, want 0600", fi.Mode().Perm())
	}
}

func TestUpdateCreatesWithDefaults(t *testing.T) {
	setupHome(t, "")
	if err := Update(map[string]any{"write_enabled": true}); err != nil {
		t.Fatal(err)
	}
	cfg, _ := Load()
	if !cfg.WriteEnabled || cfg.BaseURL != DefaultBaseURL || cfg.APIVersion != DefaultAPIVersion {
		t.Errorf("cfg = %+v", cfg)
	}
}

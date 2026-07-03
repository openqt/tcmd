package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Settings holds user preferences.
type Settings struct {
	ShowHidden   bool   `yaml:"show_hidden"`
	SortBy       string `yaml:"sort_by"`
	ConfirmDelete bool  `yaml:"confirm_delete"`
	Editor       string `yaml:"editor"`
}

// DefaultSettings returns factory defaults.
func DefaultSettings() Settings {
	return Settings{SortBy: "name", ConfirmDelete: true, Editor: "vi"}
}

// SettingsPath returns config file path.
func SettingsPath() (string, error) {
	dir, err := EnsureDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "settings.yaml"), nil
}

// LoadSettings reads settings.yaml.
func LoadSettings() (Settings, error) {
	path, err := SettingsPath()
	if err != nil {
		return DefaultSettings(), err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultSettings(), nil
		}
		return DefaultSettings(), err
	}
	var s Settings
	if err := yaml.Unmarshal(data, &s); err != nil {
		return DefaultSettings(), err
	}
	return s, nil
}

// SaveSettings writes settings.yaml.
func SaveSettings(s Settings) error {
	path, err := SettingsPath()
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

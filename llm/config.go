package llm

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Configuration struct {
	Provider    string `json:"provider,omitempty"`
	APIKey      string `json:"api_key,omitempty"`
	Model       string `json:"model,omitempty"`
	BaseURL     string `json:"base_url,omitempty"`
	Personality string `json:"personality,omitempty"`
}

func ConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("locate user config dir: %w", err)
	}
	return filepath.Join(dir, "nova", "config.json"), nil
}

func LoadConfig() (*Configuration, error) {
	cfg := &Configuration{}

	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	switch {
	case err == nil:
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
	case os.IsNotExist(err):
	default:
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	if v := os.Getenv("NOVA_PROVIDER"); v != "" {
		cfg.Provider = v
	}
	if v := os.Getenv("NOVA_API_KEY"); v != "" {
		cfg.APIKey = v
	}
	if v := os.Getenv("NOVA_MODEL"); v != "" {
		cfg.Model = v
	}
	if v := os.Getenv("NOVA_BASE_URL"); v != "" {
		cfg.BaseURL = v
	}
	if v := os.Getenv("NOVA_PERSONALITY"); v != "" {
		cfg.Personality = v
	}

	return cfg, nil
}

func SaveConfig(cfg *Configuration) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func (c *Configuration) IsConfigured() bool {
	if c.Model == "" {
		return false
	}
	if c.APIKey == "" {
		return c.BaseURL != ""
	}
	return true
}

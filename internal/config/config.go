package config

import (
	"fmt"
	"os"
	"os/user"

	"github.com/BurntSushi/toml"
)

// Config holds the global CLI config.
type Config struct {
	BaseURL string `toml:"base_url"`
	Token   string `toml:"token"`
}

// path returns the default config path.
func path() string {
	usr, err := user.Current()
	if err != nil {
		return "~/.config/memos-cli/config.toml"
	}
	return usr.HomeDir + "/.config/memos-cli/config.toml"
}

// Load reads the config file. Returns nil if the file does not exist.
func Load() (*Config, error) {
	p := path()
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}

// Save writes the config to the default path.
func Save(cfg *Config) error {
	p := path()

	dir := p[:len(p)-len("/config.toml")]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data, err := toml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	return os.WriteFile(p, data, 0600)
}

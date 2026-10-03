package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/vinegm/tuday/internal/doomsday"
)

// validate reports an error for any setting
// the rest of the program can't work with.
func (config Config) validate() error {
	years := doomsday.Years{Start: config.StartYear, End: config.EndYear}
	if err := years.Validate(); err != nil {
		return fmt.Errorf("start_year/end_year: %w", err)
	}

	return nil
}

// Loads configuration from file. The returned config is always usable:
// on error the defaults are returned alongside it.
func Load(path string) (Config, error) {
	config := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return Default(), err
	}

	if err := config.validate(); err != nil {
		return Default(), err
	}

	return config, nil
}

func DefaultPath() string {
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".config", "dday", "config.yaml")
	}

	// crappy fallback, might not even work
	return "$HOME/.config/dday/config.yaml"
}

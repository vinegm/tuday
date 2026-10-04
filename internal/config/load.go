package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/vinegm/tuday/internal/storage/paths"
)

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

	return config, nil
}

// DefaultPath is the config file's path.
func DefaultPath() string {
	return paths.Config("config.yaml")
}

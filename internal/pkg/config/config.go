package config

import "proteinreminder/internal/pkg/log"

var (
	// Use this interface for managing config.
	config Config
)

// Get config values used in the app.
type Config interface {
	Get(name string, defaultValue string) string
}

// Get config value.
func Get(name string, defaultValue string) string {
	if config == nil {
		log.Debug("config is null")
		return ""
	}
	return config.Get(name, defaultValue)
}

// Set config used "config.Get".
func SetConfig(c Config) {
	config = c
}

package config

import (
	"slacktimer/internal/pkg/log"
)

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
		log.Error("config is null")
		return ""
	}
	return config.Get(name, defaultValue)
}

// MustGet is like Get but panics if the value is empty.
func MustGet(name string) string {
	v := config.Get(name, "")
	if v == "" {
		panic("error MustGet")
	}
	return v
}

// Set config used "config.Get".
func SetConfig(c Config) {
	config = c
}

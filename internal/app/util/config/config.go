// Package config provides config values used in the app.
package config

import (
	"fmt"
	"slacktimer/internal/app/util/log"
)

var (
	// Use this interface for managing config.
	config Config
)

// Config defines methods that returns config values used in the app.
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
		panic(fmt.Sprintf("error MustGet name=%s", name))
	}
	return v
}

// SetConfig sets config used "config.Get".
func SetConfig(c Config) {
	config = c
}

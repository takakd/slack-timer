// Package driver provides the implementation of config methods.
package driver

import (
	"fmt"
	"os"
	"slacktimer/internal/app/util/config"
	"strings"

	"github.com/joho/godotenv"
)

// EnvConfig provides implementation of config.Config based on .env file.
type EnvConfig struct {
}

var _ config.Config = (*EnvConfig)(nil)

// NewEnvConfig creates new struct.
func NewEnvConfig(filepathList ...string) *EnvConfig {
	if len(filepathList) > 0 {
		if err := godotenv.Load(filepathList...); err != nil {
			panic(fmt.Sprintf("failed to load .env files. %v", filepathList))
		}
	}
	return &EnvConfig{}
}

// Get returns value corresponding name.
func (e EnvConfig) Get(name string, defaultValue string) string {
	v := os.Getenv(strings.ToUpper(name))
	if v == "" {
		v = defaultValue
	}
	return v
}

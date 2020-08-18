package config

import (
	"os"
	"strings"
)

type Config interface {
	Get(name string) string
}

//
// Environment Value Config
//

type EnvConfig struct {
}

func NewEnvConfig() Config {
	return &EnvConfig{}
}

func (e *EnvConfig) Get(name string) string {
	return os.Getenv(strings.ToUpper(name))
}

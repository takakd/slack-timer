// Deprecated
package config

import (
	"github.com/joho/godotenv"
	"os"
	"proteinreminder/internal/errorutil"
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

var envConfig *EnvConfig = nil

func (e *EnvConfig) Get(name string) string {
	return os.Getenv(strings.ToUpper(name))
}

func GetEnvConfig() *EnvConfig {
	if envConfig != nil {
		return envConfig
	}

	if err := godotenv.Load(); err != nil {
		panic(errorutil.MakePanicMessage("failed to load .env."))
	}

	return &EnvConfig{}
}

// @deprecated
func NewEnvConfig() Config {
	return &EnvConfig{}
}

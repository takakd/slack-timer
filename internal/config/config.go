// Deprecated
package config

import (
	"os"
	"strings"
	"github.com/joho/godotenv"
	"proteinreminder/internal/errorutil"
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


	return &EnvConfig{
	}
}

// @deprecated
func NewEnvConfig() Config {
	return &EnvConfig{}
}

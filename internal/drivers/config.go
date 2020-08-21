package drivers

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

var envConfig *EnvConfig = nil

func (e *EnvConfig) Get(name string) string {
	return os.Getenv(strings.ToUpper(name))
}

func GetEnvConfig() *EnvConfig {
	if envConfig != nil {
		return envConfig
	}

	// NOTE: Not need, cause using commandline mode.
	//if err := godotenv.Load(); err != nil {
	//	panic(errorutil.MakePanicMessage(err))
	//}

	return &EnvConfig{}
}

// @deprecated
func NewEnvConfig() Config {
	return &EnvConfig{}
}

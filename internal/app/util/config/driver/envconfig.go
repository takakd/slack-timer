// Deprecated
package driver

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"slacktimer/internal/pkg/helper"
	"strings"
)

type EnvConfig struct {
}

func NewEnvConfig(filepathList ...string) *EnvConfig {
	if len(filepathList) > 0 {
		if err := godotenv.Load(filepathList...); err != nil {
			panic(helper.MakePanicMessage(fmt.Sprintf("failed to load .env files. %v", filepathList)))
		}
	}
	return &EnvConfig{}
}

func (e *EnvConfig) Get(name string, defaultValue string) string {
	v := os.Getenv(strings.ToUpper(name))
	if v == "" {
		v = defaultValue
	}
	return v
}

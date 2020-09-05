// Deprecated
package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"proteinreminder/internal/pkg/errorutil"
	"strings"
)

type EnvConfig struct {
}

func NewEnvConfig(filepathList ...string) *EnvConfig {
	if len(filepathList) > 0 {
		if err := godotenv.Load(filepathList...); err != nil {
			panic(errorutil.MakePanicMessage(fmt.Sprintf("failed to load .env files. %v", filepathList)))
		}
	}
	return &EnvConfig{}
}

func (e *EnvConfig) Get(name string) string {
	return os.Getenv(strings.ToUpper(name))
}

// Deprecated
package config

import (
	"path/filepath"
	"proteinreminder/internal/pkg/fileutil"
	"proteinreminder/internal/pkg/panicutil"
)

// Get config values used in the app.
type Config interface {
	Get(name string) string
}

// Use this interface for managing config.
var config Config

func init() {
	config = GetConfig("", envPath())
}

func envPath() string {
	var err error
	appDir, err := fileutil.GetAppDir()
	if err != nil {
		panic(panicutil.MakePanicMessage("need app directory path."))
	}
	path := filepath.Join(appDir, ".env")
	if !fileutil.FileExists(path) {
		path = ""
	}
	return path
}

// Get config implementation.
func GetConfig(name string, params ...interface{}) Config {
	var c Config = nil
	if name == "" {
		var names []string
		if params != nil {
			for _, p := range params {
				if name, ok := p.(string); ok && name != "" {
					names = append(names, name)
				}
			}
		}
		c = NewEnvConfig(names...)
	}
	return c
}

// TODO: deprecated because global function complicates unit tests.
// Get config value.
func Get(name string) string {
	return config.Get(name)
}

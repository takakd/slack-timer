// Deprecated
package config

import (
	"path/filepath"
	"proteinreminder/internal/pkg/errorutil"
	"proteinreminder/internal/pkg/fileutil"
)

const (
	defaultConfig string = "env"
)

// Get config values used in the app.
type Config interface {
	Get(name string, defaultValue string) string
}

// Use this interface for managing config.
var config Config

func init() {
	config = GetConfig("")
}

func envPath() string {
	var err error
	appDir, err := fileutil.GetAppDir()
	if err != nil {
		panic(errorutil.MakePanicMessage("need app directory path."))
	}
	path := filepath.Join(appDir, ".env")
	if !fileutil.FileExists(path) {
		path = ""
	}
	return path
}

// Get config implementation.
// Returns config corresponding to name and params.
// Currently, name supports only "env".
func GetConfig(name string, params ...interface{}) Config {

	if name == "" {
		name = defaultConfig
	}

	if name == defaultConfig {
		var names []string
		path := envPath()
		if path != "" {
			names = append(names, path)
		}
		// Receive .env path by params.
		if params != nil {
			for _, p := range params {
				if name, ok := p.(string); ok && name != "" {
					names = append(names, name)
				}
			}
		}
		return NewEnvConfig(names...)
	}

	return nil
}

// Deprecated
///*
//Get config value.
//This is utility function.
// */
//func Get(name string) string {
//	return config.Get(name)
//}

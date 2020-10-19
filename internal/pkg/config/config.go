package config

import (
	"os"
	"path/filepath"
	"proteinreminder/internal/pkg/errorutil"
	"proteinreminder/internal/pkg/fileutil"
)

const (
	defaultConfig = "env"
)

var (
	// Use this interface for managing config.
	config Config
)

// Get config values used in the app.
type Config interface {
	Get(name string, defaultValue string) string
}

// Get config value.
func Get(name string, defaultValue string) string {
	return config.Get(name, defaultValue)
}

func init() {
	configType := os.Getenv("APP_CONFIG_TYPE")
	if configType == "" {
		configType = defaultConfig
	}

	if configType == defaultConfig {
		var names []string
		path := envPath()
		if path != "" {
			names = append(names, path)
		}
		SetConfig(NewEnvConfig(names...))
	}
}

// Set config used "config.Get".
func SetConfig(c Config) {
	config = c
}

// Returns .env file path in app directory.
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

//// Get config implementation.
//// Returns config corresponding to name and params.
//// Currently, name supports only "env".
//func GetConfig(name string, params ...interface{}) Config {
//
//	if name == "" {
//		name = defaultConfig
//	}
//
//	if name == defaultConfig {
//		var names []string
//		path := envPath()
//		if path != "" {
//			names = append(names, path)
//		}
//		// Receive .env path by params.
//		if params != nil {
//			for _, p := range params {
//				if name, ok := p.(string); ok && name != "" {
//					names = append(names, name)
//				}
//			}
//		}
//		return NewEnvConfig(names...)
//	}
//
//	return nil
//}

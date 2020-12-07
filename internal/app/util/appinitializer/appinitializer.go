// Package appinitializer set up the base packages in the app.
package appinitializer

import (
	"os"
	"path/filepath"
	"slacktimer/internal/app/util/config"
	"slacktimer/internal/app/util/config/driver"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/di/container/dev"
	"slacktimer/internal/app/util/di/container/prod"
	"slacktimer/internal/app/util/di/container/test"
	"slacktimer/internal/app/util/log"
	driver2 "slacktimer/internal/app/util/log/driver"
	"slacktimer/internal/pkg/helper"
)

// AppInit calls when the app launch.
func AppInit() {
	setConfig()
	setDi()
	setLogger()
}

// Setup config.
func setConfig() {
	configType := os.Getenv("APP_CONFIG_TYPE")
	if configType == "" {
		configType = "env"
	}

	if configType == "env" {
		// Get .env path
		appDir, err := helper.GetAppDir()
		if err != nil {
			panic("need app directory path.")
		}
		names := make([]string, 0)
		path := filepath.Join(appDir, ".env")
		if helper.FileExists(path) {
			names = append(names, path)
		}
		config.SetConfig(driver.NewEnvConfig(names...))
	}
}

// Setup DI container.
func setDi() {
	env := config.Get("APP_ENV", "test")

	if env == "prod" {
		di.SetDi(&prod.Container{})
	} else if env == "dev" {
		di.SetDi(&dev.Container{})
	} else if env == "test" {
		di.SetDi(&test.Container{})
	}
}

// Setup logger.
func setLogger() {
	log.SetDefaultLogger(driver2.NewCloudWatchLogger())
	log.SetLevel(config.Get("APP_LOG_LEVEL", "debug"))
}

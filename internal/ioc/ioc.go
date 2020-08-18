package ioc

import (
	"proteinreminder/internal/config"
	"proteinreminder/internal/log"
)

// Contain each struct.
type Container struct {
	logger log.Logger
	config config.Config
}

// Caching
var container *Container

func getContainer() *Container {
	if container == nil {
		container = &Container{}
	}
	return container
}

//
// Struct managed IOC are below.
//

func GetLogger() log.Logger {
	c := getContainer()
	if c.logger == nil {
		c.logger = log.NewStdoutLogger(log.LevelDebug)
	}
	return c.logger
}

func GetConfig() config.Config {
	c := getContainer()
	if c.config == nil {
		c.config = config.NewEnvConfig()
	}
	return c.config
}

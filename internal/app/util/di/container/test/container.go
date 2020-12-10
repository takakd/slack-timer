package test

import (
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
)

// Container implements DI on production env.
type Container struct {
}

var _ di.DI = (*Container)(nil)

// Get returns interfaces corresponding name.
func (t *Container) Get(name string) interface{} {
	log.Info("di.Get", name)
	return nil
}

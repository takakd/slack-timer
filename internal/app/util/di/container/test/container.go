package test

import (
	"fmt"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
)

// Container implements DI on production env.
type Container struct {
}

var _ di.DI = (*Container)(nil)

// Get returns interfaces corresponding name.
func (t *Container) Get(name string) interface{} {
	log.Info(fmt.Sprintf("call di.Get name=%s", name))
	return nil
}

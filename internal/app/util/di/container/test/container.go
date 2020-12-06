package test

import (
	"fmt"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
)

type Container struct {
}

var _ di.DI = (*Container)(nil)

// Returns interfaces in test environment.
func (t *Container) Get(name string) interface{} {
	log.Info(fmt.Sprintf("call di.Get name=%s", name))
	return nil
}

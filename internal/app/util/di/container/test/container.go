package test

import (
	"fmt"
	"slacktimer/internal/app/util/log"
)

type Container struct {
}

// Returns interfaces in test environment.
func (t *Container) Get(name string) interface{} {
	log.Debug(fmt.Sprintf("call di.Get name=%s", name))
	return nil
}

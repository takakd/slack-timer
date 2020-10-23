package container

import (
	"fmt"
	"proteinreminder/internal/pkg/log"
)

type Test struct {
}

// Returns interfaces in test environment.
func (t *Test) Get(name string) interface{} {
	log.Debug(fmt.Sprintf("call di.Get name=%s", name))
	return nil
}

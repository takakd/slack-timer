package di

import "slacktimer/internal/app/util/log"

var (
	di DI
)

// DI returns interfaces implemented per environment.
type DI interface {
	Get(name string) interface{}
}

// Helper function of DI.Get
func Get(name string) interface{} {
	if di == nil {
		log.Error("di is null")
		return nil
	}
	return di.Get(name)
}

func SetDi(d DI) {
	di = d
}

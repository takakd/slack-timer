// Package di provides dependency injection container.
package di

var (
	di DI
)

// DI returns interfaces implemented per environment.
type DI interface {
	Get(name string) interface{}
}

// Get is helper function of DI.Get
func Get(name string) interface{} {
	if di == nil {
		return nil
	}
	return di.Get(name)
}

// SetDi sets DI used throughout the application.
func SetDi(d DI) {
	di = d
}

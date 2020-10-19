package di

import "proteinreminder/internal/pkg/config"

var (
	di DI
)

// DI returns interfaces implemented per environment.
type DI interface {
	Get(name string) interface{}
}

// Helper function of DI.Get
func Get(name string) interface{} {
	return di.Get(name)
}

func SetDi(d DI) {
	di = d
}

func setDi() {
	// Set DI interface per environment.
	switch appEnv := config.Get("APP_ENV", "dev"); appEnv {
	case "production":
		SetDi(nil)
	case "test":
		SetDi(&TestDi{})
	default:
		SetDi(nil)
	}
}

func init() {
	setDi()
}

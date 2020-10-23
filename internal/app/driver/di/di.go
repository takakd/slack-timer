package di

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

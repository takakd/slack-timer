package di

type TestDi struct {
}

func (t *TestDi) Get(name string) interface{} {
	var i interface{}
	switch name {
	default:
		i = nil
	}
	return i
}

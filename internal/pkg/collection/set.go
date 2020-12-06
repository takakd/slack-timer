// Package collection provides collection data structure.
package collection

const (
	// Use as the stuffing data.
	// Use "map" for implementing this struct but it does not use the value of map,
	// so set dummy value to the map value.
	Stuffing = 0
)

type Set struct {
	values map[interface{}]int
}

func NewSet() *Set {
	set := &Set{}
	set.values = make(map[interface{}]int)
	return set
}

func (s *Set) Set(value interface{}) {
	s.values[value] = Stuffing
}

func (s *Set) Remove(value interface{}) {
	if s.Contains(value) {
		delete(s.values, value)
	}
}

// Check if the value contains.
func (s Set) Contains(value interface{}) bool {
	_, contain := s.values[value]
	return contain
}

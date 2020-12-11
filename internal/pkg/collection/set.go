// Package collection provides collection data structure.
package collection

const (
	// Use as the stuffing data.
	// Use "map" for implementing this struct but it does not use the value of map,
	// so set dummy value to the map value.
	_stuffing = 0
)

// Set provides "set" structure.
type Set struct {
	values map[interface{}]int
}

// NewSet creates new struct.
func NewSet() *Set {
	set := &Set{}
	set.values = make(map[interface{}]int)
	return set
}

// Set sets value.
func (s *Set) Set(value interface{}) {
	s.values[value] = _stuffing
}

// Remove removes value.
func (s *Set) Remove(value interface{}) {
	if s.Contains(value) {
		delete(s.values, value)
	}
}

// Contains check if the value contains.
func (s Set) Contains(value interface{}) bool {
	_, contain := s.values[value]
	return contain
}

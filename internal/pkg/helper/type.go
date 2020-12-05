package helper

import "reflect"

// Check if the type is struct.
func IsStruct(v interface{}) bool {
	c := reflect.ValueOf(v)

	// value
	if c.Kind() == reflect.Struct {
		return true
	}

	// pointer
	if c.Kind() == reflect.Ptr {
		if c.Elem().Type().Kind() == reflect.Struct {
			return true
		}
	}

	return false
}

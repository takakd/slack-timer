package typeutil

import "reflect"

// Check if the type is struct.
func IsStruct(v interface{}) bool {
	return reflect.ValueOf(v).Kind() == reflect.Struct
}

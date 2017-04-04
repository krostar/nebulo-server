package tools

import "reflect"

// IsZeroOrNil return true if x is the zero value of the its type
func IsZeroOrNil(x interface{}) bool {
	if x.(reflect.Value).Kind() == reflect.Ptr {
		return x.(reflect.Value).IsNil()
	}
	return x == reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

// Package first return the first element.
package first

import "reflect"

func isZero(v reflect.Value) (out bool) {
	switch v.Kind() {
	case reflect.Slice:
		out = v.Len() == 0
	default:
		zi := reflect.Zero(v.Type()).Interface()
		out = reflect.DeepEqual(zi, v.Interface())
	}

	return
}

// One return the first non zero value.
// If all value is zero, return the last type of input's default value.
func One(items ...interface{}) interface{} {
	if len(items) < 1 {
		panic("first: should give at least one argument")
	}

	var v reflect.Value

	for _, item := range items {
		v = reflect.ValueOf(item)

		if !isZero(v) {
			return item
		}
	}

	return reflect.Zero(v.Type()).Interface()
}

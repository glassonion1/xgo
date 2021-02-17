package xgo

import (
	"reflect"
	"strings"
)

func IsFirstUpper(v string) bool {
	if len(v) == 0 {
		return false
	}
	return strings.HasPrefix(v, strings.ToUpper(string(v[0])))
}

func IsBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.IsZero()
	case reflect.Bool:
		return value.IsZero()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.IsZero()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.IsZero()
	case reflect.Float32, reflect.Float64:
		return value.IsZero()
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}

	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

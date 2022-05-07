package xgo

import "reflect"

// Contains returns true if an element is present in a slice
func Contains[T any](list []T, elem T) bool {
	defer func() {
		recover()
	}()

	for _, item := range list {
		if ok := reflect.DeepEqual(item, elem); ok {
			return true
		}
	}
	return false

	/*
		listV := reflect.ValueOf(list)

		if listV.Kind() == reflect.Slice {
			for i := 0; i < listV.Len(); i++ {
				item := listV.Index(i).Interface()

				target := reflect.ValueOf(elem).Convert(reflect.TypeOf(item)).Interface()

				if ok := reflect.DeepEqual(item, target); ok {
					return true
				}
			}
		}
		return false
	*/
}

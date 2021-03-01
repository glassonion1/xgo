package xgo

import "reflect"

// Contains returns true if an element is present in a slice
func Contains(list interface{}, elem interface{}) bool {
	defer func() {
		recover()
	}()

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
}

package xgo

import "reflect"

// Contains returns true if an element is present in a slice
func Contains[T any](list []T, elem T) bool {
	for _, item := range list {
		if ok := reflect.DeepEqual(item, elem); ok {
			return true
		}
	}
	return false
}

// IndexOf returns index of slice
func IndexOf[T any](list []T, elem T) int {
	for i, item := range list {
		if ok := reflect.DeepEqual(item, elem); ok {
			return i
		}
	}
	return -1
}

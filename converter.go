package xgo

import "reflect"

// StructToMap converts a struct to map
func StructToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	v := reflect.ValueOf(data)
	t := v.Type()

	size := t.NumField()
	for i := 0; i < size; i++ {
		name := t.Field(i).Name
		if !IsBlank(v.FieldByName(name)) {
			result[name] = v.FieldByName(name).Interface()
		}
	}
	return result
}

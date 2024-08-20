package xgo

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

// tagCopier is tag for deep copy target
const tagCopier = "copier"

type SetCustomField func(src, dst reflect.Value) (bool, error)

// DeepCopy
func DeepCopy(srcModel interface{}, dstModel interface{}) error {
	f := func(src, dst reflect.Value) (bool, error) { return false, nil }
	return DeepCopyWithCustomSetter(srcModel, dstModel, f)
}

// DeepCopy with custom setter
func DeepCopyWithCustomSetter(
	srcModel interface{},
	dstModel interface{},
	customSetter SetCustomField,
) error {

	src := reflect.Indirect(reflect.ValueOf(srcModel))
	dst := reflect.Indirect(reflect.ValueOf(dstModel))

	if !dst.CanAddr() {
		return errors.New("copy to value is unaddressable")
	}

	if src.Kind() == reflect.Slice {
		if err := copySlice(src, dst); err != nil {
			return fmt.Errorf("%v", err)
		}
		return nil
	}

	// What to do if the deepcopy destination model has a tag
	var srcToDstTagMap = map[string]string{}
	for i := 0; i < dst.NumField(); i++ {
		dstF := dst.Type().Field(i)
		if tag, ok := dstF.Tag.Lookup(tagCopier); ok {
			srcToDstTagMap[tag] = dstF.Name
		}
	}

	for i := 0; i < src.NumField(); i++ {
		field := src.Type().Field(i)
		srcFieldValue := src.FieldByName(field.Name)

		dstFieldName := field.Name
		if tag, ok := field.Tag.Lookup(tagCopier); ok {
			dstFieldName = tag
		}
		if tag, ok := srcToDstTagMap[field.Name]; ok {
			dstFieldName = tag
		}

		_, ok := dst.Type().FieldByName(dstFieldName)
		dstFieldValue := dst.FieldByName(dstFieldName)
		if !ok {
			continue
		}

		// Ignores private field
		if !IsFirstUpper(dstFieldName) {
			continue
		}

		// string, int, float
		if srcFieldValue.Type().ConvertibleTo(dstFieldValue.Type()) {
			dstFieldValue.Set(srcFieldValue.Convert(dstFieldValue.Type()))
			continue
		}

		// *string, *int, *float to string, int, float
		if srcFieldValue.Type().Kind() == reflect.Ptr {
			if srcFieldValue.IsNil() {
				continue
			}
			if srcFieldValue.Type().Elem().ConvertibleTo(dstFieldValue.Type()) {
				dstFieldValue.Set(srcFieldValue.Elem().Convert(dstFieldValue.Type()))
				continue
			}
		}

		// string, int, float to *string, *int, *float
		if dstFieldValue.Type().Kind() == reflect.Ptr {
			if srcFieldValue.Type().ConvertibleTo(dstFieldValue.Type().Elem()) {
				rv := reflect.New(dstFieldValue.Type().Elem())
				rv.Elem().Set(srcFieldValue.Convert(dstFieldValue.Type().Elem()))
				dstFieldValue.Set(rv)
				continue
			}
		}

		isSet, err := customSetter(srcFieldValue, dstFieldValue)
		if err != nil {
			return fmt.Errorf("%s: %v", field.Name, err)
		}
		if isSet {
			continue
		}

		// set the time.Time field
		isSet, err = setTimeField(srcFieldValue, dstFieldValue)
		if err != nil {
			return fmt.Errorf("%s: %v", field.Name, err)
		}
		if isSet {
			continue
		}

		// struct, pointer, slice
		switch srcFieldValue.Kind() {
		case reflect.Struct:
			if !field.Anonymous {
				dv, vFunc := instantiate(dstFieldValue)
				if err := DeepCopy(srcFieldValue.Interface(), dv.Interface()); err != nil {
					return fmt.Errorf("%s: %v", field.Name, err)
				}
				dstFieldValue.Set(vFunc())
				continue
			}
			dstFieldValue.SetInt(srcFieldValue.Int())
		case reflect.Ptr:

			if srcFieldValue.IsNil() {
				continue
			}
			// copy to indirect
			indirect := reflect.Indirect(srcFieldValue)
			if indirect.Type().AssignableTo(dstFieldValue.Type()) && dstFieldValue.Type().Kind() != reflect.Ptr {
				dstFieldValue.Set(indirect)
				continue
			}

			dv, vFunc := instantiate(dstFieldValue)
			if err := DeepCopy(srcFieldValue.Interface(), dv.Interface()); err != nil {
				return fmt.Errorf("%s: %v", field.Name, err)
			}
			dstFieldValue.Set(vFunc())
			continue
		case reflect.Slice:
			if err := copySlice(srcFieldValue, dstFieldValue); err != nil {
				return fmt.Errorf("%s: %v", field.Name, err)
			}
			continue
		}

	}

	return nil
}

// Instantiates a value that can handle copying in both directions - from a pointer to a struct and from a struct to a pointer.
func instantiate(v reflect.Value) (reflect.Value, func() reflect.Value) {
	// ptr
	if v.Type().Kind() == reflect.Ptr {
		rv := reflect.New(v.Type().Elem())
		vFunc := func() reflect.Value { return rv }
		return rv, vFunc
	}

	// struct
	rv := reflect.New(v.Type())
	vFunc := func() reflect.Value { return reflect.Indirect(rv) }
	return rv, vFunc
}

func copySlice(src, dst reflect.Value) error {
	if src.IsNil() {
		return nil
	}
	slice := reflect.MakeSlice(reflect.SliceOf(dst.Type().Elem()), src.Len(), src.Cap())
	dst.Set(slice)

	for i := 0; i < src.Len(); i++ {
		d := dst.Index(i)

		// Other than pointer and struct
		if src.Index(i).Type().ConvertibleTo(d.Type()) {
			d.Set(src.Index(i).Convert(d.Type()))
			continue
		}

		// pointer or struct
		dv, vFunc := instantiate(d)
		if err := DeepCopy(src.Index(i).Interface(), dv.Interface()); err != nil {
			return err
		}
		d.Set(vFunc())
	}
	return nil
}

func setTimeField(src, dst reflect.Value) (bool, error) {
	switch t := src.Interface().(type) {
	case time.Time:
		// time.Time -> *timestampps.Timestamp or int64
		switch dst.Interface().(type) {
		case int64:
			dst.Set(reflect.ValueOf(t.Unix()))
			return true, nil
		}

	case *time.Time:
		if t == nil {
			return false, nil
		}
		// *time.Time -> timestampps.Timestamp or int64
		switch dst.Interface().(type) {
		case int64:
			dst.Set(reflect.ValueOf(t.Unix()))
			return true, nil
		}

	case int64:
		// int64 -> time.Time or *timestamppb.Timestamp or *durationpb.Duration
		switch dst.Interface().(type) {
		case time.Time:
			dst.Set(reflect.ValueOf(time.Unix(t, 0)))
			return true, nil
		}
	}
	return false, nil
}

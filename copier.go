package xgo

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// tagCopier is tag for deep copy target
const tagCopier = "copier"

// DeepCopy deepcopy a struct to struct.
func DeepCopy(srcModel interface{}, dstModel interface{}) error {
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
		srcFieldType, ok := src.Type().FieldByName(field.Name)
		srcFieldValue := src.FieldByName(field.Name)
		if !ok {
			continue
		}

		dstFieldName := field.Name
		if tag, ok := field.Tag.Lookup(tagCopier); ok {
			dstFieldName = tag
		}
		if tag, ok := srcToDstTagMap[field.Name]; ok {
			dstFieldName = tag
		}

		dstFieldType, ok := dst.Type().FieldByName(dstFieldName)
		dstFieldValue := dst.FieldByName(dstFieldName)
		if !ok {
			continue
		}

		// Ignores private field
		if !IsFirstUpper(dstFieldName) {
			continue
		}

		if srcFieldType.Type.ConvertibleTo(dstFieldType.Type) {
			dstFieldValue.Set(srcFieldValue.Convert(dstFieldType.Type))
			continue
		}

		isSet, err := setTimeField(srcFieldValue, dstFieldValue)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		if isSet {
			continue
		}

		// from struct to pointer
		switch srcFieldValue.Kind() {
		case reflect.Int64, reflect.Int32:
			dstFieldValue.SetInt(srcFieldValue.Int())
			continue
		case reflect.Struct:
			if !field.Anonymous {
				dv, vFunc := instantiate(dstFieldValue)
				if err := DeepCopy(srcFieldValue.Interface(), dv.Interface()); err != nil {
					return fmt.Errorf("%v", err)
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
			if indirect.Type().AssignableTo(dstFieldType.Type) && dstFieldType.Type.Kind() != reflect.Ptr {
				dstFieldValue.Set(indirect)
				continue
			}
			dv, vFunc := instantiate(dstFieldValue)
			if err := DeepCopy(srcFieldValue.Interface(), dv.Interface()); err != nil {
				return fmt.Errorf("%v", err)
			}
			dstFieldValue.Set(vFunc())
			continue

		case reflect.Slice:
			if err := copySlice(srcFieldValue, dstFieldValue); err != nil {
				return fmt.Errorf("%v", err)
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
		case *timestamppb.Timestamp:
			if t.Unix() > 0 {
				ts := timestamppb.New(t)
				if err := ts.CheckValid(); err != nil {
					return false, err
				}
				dst.Set(reflect.ValueOf(ts))
			}
		case int64:
			dst.Set(reflect.ValueOf(t.Unix()))
		case *time.Time:
			dst.Set(reflect.ValueOf(&t))
		}
		return true, nil
	case *time.Time:
		if src.IsNil() {
			dst.Set(reflect.Zero(dst.Type()))
			return true, nil
		}
		// *time.Time -> timestampps.Timestamp or int64
		switch dst.Interface().(type) {
		case *timestamppb.Timestamp:
			if t.Unix() > 0 {
				ts := timestamppb.New(*t)
				if err := ts.CheckValid(); err != nil {
					return false, err
				}
				dst.Set(reflect.ValueOf(ts))
			}
		case int64:
			dst.Set(reflect.ValueOf(t.Unix()))
		case time.Time:
			dst.Set(reflect.ValueOf(*t))
		}
		return true, nil
	case *timestamppb.Timestamp:
		// *timestamppb.Timestamp -> time.Time
		switch dst.Interface().(type) {
		case time.Time:
			if t.GetSeconds() > 0 {
				ts := timestamppb.New(t.AsTime())
				if err := ts.CheckValid(); err != nil {
					return false, err
				}
				dst.Set(reflect.ValueOf(ts.AsTime()))
			}
		case *time.Time:
			if t.GetSeconds() > 0 {
				ts := timestamppb.New(t.AsTime())
				if err := ts.CheckValid(); err != nil {
					return false, err
				}
				time := ts.AsTime()
				dst.Set(reflect.ValueOf(&time))
			}
		}
		return true, nil

	case time.Duration:
		// time.Duration -> *durationpb.Duration
		switch dst.Interface().(type) {
		case *durationpb.Duration:
			if t.Nanoseconds() > 0 {
				dst.Set(reflect.ValueOf(durationpb.New(t)))
			}
		}
		return true,
			nil
	case *durationpb.Duration:
		// *durationpb.Duration -> time.Duration
		switch dst.Interface().(type) {
		case time.Duration:
			if t == nil {
				return false, nil
			}

			d := durationpb.New(t.AsDuration())
			if err := d.CheckValid(); err != nil {
				return false, err
			}
			dst.Set(reflect.ValueOf(d.AsDuration()))
		}
		return true, nil
	case int64:
		// int64 -> time.Time or *timestamppb.Timestamp or *durationpb.Duration
		switch dst.Interface().(type) {
		case time.Time:
			dst.Set(reflect.ValueOf(time.Unix(t, 0)))
		case *timestamppb.Timestamp:
			if t > 0 {
				ts := timestamppb.New(time.Unix(t, 0))
				if err := ts.CheckValid(); err != nil {
					return false, err
				}
				dst.Set(reflect.ValueOf(ts))
			}
		case *durationpb.Duration:
			if t > 0 {
				dst.Set(reflect.ValueOf(durationpb.New(time.Duration(t))))
			}
		}
		return true, nil
	}
	return false, nil
}

package xgo

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
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
			return err
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
				// struct to struct
				s := reflect.New(dstFieldType.Type)
				v := func() reflect.Value { return reflect.Indirect(s) }
				if dstFieldType.Type.Kind() == reflect.Ptr {
					// struct to ptr
					s, v = reflect.New(dstFieldType.Type.Elem()), func() reflect.Value { return s }
				}

				if err := DeepCopy(srcFieldValue.Interface(), s.Interface()); err != nil {
					return fmt.Errorf("%v", err)
				}
				dstFieldValue.Set(v())
				continue
			}
			dstFieldValue.SetInt(srcFieldValue.Int())

		case reflect.Ptr:
			if !srcFieldValue.IsNil() {
				// copy to indirect
				indirect := reflect.Indirect(srcFieldValue)
				if indirect.Type().AssignableTo(dstFieldType.Type) && dstFieldType.Type.Kind() != reflect.Ptr {
					dstFieldValue.Set(indirect)
					continue
				}

				// ptr to struct
				s := reflect.New(dstFieldType.Type)
				v := func() reflect.Value { return reflect.Indirect(s) }
				if dstFieldType.Type.Kind() == reflect.Ptr {
					// ptr to ptr
					s, v = reflect.New(dstFieldType.Type.Elem()), func() reflect.Value { return s }
				}

				if err := DeepCopy(srcFieldValue.Interface(), s.Interface()); err != nil {
					return fmt.Errorf("%v", err)
				}
				dstFieldValue.Set(v())
				continue
			}
		}
	}

	return nil
}

func setTimeField(src, dst reflect.Value) (bool, error) {

	switch t := src.Interface().(type) {
	case time.Time:
		// time.Time -> *timestamp.Timestamp or int64
		switch dst.Interface().(type) {
		case *timestamp.Timestamp:
			if t.Unix() > 0 {
				tp, err := ptypes.TimestampProto(t)
				if err != nil {
					return false, err
				}
				dst.Set(reflect.ValueOf(tp))
			}
		case int64:
			dst.Set(reflect.ValueOf(t.Unix()))
		}
		return true, nil
	case *time.Time:
		if src.IsNil() {
			dst.Set(reflect.Zero(dst.Type()))
			return true, nil
		}
		tv := *t
		// *time.Time -> timestamp.Timestamp or int64
		switch dst.Interface().(type) {
		case *timestamp.Timestamp:
			if tv.Unix() > 0 {
				tp, err := ptypes.TimestampProto(tv)
				if err != nil {
					return false, err
				}
				dst.Set(reflect.ValueOf(tp))
			}
		case int64:
			dst.Set(reflect.ValueOf(t.Unix()))
		}
		return true, nil
	case *timestamp.Timestamp:
		// *timestamp.Timestamp -> time.Time
		switch dst.Interface().(type) {
		case time.Time:
			if t.GetSeconds() > 0 {
				ti, err := ptypes.Timestamp(t)
				if err != nil {
					return false, err
				}
				dst.Set(reflect.ValueOf(ti))
			}
		case *time.Time:
			if t.GetSeconds() > 0 {
				ti, err := ptypes.Timestamp(t)
				if err != nil {
					return false, err
				}
				dst.Set(reflect.ValueOf(&ti))
			}
		}
		return true, nil

	case time.Duration:
		// time.Duration -> *duration.Duration
		switch dst.Interface().(type) {
		case *duration.Duration:
			if t.Nanoseconds() > 0 {
				dst.Set(reflect.ValueOf(ptypes.DurationProto(t)))
			}
		}
		return true,
			nil
	case *duration.Duration:
		// *duration.Duration -> time.Duration
		switch dst.Interface().(type) {
		case time.Duration:
			if t == nil {
				return false, nil
			}
			d, err := ptypes.Duration(t)
			if err != nil {
				return false, err
			}
			dst.Set(reflect.ValueOf(d))
		}
		return true, nil
	case int64:
		// int64 -> time.Time or *timestamp.Timestamp or *duration.Duration
		switch dst.Interface().(type) {
		case time.Time:
			dst.Set(reflect.ValueOf(time.Unix(t, 0)))
		case *timestamp.Timestamp:
			if t > 0 {
				tp, err := ptypes.TimestampProto(time.Unix(t, 0))
				if err != nil {
					return false, err
				}
				dst.Set(reflect.ValueOf(tp))
			}
		case *duration.Duration:
			if t > 0 {
				dst.Set(reflect.ValueOf(ptypes.DurationProto(time.Duration(t))))
			}
		}
		return true, nil
	}
	return false, nil
}

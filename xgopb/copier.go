package xgopb

import (
	"reflect"
	"time"

	"github.com/glassonion1/xgo"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DeepCopy
func DeepCopy(srcModel interface{}, dstModel interface{}) error {
	return xgo.DeepCopyWithCustomSetter(srcModel, dstModel, setTimeField)
}

func setTimeField(src, dst reflect.Value) (bool, error) {

	switch t := src.Interface().(type) {
	case time.Time:
		if t.Unix() <= 0 {
			return false, nil
		}
		// time.Time -> *timestampps.Timestamp
		switch dst.Interface().(type) {
		case *timestamppb.Timestamp:
			ts := timestamppb.New(t)
			if err := ts.CheckValid(); err != nil {
				return false, err
			}
			dst.Set(reflect.ValueOf(ts))
			return true, nil
		}

	case *time.Time:
		if t == nil {
			return false, nil
		}
		if t.Unix() <= 0 {
			return false, nil
		}
		// *time.Time -> timestampps.Timestamp or int64
		switch dst.Interface().(type) {
		case *timestamppb.Timestamp:
			ts := timestamppb.New(*t)
			if err := ts.CheckValid(); err != nil {
				return false, err
			}
			dst.Set(reflect.ValueOf(ts))
			return true, nil
		}

	case *timestamppb.Timestamp:
		if t == nil {
			return false, nil
		}
		if t.GetSeconds() <= 0 {
			return false, nil
		}
		if err := t.CheckValid(); err != nil {
			return false, err
		}
		// *timestamppb.Timestamp -> time.Time
		switch dst.Interface().(type) {
		case time.Time:
			dst.Set(reflect.ValueOf(t.AsTime()))
			return true, nil
		case *time.Time:
			dst.Set(reflect.ValueOf(xgo.ToPtr(t.AsTime())))
			return true, nil
		}

	case time.Duration:
		if t.Nanoseconds() <= 0 {
			return false, nil
		}
		switch dst.Interface().(type) {
		case *durationpb.Duration:
			d := durationpb.New(t)
			if err := d.CheckValid(); err != nil {
				return false, err
			}
			dst.Set(reflect.ValueOf(d))
			return true, nil
		}

	case *durationpb.Duration:
		if t == nil {
			return false, nil
		}
		if err := t.CheckValid(); err != nil {
			return false, err
		}
		// *durationpb.Duration -> time.Duration
		switch dst.Interface().(type) {
		case time.Duration:
			dst.Set(reflect.ValueOf(t.AsDuration()))
			return true, nil
		}

	case int64:
		if t <= 0 {
			return false, nil
		}
		// int64 -> time.Time or *timestamppb.Timestamp or *durationpb.Duration
		switch dst.Interface().(type) {
		case *timestamppb.Timestamp:
			ts := timestamppb.New(time.Unix(t, 0))
			if err := ts.CheckValid(); err != nil {
				return false, err
			}
			dst.Set(reflect.ValueOf(ts))
			return true, nil
		case *durationpb.Duration:
			dst.Set(reflect.ValueOf(durationpb.New(time.Duration(t))))
			return true, nil
		}
	}
	return false, nil
}

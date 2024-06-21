package xgo_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/glassonion1/xgo"
)

type Fruits int

const (
	Apple Fruits = iota + 1
	Banana
	Cherry
)

type Weekday string

const (
	Sunday    = Weekday("Sunday")
	Monday    = Weekday("Monday")
	Tuesday   = Weekday("Tuesday")
	Wednesday = Weekday("Wednesday")
	Thursday  = Weekday("Thursday")
	Friday    = Weekday("Friday")
	Saturday  = Weekday("Saturday")
)

func TestNewValues(t *testing.T) {
	type is []interface{}
	cases := []struct {
		name   string        // Name of the test case
		params []interface{} // Params to pass in subsequent calls
	}{
		{"New of bool", is{false, true}},
		{"New of string", is{"", "a"}},
		{"New of int", is{0, 1}},
		{"New of int8", is{int8(0), int8(1)}},
		{"New of int16", is{int16(0), int16(1)}},
		{"New of int32", is{int32(0), int32(1)}},
		{"New of int64", is{int64(0), int64(1)}},
		{"New of uint", is{uint(0), uint(1)}},
		{"New of uint8", is{uint8(0), uint8(1)}},
		{"New of uint16", is{uint16(0), uint16(1)}},
		{"New of uint32", is{uint32(0), uint32(1)}},
		{"New of uint64", is{uint64(0), uint64(1)}},
		{"New of float32", is{float32(0), float32(1)}},
		{"New of float64", is{float64(0), float64(1)}},
		{"New of byte", is{byte(0), byte(1)}},
		{"New of rune", is{rune(0), rune(1)}},
		{"New of Date", is{
			time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1),
		}},
		{"New of Fruit", is{Apple, Banana, Cherry}},
		{"New of Weekday", is{Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday}},
	}

	for _, c := range cases {
		for _, param := range c.params {
			result := xgo.New(param)
			got := *result
			if got != param {
				t.Errorf("[%s(%v)] Expected: %v, got: %v", c.name, param, param, got)
			}
			if reflect.TypeOf(got) != reflect.TypeOf(param) {
				t.Errorf("[%s(%v)] Expected: %T, got: %T", c.name, param, param, got)
			}
		}
	}
}

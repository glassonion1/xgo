package xgo_test

import (
	"fmt"
	"time"

	"github.com/glassonion1/xgo"
)

func ExampleDeepCopy_struct() {
	type FromModel struct {
		ID        string `copier:"Id"`
		Name      string
		CreatedAt time.Time
		UpdatedAt *time.Time
	}

	type ToModel struct {
		Id        string
		Name      string
		CreatedAt time.Time
		UpdatedAt *time.Time
	}

	now := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	from := FromModel{
		ID:        "xxxx",
		Name:      "R2D2",
		CreatedAt: now,
		UpdatedAt: &now,
	}
	to := &ToModel{}
	err := xgo.DeepCopy(from, to)
	if err != nil {
		// handles error
	}
	fmt.Println("ToModel object:", to)

	// Output: ToModel object: &{xxxx R2D2 2025-06-01 00:00:00 +0000 UTC 2025-06-01 00:00:00 +0000 UTC}
}

func ExampleDeepCopy_slice() {
	type FromModel struct {
		ID   string
		Name string
	}
	type ToModel struct {
		ID   string
		Name string
	}
	from := []FromModel{
		{
			ID:   "xxxx1",
			Name: "R2D2",
		},
		{
			ID:   "xxxx2",
			Name: "C3PO",
		},
	}
	to := &[]ToModel{}
	err := xgo.DeepCopy(from, to)
	if err != nil {
		// handles error
	}
	fmt.Println("ToModel object:", to)

	// Output: ToModel object: &[{xxxx1 R2D2} {xxxx2 C3PO}]
}

func ExampleContains() {
	// slice of int32
	containsInt32 := xgo.Contains([]int32{1, 2, 3, 4, 5}, 3)
	fmt.Println("contains int32:", containsInt32)

	// slice of int
	containsInt := xgo.Contains([]int{1, 2, 3, 4, 5}, 2)
	fmt.Println("contains int:", containsInt)

	// slice of float64
	containsFloat64 := xgo.Contains([]float64{1.1, 2.2, 3.3, 4.4, 5.5}, 4.4)
	fmt.Println("contains float64:", containsFloat64)

	// slice of string
	containsString := xgo.Contains([]string{"r2d2", "c3po", "bb8"}, "c3po")
	fmt.Println("contains string:", containsString) // -> true

	// slice of struct
	type hero struct {
		ID   string
		Name string
	}
	list := []hero{
		{
			ID:   "1",
			Name: "Luke Skywalker",
		},
		{
			ID:   "2",
			Name: "Han Solo",
		},
		{
			ID:   "3",
			Name: "Leia Organa",
		},
	}
	target := hero{
		ID:   "2",
		Name: "Han Solo",
	}
	containsStruct := xgo.Contains(list, target)
	fmt.Println("contains struct:", containsStruct)

	// Output:
	// contains int32: true
	// contains int: true
	// contains float64: true
	// contains string: true
	// contains struct: true
}

type Vegetables string

const (
	Pea     Vegetables = "Pea"
	Okra    Vegetables = "Okra"
	Pumpkin Vegetables = "Pumpkin"
)

type Model struct {
	ID        *int
	Name      *string
	Material  *Vegetables
	CreatedAt *time.Time
}

func ExampleToPtr() {
	obj := Model{
		ID:        xgo.ToPtr(123),
		Name:      xgo.ToPtr("R2D2"),
		Material:  xgo.ToPtr(Pea),
		CreatedAt: xgo.ToPtr(time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC)),
	}
	fmt.Println("object:", obj)
}

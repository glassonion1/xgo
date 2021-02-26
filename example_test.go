package xgo_test

import (
	"fmt"
	"time"

	"github.com/glassonion1/xgo"
	"github.com/golang/protobuf/ptypes/timestamp"
)

// It is a common, ordinary struct
type FromModel struct {
	ID        string `copier:"Id"`
	Name      string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// It is like a protobuf struct on gRPC
type ToModel struct {
	Id        string
	Name      string
	CreatedAt *timestamp.Timestamp
	UpdatedAt *timestamp.Timestamp
}

func ExampleDeepCopy() {
	now := time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC)
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

	// Output: ToModel object: &{xxxx R2D2 seconds:1590969600 seconds:1590969600}
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

	// slice of struct
	type hero struct {
		ID   string
		Name string
	}
	list := []hero{
		hero{
			ID:   "1",
			Name: "Luke Skywalker",
		},
		hero{
			ID:   "2",
			Name: "Han Solo",
		},
		hero{
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
	// contains struct: true
}

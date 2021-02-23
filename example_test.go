package xgo_test

import (
	"fmt"
	"time"

	"github.com/glassonion1/xgo"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type FromModel struct {
	ID        string `copier:"Id"`
	Name      string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
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
		Name:      "foo",
		CreatedAt: now,
		UpdatedAt: &now,
	}
	to := &ToModel{}
	err := xgo.DeepCopy(from, to)
	if err != nil {
		// handles error
	}
	fmt.Printf("prints ToModel object: %v", to)

	// Output: prints ToModel object: &{xxxx foo seconds:1590969600 seconds:1590969600}
}

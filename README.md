# xgo - Useful Go library

[![Test CLI](https://github.com/glassonion1/xgo/actions/workflows/test.yml/badge.svg)](https://github.com/glassonion1/xgo/actions/workflows/test.yml)
[![Godoc](https://img.shields.io/badge/godoc-reference-blue)](https://godoc.org/github.com/glassonion1/xgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/glassonion1/xgo)](https://goreportcard.com/report/github.com/glassonion1/xgo)
[![GitHub license](https://img.shields.io/github/license/glassonion1/xgo)](https://github.com/glassonion1/xgo/blob/main/LICENSE)

The xgo contains various useful features for gohers.

## Features
- Deep copy
- Contains
- Chunk
- Exponential backoff
- Struct to map
- Obtain pointers to types

### Protobuf support
There is an extension package that supports deep copying Protobuf timestamp types.
- [xgo pb](https://github.com/glassonion1/xgo/tree/main/xgopb)

## Install
```
$ go get github.com/glassonion1/xgo
```

## Import
```go
import "github.com/glassonion1/xgo"
```

## Usage
### Deep copy
field-to-field copying based on matching names.

In layered architecture, there are times when we need to copy values between structs that have different types but share the same field names. 
The DeepCopy function is designed to be used in such scenarios.

Support for copying data:
- from struct to struct
- from struct to pointer
- from pointer to struct
- from slice to slice
#### from struct to struct
```go
package xgo_test

import (
    "fmt"
    "time"

    "github.com/glassonion1/xgo"
)

type FromModel struct {
    ID         string `copier:"Id"`
    Name       string
    CreatedAt  time.Time
    UpdatedAt  *time.Time
}
type ToModel struct {
    Id         string
    Name       string
    CreatedAt  time.Time
    UpdatedAt  *time.Time
}

func Example() {
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
```

#### from slice to slice
```go
type FromModel struct {
    ID         string
    Name       string
}
type ToModel struct {
    ID         string
    Name       string
}

func Example() {
    from := []FromModel{
        {
            ID: "xxxx1",
            Name: "R2D2",
        },
    {
            ID: "xxxx2",
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
```

### Contains
Contains method for a slice.
```go
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
fmt.Println(containsString) // -> true

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
// contains string: true
// contains struct: true
```

### ToPtr
Obtain pointers to types
```go
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
```

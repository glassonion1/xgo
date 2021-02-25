# XGo - Useful Go libraries

[![Test CLI](https://github.com/glassonion1/xgo/actions/workflows/test.yml/badge.svg)](https://github.com/glassonion1/xgo/actions/workflows/test.yml)
[![Godoc](https://img.shields.io/badge/godoc-reference-blue)](https://godoc.org/github.com/glassonion1/xgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/glassonion1/xgo)](https://goreportcard.com/report/github.com/glassonion1/xgo)
[![GitHub license](https://img.shields.io/github/license/glassonion1/xgo)](https://github.com/glassonion1/xgo/blob/main/LICENSE)

XGo contains various useful Go libraries

## Features
- Deep copy
- Contains
- Chunk
- Exponential backoff
- Struct to map

## Install
```
$ go get github.com/glassonion1/xgo
```

## Usage
### Deep copy
```go
package xgo_test

import (
    "fmt"
    "time"

    "github.com/golang/protobuf/ptypes/timestamp"
    "github.com/glassonion1/xgo"
)

// It is a common, ordinary struct
type FromModel struct {
    ID         string `copier:"Id"`
    Name       string
    CreatedAt  time.Time
    UpdatedAt  *time.Time
}
// It is like a protobuf struct on gRPC
type ToModel struct {
    Id         string
    Name       string
    CreatedAt  *timestamp.Timestamp
    UpdatedAt  *timestamp.Timestamp
}

func Example() {
    now := time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC)
    from := FromModel{
        ID: "xxxx",
        Name: "foo",
        CreatedAt: now,
        UpdatedAt: &now,
    }
    to := &ToModel{}
    err := xgo.DeepCopy(from, to)
    if err != nil {
        // handles error
    }
    fmt.Println("ToModel object:", to)
    
    // Output: ToModel object: &{xxxx foo seconds:1590969600 seconds:1590969600}
}
```

### Contains
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

// slice of struct
type item struct {
    ID   string
    Name string
}
list := []item{
    item{
        ID:   "1",
        Name: "test1",
    },
    item{
        ID:   "2",
        Name: "test2",
    },
    item{
        ID:   "3",
        Name: "test3",
    },
}
target := item{
	ID:   "2",
	Name: "test2",
}
containsStruct := xgo.Contains(list, target)
fmt.Println("contains struct:", containsStruct)

// Output:
// contains int32: true
// contains int: true
// contains float64: true
// contains struct: true
```
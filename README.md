# XGo - Useful golang libraries

[![Test CLI](https://github.com/glassonion1/xgo/actions/workflows/test.yml/badge.svg)](https://github.com/glassonion1/xgo/actions/workflows/test.yml)
[![Godoc](https://img.shields.io/badge/godoc-reference-blue)](https://godoc.org/github.com/glassonion1/xgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/glassonion1/xgo)](https://goreportcard.com/report/github.com/glassonion1/xgo)
[![GitHub license](https://img.shields.io/github/license/glassonion1/xgo)](https://github.com/glassonion1/xgo/blob/main/LICENSE)

XGo contains various useful golang libraries

## Features
- Deep copy
- Exponential backoff
- Chunk
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

type FromModel struct {
    ID         string `copier:"Id"`
    Name       string
    CreatedAt  time.Time
    UpdatedAt  *time.Time
}
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
    fmt.Printf("prints ToModel object: %v", to)
    
    // Output: prints ToModel object: &{xxxx foo seconds:1590969600 seconds:1590969600}
}
```

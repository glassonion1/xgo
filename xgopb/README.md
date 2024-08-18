# xgo pb - protobuf extension
The xgo pb is a protobuf extension of xgo.

## Features
- Deep copy

## Install
```
$ go get github.com/glassonion1/xgo/xgopb
```

## Import
```go
import "github.com/glassonion1/xgo/xgopb"
```

## Usage
### Deep copy

```go
package xgopb_test

import (
    "fmt"
    "time"

    "github.com/golang/protobuf/ptypes/timestamp"
    "github.com/glassonion1/xgo/xgopb"
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
    now := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
    from := FromModel{
        ID: "xxxx",
        Name: "R2D2",
        CreatedAt: now,
        UpdatedAt: &now,
    }
    to := &ToModel{}
    err := xgo.DeepCopy(from, to)
    if err != nil {
        // handles error
    }
    fmt.Println("ToModel object:", to)
    
    // Output: ToModel object: &{xxxx R2D2 seconds:1748736000 seconds:1748736000}
}
```

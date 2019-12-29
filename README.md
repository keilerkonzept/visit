# visit

[![](https://godoc.org/github.com/keilerkonzept/visit?status.svg)](http://godoc.org/github.com/keilerkonzept/visit) [![Go Report Card](https://goreportcard.com/badge/github.com/keilerkonzept/visit)](https://goreportcard.com/report/github.com/keilerkonzept/visit)

A Go library to recursively visit data structures using reflection.

```go
import "github.com/keilerkonzept/visit"
```

## Get it

```sh
go get -u "github.com/keilerkonzept/visit"
```

## Use it

```go
import (
    "github.com/keilerkonzept/visit"

    "fmt"
    "reflect"
)

type myStruct struct {
    String string
    Map    map[string]myStruct
    Ptr    *myStruct
}

func main() {
	obj := &myStruct{
		String: "hello",
		Map: map[string]myStruct{
			"world": myStruct{String: "!"},
		},
	}
	obj.Ptr = obj

	var strings []string
	visit.Values(obj, func(value visit.ValueWithParent) (visit.Action, error) {
		if value.Kind() == reflect.String {
			strings = append(strings, value.String())
		}
		return visit.Continue, nil
	})
	fmt.Println(strings)
	// Output:
	// [hello world !]
}
```

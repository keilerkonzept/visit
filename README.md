# visitstruct

[![](https://godoc.org/github.com/keilerkonzept/visitstruct?status.svg)](http://godoc.org/github.com/keilerkonzept/visitstruct) [![](https://goreportcard.com/badge/github.com/keilerkonzept/visitstruct/goreportcard)](https://goreportcard.com/report/github.com/keilerkonzept/visitstruct)

A Go library to recursively visit data structures using reflection.

```go
import "github.com/keilerkonzept/visitstruct"
```

## Get it

```sh
go get -u "github.com/keilerkonzept/visitstruct"
```

## Use it

```go
import "github.com/keilerkonzept/visitstruct"

func main() {
   type myStruct struct {
		String string
		Map    map[string]myStruct
		Ptr    *myStruct
	}
	obj := &myStruct{
		String: "hello",
		Map: map[string]myStruct{
			"world": myStruct{String: "!"},
		},
	}
	obj.Ptr = obj

	var strings []string
	Any(obj, func(v reflect.Value) (action, error) {
		if v.Kind() == reflect.String {
			strings = append(strings, v.String())
		}
		return Continue, nil
	})
	fmt.Println(strings)
	// Output:
	// [hello world !]
}
```

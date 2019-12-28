// Package visitstruct is a library to visit Go data structures (using reflection)
package visitstruct

import (
	"fmt"
	"reflect"
	"testing"
)

func Example() {
	type myStruct struct {
		String string
		Map    map[string]myStruct
		Ptr    *myStruct
	}
	obj := &myStruct{
		String: "hello",
		Map: map[string]myStruct{
			"world": {String: "!"},
		},
	}
	obj.Ptr = obj

	var strings []string
	Any(obj, func(value, parent, index reflect.Value) (Action, error) {
		if value.Kind() == reflect.String {
			strings = append(strings, value.String())
		}
		return Continue, nil
	})
	fmt.Println(strings)
	// Output:
	// [hello world !]
}

func TestAny(t *testing.T) {
	type kitchenSink struct {
		ptr     *kitchenSink
		structs []kitchenSink
		strings []string
		maps    []map[string]interface{}
		single  string
	}
	loopy := kitchenSink{
		structs: []kitchenSink{
			{single: "abc"},
		},
		maps: []map[string]interface{}{
			{
				"hello": 123,
				"world": 456,
			},
		},
		single:  "baz",
		strings: []string{"foo", "bar"},
	}
	loopy.ptr = &loopy
	accumulatedStrings := make(map[string]int)

	type args struct {
		obj interface{}
		f   func(v, p, i reflect.Value) (Action, error)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		out     func() interface{}
		wantOut interface{}
	}{
		{
			name: "kitchen sink",
			args: args{
				obj: &loopy,
				f: func(v, _, _ reflect.Value) (Action, error) {
					if v.Kind() == reflect.String {
						accumulatedStrings[v.String()]++
					}
					return Continue, nil
				},
			},
			out: func() interface{} { return accumulatedStrings },
			wantOut: map[string]int{
				"abc":   1,
				"hello": 1,
				"world": 1,
				"foo":   1,
				"bar":   1,
				"baz":   1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Any(tt.args.obj, tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("Any() error = %v, wantErr %v", err, tt.wantErr)
			}
			if out := tt.out(); !reflect.DeepEqual(out, tt.wantOut) {
				t.Errorf("Any() out = %v, wantOut %v", out, tt.wantOut)

			}
		})
	}
}

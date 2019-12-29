// Package visit is a library to visit Go data structures (using reflection)
package visit

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
	Values(obj, func(value ValueWithParent) (Action, error) {
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
		Ptr     *kitchenSink
		Structs []kitchenSink
		Strings []string
		Maps    []map[string]interface{}
		Single  string
	}
	loopy := kitchenSink{
		Structs: []kitchenSink{
			{Single: "abc"},
		},
		Maps: []map[string]interface{}{
			{
				"hello": 123,
				"world": 456,
			},
		},
		Single:  "baz",
		Strings: []string{"foo", "bar"},
	}
	loopy.Ptr = &loopy
	accumulatedStrings := make(map[string]int)

	rewrite := kitchenSink{
		Single: "abc",
		Maps: []map[string]interface{}{
			{
				"def": []string{"xyz", "uvw"},
			},
		},
		Strings: []string{"foo", "bar"},
	}
	type args struct {
		obj interface{}
		f   func(v ValueWithParent) (Action, error)
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
				f: func(v ValueWithParent) (Action, error) {
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
		{
			name: "rewrite",
			args: args{
				obj: &rewrite,
				f: func(v ValueWithParent) (Action, error) {
					if v.Kind() == reflect.String {
						Assign(v, reflect.ValueOf(v.String()+"(edited)"))
					}
					return Continue, nil
				},
			},
			out: func() interface{} { return rewrite },
			wantOut: kitchenSink{
				Single: "abc(edited)",
				Maps: []map[string]interface{}{
					{
						"def": []string{"xyz(edited)", "uvw(edited)"},
					},
				},
				Strings: []string{"foo(edited)", "bar(edited)"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Values(tt.args.obj, tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("Values() error = %v, wantErr %v", err, tt.wantErr)
			}
			if out := tt.out(); !reflect.DeepEqual(out, tt.wantOut) {
				t.Errorf("Values() out = %v, wantOut %v", out, tt.wantOut)

			}
		})
	}
}

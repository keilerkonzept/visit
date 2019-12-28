// Package visitstruct is a library to visit Go data structures (using reflection)
package visitstruct

import (
	"reflect"
)

type action string

// Visitor actions
const (
	Skip     action = "Skip"
	Stop     action = "Stop"
	Continue action = "Continue"
)

// Any visits a structure using cycle detection
func Any(obj interface{}, f func(reflect.Value) (action, error)) error {
	seen := make(map[uintptr]bool)
	return CycleFree(obj, func(v reflect.Value) (action, error) {
		switch v.Kind() {
		case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
			ptr := v.Pointer()
			if seen[ptr] {
				return Skip, nil
			}
			seen[ptr] = true
		}
		return f(v)
	})
}

// CycleFree visits a structure *without* performing cycle detection
func CycleFree(obj interface{}, f func(reflect.Value) (action, error)) error {
	worklist := []reflect.Value{reflect.ValueOf(obj)}
	for len(worklist) > 0 {
		v := worklist[0]
		worklist[0] = worklist[len(worklist)-1]
		worklist = worklist[:len(worklist)-1]
		action, err := f(v)
		if err != nil {
			return err
		}
		switch action {
		case Skip:
			continue
		case Stop:
			return nil
		case Continue:
		}

		switch v.Kind() {
		case reflect.Map:
			i := v.MapRange()
			for i.Next() {
				worklist = append(worklist, i.Key())
				worklist = append(worklist, i.Value())
			}
		case reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				worklist = append(worklist, v.Index(i))
			}
		case reflect.Array:
			for i := 0; i < v.Len(); i++ {
				worklist = append(worklist, v.Index(i))
			}
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				worklist = append(worklist, v.FieldByIndex([]int{i}))
			}
		case reflect.Ptr:
			if !v.IsNil() {
				worklist = append(worklist, v.Elem())
			}
		}
	}

	return nil
}

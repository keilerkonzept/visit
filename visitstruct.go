// Package visitstruct is a library to visit Go data structures (using reflection)
package visitstruct

import (
	"reflect"
)

type action string

// Action is the type of visitor actions
type Action action

// Visitor actions
const (
	Skip     Action = "Skip"
	Stop     Action = "Stop"
	Continue Action = "Continue"
)

// VisitFunc is a visitor function
type VisitFunc func(value, parent, index reflect.Value) (Action, error)

// Any visits a structure using cycle detection
func Any(obj interface{}, f VisitFunc) error {
	seen := make(map[uintptr]bool)
	return CycleFree(obj, func(value, parent, index reflect.Value) (Action, error) {
		switch value.Kind() {
		case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
			ptr := value.Pointer()
			if seen[ptr] {
				return Skip, nil
			}
			seen[ptr] = true
		}
		return f(value, parent, index)
	})
}

// CycleFree visits a structure *without* performing cycle detection
func CycleFree(obj interface{}, f VisitFunc) error {
	type item struct {
		value, parent, index reflect.Value
	}
	worklist := []item{{value: reflect.ValueOf(obj)}}
	for len(worklist) > 0 {
		current := worklist[0]
		value := current.value
		worklist[0] = worklist[len(worklist)-1]
		worklist = worklist[:len(worklist)-1]
		action, err := f(value, current.parent, current.index)
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

		switch value.Kind() {
		case reflect.Map:
			i := value.MapRange()
			for i.Next() {
				worklist = append(worklist, item{value: i.Key(), parent: value})
				worklist = append(worklist, item{i.Value(), value, i.Key()})
			}
		case reflect.Slice:
			for i := 0; i < value.Len(); i++ {
				worklist = append(worklist, item{value.Index(i), value, reflect.ValueOf(i)})
			}
		case reflect.Array:
			for i := 0; i < value.Len(); i++ {
				worklist = append(worklist, item{value.Index(i), value, reflect.ValueOf(i)})
			}
		case reflect.Struct:
			for i := 0; i < value.NumField(); i++ {
				worklist = append(worklist, item{value.FieldByIndex([]int{i}), value, reflect.ValueOf(i)})
			}
		case reflect.Interface:
			if !value.IsNil() {
				worklist = append(worklist, item{value: value.Elem(), parent: value})
			}
		case reflect.Ptr:
			if !value.IsNil() {
				worklist = append(worklist, item{value: value.Elem(), parent: value})
			}
		}
	}

	return nil
}

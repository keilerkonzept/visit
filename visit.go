// Package visit is a library to visit Go data structures (using reflection)
package visit

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

// Any visits a structure, with cycle detection
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

type worklistItem struct {
	value, parent, index reflect.Value
}

// CycleFree visits a structure *without* performing cycle detection
func CycleFree(obj interface{}, f VisitFunc) error {
	worklist := []worklistItem{{value: reflect.ValueOf(obj)}}
	for len(worklist) > 0 {
		var top worklistItem
		top, worklist = pop(worklist)
		action, err := f(top.value, top.parent, top.index)
		if err != nil {
			return err
		}
		switch action {
		case Skip:
			continue
		case Stop:
			return nil
		}
		worklist = queue(worklist, top.value)
	}
	return nil
}

func pop(worklist []worklistItem) (worklistItem, []worklistItem) {
	top := worklist[0]
	worklist[0] = worklist[len(worklist)-1]
	worklist = worklist[:len(worklist)-1]
	return top, worklist
}

func queue(worklist []worklistItem, childrenOf reflect.Value) []worklistItem {
	switch childrenOf.Kind() {
	case reflect.Map:
		i := childrenOf.MapRange()
		for i.Next() {
			worklist = append(worklist, worklistItem{value: i.Key(), parent: childrenOf})
			worklist = append(worklist, worklistItem{i.Value(), childrenOf, i.Key()})
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < childrenOf.Len(); i++ {
			worklist = append(worklist, worklistItem{childrenOf.Index(i), childrenOf, reflect.ValueOf(i)})
		}
	case reflect.Interface, reflect.Ptr:
		if !childrenOf.IsNil() {
			worklist = append(worklist, worklistItem{value: childrenOf.Elem(), parent: childrenOf})
		}
	case reflect.Struct:
		for i := 0; i < childrenOf.NumField(); i++ {
			worklist = append(worklist, worklistItem{childrenOf.FieldByIndex([]int{i}), childrenOf, reflect.ValueOf(i)})
		}
	}
	return worklist
}

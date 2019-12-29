// Package visit is a library to visit Go data structures (using reflection)
package visit

import "reflect"

type action string

// Action is the type of visitor actions
type Action action

// Visitor actions
const (
	Skip     Action = "Skip"
	Stop     Action = "Stop"
	Continue Action = "Continue"
)

// Func is a visitor function
type Func func(ValueWithParent) (Action, error)

// Values visits a structure, with cycle detection
func Values(obj interface{}, f Func) error {
	seen := make(map[uintptr]bool)
	return ValuesUnsafe(obj, func(v ValueWithParent) (Action, error) {
		switch v.Value.Kind() {
		case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
			ptr := v.Value.Pointer()
			if seen[ptr] {
				return Skip, nil
			}
			seen[ptr] = true
		}
		return f(v)
	})
}

// ValueWithParent is a reflect.Value with its parent container (if any) and the corresponding index value.
type ValueWithParent struct {
	reflect.Value
	Parent *ValueWithParent
	Index  reflect.Value
}

// ValuesUnsafe visits a structure *without* performing cycle detection
func ValuesUnsafe(obj interface{}, f Func) error {
	worklist := []ValueWithParent{{Value: reflect.ValueOf(obj)}}
	for len(worklist) > 0 {
		var top ValueWithParent
		top, worklist = pop(worklist)
		action, err := f(top)
		if err != nil {
			return err
		}
		switch action {
		case Skip:
			continue
		case Stop:
			return nil
		}
		worklist = queue(worklist, top)
	}
	return nil
}

func pop(worklist []ValueWithParent) (ValueWithParent, []ValueWithParent) {
	top := worklist[0]
	worklist[0] = worklist[len(worklist)-1]
	worklist = worklist[:len(worklist)-1]
	return top, worklist
}

func queue(worklist []ValueWithParent, childrenOf ValueWithParent) []ValueWithParent {
	value := childrenOf.Value
	switch value.Kind() {
	case reflect.Map:
		i := value.MapRange()
		for i.Next() {
			worklist = append(worklist, ValueWithParent{Value: i.Key(), Parent: &childrenOf})
			worklist = append(worklist, ValueWithParent{i.Value(), &childrenOf, i.Key()})
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			worklist = append(worklist, ValueWithParent{value.Index(i), &childrenOf, reflect.ValueOf(i)})
		}
	case reflect.Interface, reflect.Ptr:
		if !value.IsNil() {
			worklist = append(worklist, ValueWithParent{Value: value.Elem(), Parent: &childrenOf})
		}
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			worklist = append(worklist, ValueWithParent{value.Field(i), &childrenOf, reflect.ValueOf(i)})
		}
	}
	return worklist
}

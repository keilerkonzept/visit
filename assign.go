package visit

import "reflect"

// Assign tries to assign `newValue` to the given `value`, and otherwise returns an error.
func Assign(value ValueWithParent, newValue reflect.Value) error {
	if !TryAssign(value, newValue) {
		return &reflect.ValueError{
			Method: "Assign",
			Kind:   value.Kind(),
		}
	}
	return nil
}

// TryAssign tries to assign `newValue` to the given `value`, and returns a boolean indicating success.
func TryAssign(value ValueWithParent, newValue reflect.Value) bool {
	if value.CanSet() {
		value.Set(newValue)
		return true
	}
	for value.Parent != nil {
		switch value.Parent.Kind() {
		case reflect.Array, reflect.Slice:
			if value.Index.IsValid() {
				value.Parent.Value.Index(int(value.Index.Int())).Set(newValue)
				return true
			}
			return false
		case reflect.Map:
			if value.Index.IsValid() {
				value.Parent.SetMapIndex(value.Index, newValue)
				return true
			}
			return false
		case reflect.Interface:
			value = *value.Parent
		default:
			return false
		}
	}
	return false
}

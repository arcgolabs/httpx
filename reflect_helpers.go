package httpx

import "reflect"

func indirectStructType[T any]() (reflect.Type, bool, bool) {
	valueType := reflect.TypeFor[T]()
	hasNestedPointer := false
	for valueType.Kind() == reflect.Pointer {
		hasNestedPointer = true
		valueType = valueType.Elem()
	}
	return valueType, hasNestedPointer, valueType.Kind() == reflect.Struct
}

func indirectStructValue[T any](input *T) (reflect.Value, bool) {
	if input == nil {
		return reflect.Value{}, false
	}

	value := reflect.ValueOf(input)
	if !value.IsValid() || value.IsNil() {
		return reflect.Value{}, false
	}

	for value.IsValid() && value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return reflect.Value{}, false
		}
		value = value.Elem()
	}

	if !value.IsValid() || value.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}
	return value, true
}

package util

import (
	"fmt"
	"reflect"
)

func CoerceToSlice(maybeSlice interface{}) ([]interface{}, error) {
	s := reflect.ValueOf(maybeSlice)
	if s.Kind() != reflect.Slice {
		return nil, fmt.Errorf("input was not of type slice")
	}

	if s.IsNil() {
		return nil, fmt.Errorf("input value was nil")
	}

	result := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		result[i] = s.Index(i).Interface()
	}

	return result, nil
}

package makefunc

import (
	"fmt"
	"reflect"
)

func MakeFunc(f interface{}, returnArgs []reflect.Type, inputParams ...interface{}) (func() error, error) {
	if err := ValidateFunc(f, returnArgs, inputParams...); err != nil {
		return nil, fmt.Errorf("invalid function: %v", err)
	}
	typ := reflect.TypeOf(f)
	paramValues := make([]reflect.Value, 0, len(inputParams))
	for idx, arg := range inputParams {
		if arg == nil {
			paramValues = append(paramValues, reflect.Zero(typ.In(idx)))
		} else {
			paramValues = append(paramValues, reflect.ValueOf(arg))
		}
	}
	return func() error {
		// TODO use reflect.Value.Convert since we validate with ConvertibleTo
		ret := reflect.ValueOf(f).Call(paramValues)
		if ret[0].IsNil() {
			return nil
		}
		// the return value being an `error` object is guaranteed by
		// `validateFunction` called above.
		return ret[0].Interface().(error)
	}, nil
}

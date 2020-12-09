package makefunc

import (
	"fmt"
	"reflect"
)

// CheckTypeCompatibility checks whether `got` is a valid type for `want`, and
// returns a descriptive error otherwise.
//
// Examples of positive compatibility:
//   want=int, got=int
//   want=interface{}, got=string
//   want=interface{}, got=nil
//   want=error, got=untyped nil
//   want=error, got=error-typed nil
//   want=io.ReadWriter, got=io.Writer
//   want=[]string, got=[]string
//   want=struct{}, got=struct{}
//
// Examples of negative compatibility:
//   want=error, got=string nil
//   want=io.Reader, got=io.Writer
//   want=[]struct{}, got []string
//   want=*string, got=string
//   want=*string, got=*int
//   want=struct{}, got=*struct{}
//   want=[]interface{}, got=[]string
func CheckTypeCompatibility(want, got reflect.Type) error {
	if want == nil {
		return fmt.Errorf("`want` cannot be nil")
	}
	if got == nil {
		// this happens when an untyped nil is passed.
		// Check if `want` is a nil-able type.
		switch want.Kind() {
		case reflect.Func, reflect.Chan, reflect.Slice, reflect.Map, reflect.Interface, reflect.Ptr:
			return nil
		default:
			return fmt.Errorf("incompatible types: got untyped `nil` but %v is not nil-able", want)
		}
	}

	switch want.Kind() {
	case reflect.Interface:
		if !got.Implements(want) {
			return fmt.Errorf("invalid type '%v', does not implement '%v'", got, want)
		}
	default:
		if !got.ConvertibleTo(want) {
			return fmt.Errorf("invalid type '%v', it is not convertible to %v", got, want)
		}
		if !got.AssignableTo(want) {
			return fmt.Errorf("invalid type '%v', it is not assignable to %v", got, want)
		}
	}
	return nil
}

// ValidateFunc checks whether `f` is a function, and whether
// the input parameters and the return arguments have the correct number and types.
func ValidateFunc(f interface{}, returnArgs []reflect.Type, inputParams ...interface{}) error {
	typ := reflect.TypeOf(f)
	if kind := typ.Kind(); kind != reflect.Func {
		return fmt.Errorf("invalid object, want %v, got %v", reflect.Func, kind)
	}
	if typ.NumOut() != len(returnArgs) {
		return fmt.Errorf("wrong number of return arguments, want %d, got %d", typ.NumOut(), len(returnArgs))
	}
	for idx := 0; idx < len(returnArgs); idx++ {
		if typ.Out(idx) != returnArgs[idx] {
			return fmt.Errorf("invalid return argument at index %d: want %v, got %v", idx, typ.Out(idx), returnArgs[idx])
		}
	}
	// check input parameter types
	if typ.IsVariadic() {
		// check non-variadic arguments
		if len(inputParams) < typ.NumIn()-1 {
			return fmt.Errorf("wrong number of parameters for variadic function, want at least %d, got %d", typ.NumIn()-1, len(inputParams))
		}
		for idx := 0; idx < typ.NumIn()-1; idx++ {
			want, got := typ.In(idx), reflect.TypeOf(inputParams[idx])
			if err := CheckTypeCompatibility(want, got); err != nil {
				return fmt.Errorf("incompatible type at index %d: %v", idx, err)
			}
		}
		// check variadic arguments
		vargs := typ.In(typ.NumIn() - 1)
		// variadic arguments are passed as a slice of same-kind objects, e.g. []string.
		// Check if this is a slice.
		if vargs.Kind() != reflect.Slice {
			return fmt.Errorf("invalid kind of variadic arguments: want %v, got %v", reflect.Slice, vargs.Kind())
		}
		// Then check that all the passed variadic arguments have compatible types with
		// the function's signture.
		want := vargs.Elem()
		for idx := typ.NumIn() - 1; idx < len(inputParams); idx++ {
			got := reflect.TypeOf(inputParams[idx])
			if err := CheckTypeCompatibility(want, got); err != nil {
				return fmt.Errorf("incompatible type at index %d: %v", idx, err)
			}
		}
	} else {
		// checking a non-variadic function
		if typ.NumIn() != len(inputParams) {
			return fmt.Errorf("wrong number of parameters, want %d, got %d", typ.NumIn(), len(inputParams))
		}
		for idx := 0; idx < typ.NumIn(); idx++ {
			want, got := typ.In(idx), reflect.TypeOf(inputParams[idx])
			if err := CheckTypeCompatibility(want, got); err != nil {
				return fmt.Errorf("incompatible type at index %d: %v", idx, err)
			}
		}
	}

	return nil
}

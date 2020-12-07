package makefunc

import (
	"fmt"
	"reflect"
)

var (
	ErrorType = reflect.TypeOf((*error)(nil)).Elem()
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
//   want=[]interface{}, got=[]string
//   want=struct{}, got=struct{}
//
// Examples of negative compatibility:
//   want=error, got=string nil
//   want=io.Reader, got=io.Writer
//   want=[]struct{}, got []string
//   want=*string, got=string
//   want=*string, got=*int
//   want=struct{}, got=*struct{}
func CheckTypeCompatibility(want, got reflect.Type) error {
	if want == nil {
		return fmt.Errorf("`want` cannot be nil")
	}
	if got == nil {
		// this happens when an untyped nil is passed
		return fmt.Errorf("not implemented yet")
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

// ValidateFunction checks whether `f` is a function, and whether
// the passed parameters have the correct number and types.
func ValidateFunction(f interface{}, args ...interface{}) error {
	typ := reflect.TypeOf(f)
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	if kind := typ.Kind(); kind != reflect.Func {
		return fmt.Errorf("invalid object, want %v, got %v", reflect.Func, kind)
	}
	// FIXME allow specifying a return type, instead of statically checking that
	//       the function returns `error`.
	switch {
	case typ.NumOut() == 1:
		// return type must be error
		if typ.Out(0) != errorType {
			return fmt.Errorf("invalid return type, want %v, got %v", errorType, typ.Out(0))
		}
	default:
		return fmt.Errorf("invalid number of returned parameters, want 1, got %d", typ.NumOut())
	}
	// check input parameter types
	if typ.IsVariadic() {
		// check non-variadic arguments
		if len(args) < typ.NumIn()-1 {
			return fmt.Errorf("wrong number of parameters for variadic function, want at least %d, got %d", typ.NumIn()-1, len(args))
		}
		for idx := 0; idx < typ.NumIn()-1; idx++ {
			want, got := typ.In(idx), reflect.TypeOf(args[idx])
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
		for idx := typ.NumIn() - 1; idx < len(args); idx++ {
			got := reflect.TypeOf(args[idx])
			if err := CheckTypeCompatibility(want, got); err != nil {
				return fmt.Errorf("incompatible type at index %d: %v", idx, err)
			}
		}
	} else {
		// checking a non-variadic function
		if typ.NumIn() != len(args) {
			return fmt.Errorf("wrong number of parameters, want %d, got %d", typ.NumIn(), len(args))
		}
		for idx := 0; idx < typ.NumIn(); idx++ {
			want, got := typ.In(idx), reflect.TypeOf(args[idx])
			if err := CheckTypeCompatibility(want, got); err != nil {
				return fmt.Errorf("incompatible type at index %d: %v", idx, err)
			}
		}
	}

	return nil
}

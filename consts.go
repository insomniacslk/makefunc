package makefunc

import "reflect"

// exported shortcuts to make it easier to call ValidateFunction
var (
	ErrorType = reflect.TypeOf((*error)(nil)).Elem()
)

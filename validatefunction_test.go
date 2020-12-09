package makefunc

import (
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckTypeCompatibilityInt(t *testing.T) {
	intType := reflect.TypeOf(int(0))
	one := int(1)
	// want=int, got=int
	assert.NoError(t, CheckTypeCompatibility(intType, reflect.TypeOf(1)))
	// want=int, got=*int
	assert.Error(t, CheckTypeCompatibility(intType, reflect.TypeOf(&one)))
	// want=int, got=bool
	assert.Error(t, CheckTypeCompatibility(intType, reflect.TypeOf(true)))
	// want=int, got=float
	assert.Error(t, CheckTypeCompatibility(intType, reflect.TypeOf(1.1)))
	// want=int, got=string
	assert.Error(t, CheckTypeCompatibility(intType, reflect.TypeOf("string")))
	// want=int, got=error
	assert.Error(t, CheckTypeCompatibility(intType, ErrorType))
	// want=int, got=[]int
	assert.Error(t, CheckTypeCompatibility(intType, reflect.TypeOf([]int{})))
	// want=int, got=[]interface{}
	assert.Error(t, CheckTypeCompatibility(intType, reflect.TypeOf([]interface{}{})))
	// want=int, got=[]struct{}
	assert.Error(t, CheckTypeCompatibility(intType, reflect.TypeOf([]struct{}{})))
}

func TestCheckTypeCompatibilityIntPtr(t *testing.T) {
	intPtr := int(0)
	intPtrType := reflect.TypeOf(&intPtr)
	one := int(1)
	onePointOne := 1.1
	trueObj := true
	// want=*int, got=*int
	assert.NoError(t, CheckTypeCompatibility(intPtrType, reflect.TypeOf(&one)))
	// want=*int, got=int
	assert.Error(t, CheckTypeCompatibility(intPtrType, reflect.TypeOf(one)))
	// want=*int, got=bool
	assert.Error(t, CheckTypeCompatibility(intPtrType, reflect.TypeOf(trueObj)))
	// want=*int, got=*bool
	assert.Error(t, CheckTypeCompatibility(intPtrType, reflect.TypeOf(&trueObj)))
	// want=*int, got=*float
	assert.Error(t, CheckTypeCompatibility(intPtrType, reflect.TypeOf(&onePointOne)))
	// want=*int, got=float
	assert.Error(t, CheckTypeCompatibility(intPtrType, reflect.TypeOf(onePointOne)))
	// want=*int, got=string
	assert.Error(t, CheckTypeCompatibility(intPtrType, reflect.TypeOf("string")))
	// want=*int, got=error
	assert.Error(t, CheckTypeCompatibility(intPtrType, ErrorType))
	// want=*int, got=[]int
	assert.Error(t, CheckTypeCompatibility(intPtrType, reflect.TypeOf([]int{})))
	// want=*int, got=[]interface{}
	assert.Error(t, CheckTypeCompatibility(intPtrType, reflect.TypeOf([]interface{}{})))
}

// TODO test error interface, slices, pointers, maps

func TestValidateFunctionValidFunctionReturningError(t *testing.T) {
	assert.NoError(t, ValidateFunction(func() error {
		return nil
	},
		[]reflect.Type{ErrorType},
	))
	assert.NoError(t, ValidateFunction(func(string) error {
		return nil
	},
		[]reflect.Type{ErrorType},
		"test",
	))
	assert.NoError(t, ValidateFunction(func(string, int) error {
		return nil
	},
		[]reflect.Type{ErrorType},
		"test", 1,
	))
	assert.NoError(t, ValidateFunction(func(string, int) error {
		return nil
	},
		[]reflect.Type{ErrorType},
		"test", 1,
	))
	assert.NoError(t, ValidateFunction(func([]string) error {
		return nil
	},
		[]reflect.Type{ErrorType},
		[]string{},
	))
	assert.NoError(t, ValidateFunction(func([]string) error {
		return nil
	},
		[]reflect.Type{ErrorType},
		[]string{"test"},
	))
	assert.NoError(t, ValidateFunction(func([]string) error {
		return nil
	},
		[]reflect.Type{ErrorType},
		[]string{"test", "tset"},
	))
	assert.NoError(t, ValidateFunction(func(string, ...interface{}) error {
		return nil
	},
		[]reflect.Type{ErrorType},
		"test", 1,
	))
	assert.NoError(t, ValidateFunction(func(string, ...interface{}) error {
		return nil
	},
		[]reflect.Type{ErrorType},
		"test", 1, 2, 3,
	))
	assert.NoError(t, ValidateFunction(func(string, ...int) error {
		return nil
	},
		[]reflect.Type{ErrorType},
		"test", 1, 2, 3,
	))
}

func TestValidateFunctionValidFunctionReturningNothing(t *testing.T) {
	assert.NoError(t, ValidateFunction(func() {
	},
		[]reflect.Type{},
	))
	assert.NoError(t, ValidateFunction(func(string) {
	},
		[]reflect.Type{},
		"test",
	))
	assert.NoError(t, ValidateFunction(func(string, int) {
	},
		[]reflect.Type{},
		"test", 1,
	))
	assert.NoError(t, ValidateFunction(func(string, int) {
	},
		[]reflect.Type{},
		"test", 1,
	))
	assert.NoError(t, ValidateFunction(func([]string) {
	},
		[]reflect.Type{},
		[]string{},
	))
	assert.NoError(t, ValidateFunction(func([]string) {
	},
		[]reflect.Type{},
		[]string{"test"},
	))
	assert.NoError(t, ValidateFunction(func([]string) {
	},
		[]reflect.Type{},
		[]string{"test", "tset"},
	))
	assert.NoError(t, ValidateFunction(func(string, ...interface{}) {
	},
		[]reflect.Type{},
		"test", 1,
	))
	assert.NoError(t, ValidateFunction(func(string, ...interface{}) {
	},
		[]reflect.Type{},
		"test", 1, 2, 3,
	))
	assert.NoError(t, ValidateFunction(func(string, ...int) {
	},
		[]reflect.Type{},
		"test", 1, 2, 3,
	))
}

func TestValidateFunctionInvalidFunctions(t *testing.T) {
	// not a function
	assert.Error(t, ValidateFunction(1, []reflect.Type{ErrorType}))
	// function with wrong return type
	assert.Error(t, ValidateFunction(func() {}, []reflect.Type{ErrorType}))
	assert.Error(t, ValidateFunction(func() int { return 1 }, []reflect.Type{ErrorType}))
	assert.Error(t, ValidateFunction(func() (int, error) { return 1, nil }, []reflect.Type{ErrorType}))
	// wrong number of passed arguments in non-variadic function
	assert.Error(t, ValidateFunction(func() error { return nil }, []reflect.Type{ErrorType}, 1))
	assert.Error(t, ValidateFunction(func(int) error { return nil }, []reflect.Type{ErrorType}))
	assert.Error(t, ValidateFunction(func(int) error { return nil }, []reflect.Type{ErrorType}, 1, 2, 3))
	// arguments that do not implement the function's input types
	assert.Error(t, ValidateFunction(func(int) error { return nil }, []reflect.Type{ErrorType}, "string"))
	// wrong number of passed arguments in variadic function
	assert.Error(t, ValidateFunction(func(string, ...string) error { return nil }, []reflect.Type{ErrorType}))
	assert.Error(t, ValidateFunction(func(string, string, ...string) error { return nil }, []reflect.Type{ErrorType}, "string"))
	// wrong type passed to variadic function
	assert.Error(t, ValidateFunction(func(...int) error { return nil }, []reflect.Type{ErrorType}, "string"))
	assert.Error(t, ValidateFunction(func(int, ...int) error { return nil }, []reflect.Type{ErrorType}, "string"))
	// wrong interface implemented in variadic function
	assert.Error(t, ValidateFunction(func(...error) error { return nil }, []reflect.Type{ErrorType}, 1))
	assert.Error(t, ValidateFunction(func(...io.Reader) error { return nil }, []reflect.Type{ErrorType}, 1))
}

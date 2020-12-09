package makefunc

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type someError struct{}

func (s someError) Error() string {
	return "some error"
}

func TestMakeFunctionThatReturnsErrorReturnTypedNilVariadic(t *testing.T) {
	var typedNil *someError
	fn, err := MakeFunc(
		func(e error, errs ...error) error { return e },
		[]reflect.Type{ErrorType},
		typedNil,
	)
	require.NoError(t, err)
	assert.Nil(t, fn())
}

func TestMakeFunctionThatReturnsErrorReturnUntypedNilVariadic(t *testing.T) {
	fn, err := MakeFunc(
		func(e error, errs ...error) error { return e },
		[]reflect.Type{ErrorType},
		nil,
	)
	require.NoError(t, err)
	assert.Nil(t, fn())
}

func TestMakeFunctionThatReturnsErrorReturnErr(t *testing.T) {
	fn, err := MakeFunc(
		func(e error) error { return e },
		[]reflect.Type{ErrorType},
		errors.New("err"),
	)
	require.NoError(t, err)
	assert.Error(t, fn())
}

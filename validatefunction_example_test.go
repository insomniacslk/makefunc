package makefunc

import (
	"encoding/json"
	"log"
	"reflect"
)

func ExampleValidateFunction() {
	kv := struct{ Key string }{}
	// the signature of json.Unmarshal is
	// func Unmarshal(data []byte, v interface{}) error
	if err := ValidateFunction(json.Unmarshal, []reflect.Type{ErrorType}, []byte(`{"key": "value"}`), kv); err != nil {
		log.Fatal(err)
	}
	// Output:
}

package makefunc

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

func ExampleMakeFunc() {
	kv := struct{ Key string }{}
	// the signature of json.Unmarshal is
	// func Unmarshal(data []byte, v interface{}) error
	fn, err := MakeFunc(json.Unmarshal, []reflect.Type{ErrorType}, []byte(`{"key": "value"}`), &kv)
	if err != nil {
		log.Fatal(err)
	}
	if err := fn(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(kv)
	// Output: {value}
}

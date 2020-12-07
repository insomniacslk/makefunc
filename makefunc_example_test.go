package makefunc

import (
	"encoding/json"
	"log"
)

func ExampleValidateFunction() {
	kv := struct{ Key string }{}
	// the signature of json.Unmarshal is
	// func Unmarshal(data []byte, v interface{}) error
	if err := ValidateFunction(json.Unmarshal, []byte(`{"key": "value"}`), kv); err != nil {
		log.Fatal(err)
	}
	// Output:
}

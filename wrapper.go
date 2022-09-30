package albatross

import (
	"encoding/json"
)

// JsonUnwrapper is an interface to retrieve a wrapped json.RawMessage
// composed in another type
type JsonUnwrapper interface {
	GetErr() error
	GetObject() json.RawMessage
}

// UnwrapObject is a generic function that takes a type of JsonUnwrapper
// and returns the object in the desired type.
//
// This function will first attempt to assert the object to the desired type.
// If this fails it will attempt to decode the object as JSON data to the desired type.
func UnwrapObject[T any](obj JsonUnwrapper) (T, error) {
	var data T
	if err := obj.GetErr(); err != nil {
		return data, err
	}

	rawJson := obj.GetObject()
	err := json.Unmarshal(rawJson, &data)
	return data, err
}

package albatross

import (
	"encoding/json"
)

// JsonUnwrapper is an interface to retrieve a wrapped json.RawMessage
// composed in another type
type JsonUnwrapper interface {
	GetErr() error
	GetWrapped() json.RawMessage
}

// UnwrapObject is a generic function that takes a type of JsonUnwrapper
// and returns the wrapped object in the desired type.
func UnwrapObject[T any](obj JsonUnwrapper) (T, error) {
	var data T
	if err := obj.GetErr(); err != nil {
		return data, err
	}

	rawJson := obj.GetWrapped()
	err := json.Unmarshal(rawJson, &data)
	return data, err
}

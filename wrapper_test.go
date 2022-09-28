package albatross

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ JsonUnwrapper = (*Objector)(nil)

// Objector is a type used to test the ObjectRetriever interface
type Objector struct {
	Obj json.RawMessage
	Err error
}

func (o *Objector) GetObject() json.RawMessage { return o.Obj }

func (o *Objector) GetErr() error { return o.Err }

func TestAssertionOk(t *testing.T) {
	obj := &Objector{Err: nil, Obj: []byte("10")}
	retrieved, err := GetObject[int](obj)
	if err != nil {
		t.Fatalf("Could not retrieve object as int")
	}
	assert.Equal(t, retrieved, 10, "Retrieved object is not an int of 10")
}

func TestUnmarshallnOk(t *testing.T) {
	obj := &Objector{Err: nil, Obj: []byte("10")}
	retrieved, err := GetObject[int](obj)
	if err != nil {
		t.Fatalf("Could not retrieve object as int: %s", err)
	}
	assert.Equal(t, retrieved, 10, "Retrieved object is not an int of 10")
}

func TestUnmarshallError(t *testing.T) {
	obj := &Objector{Err: nil, Obj: []byte(`{"value":10}`)}
	_, err := GetObject[int](obj)
	if err == nil {
		t.Fatalf("Test should fail to cast an object to an int")
	}
}

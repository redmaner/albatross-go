package albatross

import (
	"encoding/json"
	"fmt"
	"time"
)

var _ JsonUnwrapper = (*JsonRPCResponse)(nil)
var _ error = (*JsonRPCError)(nil)

type ID interface {
	int | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | string | *string
}

// JsonRPCRequest represents a JSON-RPC 2.0 request
type JsonRPCRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Id      any           `json:"id,omitempty"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func NewRPCRequest(method string, params ...interface{}) *JsonRPCRequest {
	return NewRPCRequestWithID(method, time.Now().Unix(), params)
}

func NewRPCRequestWithID[T ID](method string, id T, params ...interface{}) *JsonRPCRequest {
	requestParams := []interface{}{}
	if len(params) > 0 {
		requestParams = append(requestParams, params...)
	}

	return &JsonRPCRequest{
		Jsonrpc: "2.0",
		Id:      id,
		Method:  method,
		Params:  requestParams,
	}
}

// JsonRPCResponse represents a JSON-RPC 2.0 response
// Data is unmarshalled as `json.RawMessage` type and can be further
// unmarshalled with the `UnwrapObject` function
type JsonRPCResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	Id      any             `json:"id"`
	Error   *JsonRPCError   `json:"error,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (r *JsonRPCResponse) GetErr() error {
	if r.Error == nil {
		return nil
	}
	return r.Error
}

func (r *JsonRPCResponse) GetObject() json.RawMessage {
	return r.Data
}

// JsonRPCError represents a JSON-RPC 2.0 error.
// This implements the Go error interface and it is usually not reqired to interact with this type directly.
type JsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

func (e *JsonRPCError) Error() string {
	if len(e.Data) > 0 {
		return fmt.Sprintf("JSON-RPC Error %d - %s. Error data: %s", e.Code, e.Message, e.Data)
	}
	return fmt.Sprintf("JSON-RPC Error %d - %s", e.Code, e.Message)
}

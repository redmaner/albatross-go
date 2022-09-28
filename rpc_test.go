package albatross

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ http.RoundTripper = (*testRoundtripper)(nil)

// this testRoundtripper is used to test the HTTP rpc client
type testRoundtripper struct {
	responseRecorder  *httptest.ResponseRecorder
	roundtripCallback func(r *http.Request) error
}

func (t *testRoundtripper) RoundTrip(r *http.Request) (*http.Response, error) {
	err := t.roundtripCallback(r)
	if err != nil {
		return nil, err
	}

	response := t.responseRecorder.Result()
	t.responseRecorder.Flush()

	return response, nil
}

func TestRpcCallOverHttpOk(t *testing.T) {
	recorder := httptest.NewRecorder()
	callback := func(r *http.Request) error {

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		expectedRawRequest := `{"jsonrpc":"2.0","id":1,"method":"getLatestBlock","Params":[]}`
		assert.Equal(t, strings.TrimSpace(string(data)), expectedRawRequest, "Request is invalid")

		expectedAuthHeader := `Basic dXNlcm5hbWU6cGFzc3dvcmQ=`
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, authHeader, expectedAuthHeader, "Auth header is invalid")

		return nil
	}

	roundtrip := &testRoundtripper{
		responseRecorder:  recorder,
		roundtripCallback: callback,
	}

	rpcClient := &HttpClient{
		client:   roundtrip,
		Url:      "https://test.albatross.example",
		UseAuth:  true,
		Username: "username",
		Password: "password",
	}

	mockResponse := `{"jsonrpc":"2.0","data":1234,"id":1}`
	recorder.WriteString(mockResponse)

	req := NewRPCRequestWithID("getLatestBlock", 1)
	resp, err := rpcClient.Call(req)
	if err != nil {
		t.Fatal(err)
	}

	blockNumber, err := GetObject[int](resp)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, blockNumber, 1234, "Latest block number invalid")
}

func TestRpcCallOverHttpRPCerror(t *testing.T) {
	recorder := httptest.NewRecorder()
	callback := func(r *http.Request) error {

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		expectedRawRequest := `{"jsonrpc":"2.0","id":1,"method":"getTransactionByHash","Params":["21cfba017cf06251846eb5085e52a2388b2c4c05bd1b155063358ea63f75ac53"]}`
		assert.Equal(t, strings.TrimSpace(string(data)), expectedRawRequest, "Request is invalid")

		expectedAuthHeader := `Basic dXNlcm5hbWU6cGFzc3dvcmQ=`
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, authHeader, expectedAuthHeader, "Auth header is invalid")

		return nil
	}

	roundtrip := &testRoundtripper{
		responseRecorder:  recorder,
		roundtripCallback: callback,
	}

	rpcClient := &HttpClient{
		client:   roundtrip,
		Url:      "https://test.albatross.example",
		UseAuth:  true,
		Username: "username",
		Password: "password",
	}

	mockResponse := `{"jsonrpc":"2.0","error":{"code":-32603,"message":"Internal error","data":"Multiple transactions found: 21cfba017cf06251846eb5085e52a2388b2c4c05bd1b155063358ea63f75ac53"},"id":1}`
	recorder.WriteString(mockResponse)

	req := NewRPCRequestWithID("getTransactionByHash", 1, "21cfba017cf06251846eb5085e52a2388b2c4c05bd1b155063358ea63f75ac53")
	resp, err := rpcClient.Call(req)
	if err != nil {
		t.Fatal(err)
	}

	_, err = GetObject[string](resp)
	expectedErr := `JSON-RPC Error -32603 - Internal error. Error data: Multiple transactions found: 21cfba017cf06251846eb5085e52a2388b2c4c05bd1b155063358ea63f75ac53`
	assert.Equal(t, err.Error(), expectedErr, "Returned error is invalid")
}

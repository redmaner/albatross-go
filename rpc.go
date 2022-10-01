package albatross

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var _ rpcClient = (*HttpClient)(nil)

type rpcClient interface {
	Call(*JsonRPCRequest) (*JsonRPCResponse, error)
}

func callAndConfirm(client rpcClient, req *JsonRPCRequest) error {
	_, err := client.Call(req)
	if err != nil {
		return err
	}

	return nil
}

func callAndUnwrap[T any](client rpcClient, req *JsonRPCRequest) (T, error) {
	rpcResp, err := client.Call(req)
	if err != nil {
		var emptyReturn T
		return emptyReturn, err
	}

	return UnwrapObject[T](rpcResp)
}

func callAndUnwrapToPointer[T any](client rpcClient, req *JsonRPCRequest) (*T, error) {
	data, err := callAndUnwrap[T](client, req)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

type HttpClient struct {
	client http.RoundTripper

	url      string
	useAuth  bool
	username string
	password string
}

// NewHttpClient returns a new HTTP RPC client to interact to the
// RPC server of a running albatross node
func NewHttpClient(url string) (*HttpClient, error) {
	if ok, err := verifyUrl(url); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("invalid url")
	}

	return &HttpClient{
		client: http.DefaultTransport,
		url:    url,
	}, nil
}

func (c *HttpClient) SetUseAuth(useAuth bool) *HttpClient {
	c.useAuth = useAuth
	return c
}

func (c *HttpClient) SetUsername(username string) *HttpClient {
	c.username = username
	return c
}

func (c *HttpClient) SetPassword(password string) *HttpClient {
	c.password = password
	return c
}

// Call executes an remote procedure call (RPC) using the given request
func (h *HttpClient) Call(r *JsonRPCRequest) (*JsonRPCResponse, error) {
	buf := bytes.NewBufferString("")
	if err := json.NewEncoder(buf).Encode(r); err != nil {
		return nil, err
	}

	body, err := h.send(buf)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var rpcResp JsonRPCResponse
	err = json.NewDecoder(body).Decode(&rpcResp)
	if err != nil {
		return nil, err
	}

	return &rpcResp, nil
}

// Batch executes a batch remote procedure call (RPC) using the given slice of requests
// Remember that according to the JSON-RPC spec the responses might be returned in a different order
// than the order in which the requests are provided.
func (h *HttpClient) Batch(r []*JsonRPCRequest) ([]*JsonRPCResponse, error) {
	buf := bytes.NewBufferString("")
	if err := json.NewEncoder(buf).Encode(r); err != nil {
		return nil, err
	}

	body, err := h.send(buf)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	rpcResp := make([]*JsonRPCResponse, len(r))
	err = json.NewDecoder(body).Decode(&rpcResp)
	if err != nil {
		return nil, err
	}

	return rpcResp, nil
}

func (h *HttpClient) send(body io.Reader) (io.ReadCloser, error) {
	httpRequest, err := http.NewRequest(http.MethodPost, h.url, body)
	if err != nil {
		return nil, err
	}

	h.setAuthHeader(httpRequest)

	httpResp, err := h.client.RoundTrip(httpRequest)
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("server responded with HTTP status code %d: %s", httpResp.StatusCode, string(data))
	}

	return httpResp.Body, nil
}

func (h *HttpClient) setAuthHeader(r *http.Request) {
	if !h.useAuth {
		return
	}

	bearerToken := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", h.username, h.password)))
	r.Header.Set("Authorization", fmt.Sprintf("Basic %s", bearerToken))
}

// GetBlockNumber retrieves the latest block number of the blockchain
func (h *HttpClient) GetBlockNumber() (blockNumber int, err error) {
	req := NewRPCRequest("getBlockNumber")

	return callAndUnwrap[int](h, req)
}

// GetBathhNumber retrieves the latest batch number of the blockchain
func (h *HttpClient) GetBatchNumber() (batchNumber int, err error) {
	req := NewRPCRequest("getBatchNumber")

	return callAndUnwrap[int](h, req)
}

// GetEpochNumber retrieves the latest epoch number of the blockchain
func (h *HttpClient) GetEpochNumber() (epochNumber int, err error) {
	req := NewRPCRequest("getEpochNumber")

	return callAndUnwrap[int](h, req)
}

// GetLatestBlock returns the latest block
func (h *HttpClient) GetLatestBlock(includeFullTransactions ...bool) (*Block, error) {
	params := []interface{}{}
	params = addOptionalParam(params, includeFullTransactions, false)
	req := NewRPCRequest("getLatestBlock", params...)

	return callAndUnwrapToPointer[Block](h, req)
}

// GetBlockByNumber retrieves the desired block by number
func (h *HttpClient) GetBlockByNumber(number int, includeFullTransactions ...bool) (*Block, error) {
	params := []interface{}{number}
	params = addOptionalParam(params, includeFullTransactions, false)
	req := NewRPCRequest("getBlockByNumber", params...)

	return callAndUnwrapToPointer[Block](h, req)
}

// GetBlockByHash retrieves the desired block by hash
func (h *HttpClient) GetBlockByHash(hash string, includeFullTransactions ...bool) (*Block, error) {
	params := []interface{}{hash}
	params = addOptionalParam(params, includeFullTransactions, false)
	req := NewRPCRequest("getBlockByHash", params...)

	return callAndUnwrapToPointer[Block](h, req)
}

// GetAccountByAddress returns the desired account by address
func (h *HttpClient) GetAccountByAddress(address string) (*Account, error) {
	req := NewRPCRequest("getAccountByAddress", address)

	return callAndUnwrapToPointer[Account](h, req)
}

// CreateAccount creates a new basic account on the Nimiq blockchain
func (h *HttpClient) CreateAccount(passphrase ...string) (*ReturnAccount, error) {
	params := addOptionalParam[string, interface{}]([]interface{}{}, passphrase, nil)

	req := NewRPCRequest("createAccount", params...)

	return callAndUnwrapToPointer[ReturnAccount](h, req)
}

// ImportAccountByRawKey import account on the node using the account's private key
func (h *HttpClient) ImportAccountByRawKey(rawKey string, passphrase ...string) error {
	params := []interface{}{rawKey}
	params = addOptionalParam[string, interface{}](params, passphrase, nil)

	req := NewRPCRequest("importRawKey", params...)

	return callAndConfirm(h, req)
}

// IsAccountImported returns whether the account is imported on the node
func (h *HttpClient) IsAccountImported(address string) (bool, error) {
	req := NewRPCRequest("isAccountImported", address)

	return callAndUnwrap[bool](h, req)
}

// LockAccount locks the given account on the node
func (h *HttpClient) LockAccount(address string) error {
	req := NewRPCRequest("lockAccount", address)

	return callAndConfirm(h, req)
}

// UnlockAccount unlocks the given account on the node
func (h *HttpClient) UnlockAccount(address string, passphrase ...string) error {
	params := []interface{}{address}
	params = addOptionalParam[string, interface{}](params, passphrase, nil)
	params = append(params, nil) // Param for duration, which is currently not supported by the server

	req := NewRPCRequest("unlockAccount", params...)

	return callAndConfirm(h, req)
}

// IsAccountImported returns whether the account is imported on the node
func (h *HttpClient) IsAccountUnlocked(address string) (bool, error) {
	req := NewRPCRequest("isAccountUnlocked", address)

	return callAndUnwrap[bool](h, req)
}

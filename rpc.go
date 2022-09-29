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
	"regexp"
)

type HttpClient struct {
	client http.RoundTripper

	url      string
	useAuth  bool
	username string
	password string
}

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

func verifyUrl(url string) (bool, error) {
	regex := `^(https|http|ws|wss):\/\/`
	return regexp.Match(regex, []byte(url))
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
	httpRequest, err := http.NewRequest(http.MethodGet, h.url, body)
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

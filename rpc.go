package albatross

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type HttpClient struct {
	client http.RoundTripper

	Url      string
	UseAuth  bool
	Username string
	Password string
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
	httpRequest, err := http.NewRequest(http.MethodGet, h.Url, body)
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
	if !h.UseAuth {
		return
	}

	bearerToken := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", h.Username, h.Password)))
	r.Header.Set("Authorization", fmt.Sprintf("Basic %s", bearerToken))
}

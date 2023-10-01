package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rotisserie/eris"
)

// HttpClient represents HTTP client.
type HttpClient struct {
	client   *http.Client
	baseURL  string
	basePath string
}

// NewAimApiClient creates new HTTP client for the aim api
func NewUserApiClient(baseURL string) *HttpClient {
	return &HttpClient{
		client:   &http.Client{},
		baseURL:  baseURL,
		basePath: "/users",
	}
}

// DoPostRequest do POST request.
func (c HttpClient) DoPostRequest(uri string, request interface{}, response interface{}) error {
	data, err := json.Marshal(request)
	if err != nil {
		return eris.Wrap(err, "error marshaling request")
	}
	return c.doRequest(http.MethodPost, uri, response, bytes.NewBuffer(data))
}

// DoPutRequest do PUT request.
func (c HttpClient) DoPutRequest(uri string, request interface{}, response interface{}) error {
	data, err := json.Marshal(request)
	if err != nil {
		return eris.Wrap(err, "error marshaling request")
	}
	return c.doRequest(http.MethodPut, uri, response, bytes.NewBuffer(data))
}

// DoGetRequest do GET request.
func (c HttpClient) DoGetRequest(uri string, response interface{}) error {
	return c.doRequest(http.MethodGet, uri, response, nil)
}

// DoDeleteRequest do DELETE request.
func (c HttpClient) DoDeleteRequest(uri string, response interface{}) error {
	return c.doRequest(http.MethodDelete, uri, response, nil)
}

// DoStreamRequest do stream request.
func (c HttpClient) DoStreamRequest(method, uri string, request interface{}) ([]byte, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, eris.Wrap(err, "error marshaling request")
	}

	// 1. create actual request object.
	str, err := StrReplace(
		fmt.Sprintf(
			"%s%s%s",
			c.baseURL,
			c.basePath,
			uri,
		),
		[]string{},
		[]interface{}{},
	)
	if err != nil {
		return nil, eris.Wrap(err, "error replacing string")
	}
	req, err := http.NewRequestWithContext(
		context.Background(),
		method,
		str,
		bytes.NewBuffer(data),
	)
	if err != nil {
		return nil, eris.Wrap(err, "error creating request")
	}
	req.Header.Set("Content-Type", "application/json")

	// 2. send request data.
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, eris.Wrap(err, "error doing request")
	}
	defer resp.Body.Close()

	// Check status code for errors
	if resp.StatusCode >= 400 {
		var errMsg string
		if resp.StatusCode >= 500 {
			errMsg = "server error"
		} else {
			errMsg = "client error"
		}
		return nil, eris.Errorf("HTTP request failed with status %d: %s", resp.StatusCode, errMsg)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, eris.Wrap(err, "error reading streaming response")
	}

	return body, nil
}

// doRequest do request of any http method
func (c HttpClient) doRequest(httpMethod string, uri string, response interface{}, requestBody io.Reader) error {
	// 1. create actual request object.
	str, err := StrReplace(
		fmt.Sprintf(
			"%s%s%s",
			c.baseURL,
			c.basePath,
			uri,
		),
		[]string{},
		[]interface{}{},
	)
	if err != nil {
		return eris.Wrap(err, "error replacing string")
	}
	req, err := http.NewRequestWithContext(
		context.Background(),
		httpMethod,
		str,
		requestBody,
	)
	if err != nil {
		return eris.Wrap(err, "error creating request")
	}
	req.Header.Set("Content-Type", "application/json")

	// 3. send request data.
	resp, err := c.client.Do(req)
	if err != nil {
		return eris.Wrap(err, "error doing request")
	}
	defer resp.Body.Close()

	// Check status code for errors
	if resp.StatusCode >= 400 {
		var errMsg string
		if resp.StatusCode >= 500 {
			errMsg = "server error"
		} else {
			errMsg = "client error"
		}
		return eris.Wrapf(err, "HTTP request failed with status %d: %s", resp.StatusCode, errMsg)
	}

	// 4. read and check response data.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return eris.Wrap(err, "error reading response data")
	}

	if err := json.Unmarshal(body, response); err != nil {
		return eris.Wrap(err, "error unmarshaling response data")
	}

	return nil
}

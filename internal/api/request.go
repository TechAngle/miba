package api

import (
	"encoding/json"
	"fmt"
	"io"
	"miba/internal/domain"
	"net/http"
	"net/url"
)

type ContentType string

const (
	URLEncoded ContentType = "application/x-www-form-urlencoded"
	JSON       ContentType = "application/json"
)

// SendRequest executes a request to endpoint.
//
// IMPORTANT: Returns response with OPENED body descriptor. It leads to memory leak if
// was not closed after execution of this method.
//
// FIX: endpoint must be concatenation of host and endpoint
func (c *MiAPIClient) MakeRequest(method, endpoint string, contentType ContentType, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", string(contentType))

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// buildUrl constructs URL with params (if provided)
func (c *MiAPIClient) buildUrl(endpoint string, params url.Values) string {
	u := *c.baseURL

	// NOTE: ignoring err because we use Mi Router hardcoded endpoint
	u.Path, _ = url.JoinPath(u.Path, endpoint)
	if params != nil {
		u.RawQuery = params.Encode()
	}

	return u.String()
}

// loginRequest performs a request to login endpoint with provided params.
func (c *MiAPIClient) loginRequest(params url.Values) (*domain.LoginResponse, error) {
	path := c.buildUrl(domain.LoginEndpoint, params)

	res, err := c.MakeRequest(http.MethodPost, path, URLEncoded, nil)
	if err != nil {
		return nil, fmt.Errorf("login req: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var data domain.LoginResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode body: %w", err)
	}

	return &data, nil
}

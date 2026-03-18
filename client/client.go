package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const jsonAPIContentType = "application/vnd.api+json"

// APIError is returned when the server responds with a non-2xx status.
type APIError struct {
	StatusCode int
	Errors     []ErrorDetail
}

// ErrorDetail is a single JSON:API error object.
type ErrorDetail struct {
	Status string `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail,omitempty"`
}

func (e *APIError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("api error %d: %s", e.StatusCode, e.Errors[0].Title)
	}
	return fmt.Sprintf("api error %d", e.StatusCode)
}

// Client communicates with the mock-fps server.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// New creates a Client with the default http.Client.
func New(baseURL string) *Client {
	return &Client{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		HTTPClient: http.DefaultClient,
	}
}

// NewWithClient creates a Client with a custom http.Client.
func NewWithClient(baseURL string, httpClient *http.Client) *Client {
	return &Client{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		HTTPClient: httpClient,
	}
}

// doJSON performs an HTTP request with JSON:API content type.
// If body is non-nil it is marshalled as the request body.
// If result is non-nil the response body is unmarshalled into it.
func (c *Client) doJSON(method, path string, body, result any) error {
	_, err := c.doJSONWithStatus(method, path, body, result)
	return err
}

// doJSONWithStatus is like doJSON but also returns the HTTP status code.
func (c *Client) doJSONWithStatus(method, path string, body, result any) (int, error) {
	return c.doRequest(method, path, jsonAPIContentType, body, result)
}

// doPlainJSON performs a request with application/json content type (for admin endpoints).
func (c *Client) doPlainJSON(method, path string, body, result any) error {
	_, err := c.doRequest(method, path, "application/json", body, result)
	return err
}

func (c *Client) doRequest(method, path, contentType string, body, result any) (int, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return 0, fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, reqBody)
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Accept", contentType)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		var errResp struct {
			Errors []ErrorDetail `json:"errors"`
		}
		if json.Unmarshal(respBody, &errResp) == nil {
			apiErr.Errors = errResp.Errors
		}
		return resp.StatusCode, apiErr
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return resp.StatusCode, fmt.Errorf("unmarshal response: %w", err)
		}
	}

	return resp.StatusCode, nil
}

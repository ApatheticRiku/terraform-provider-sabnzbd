// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is the SABnzbd API client.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new SABnzbd API client.
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// APIError represents an error returned by the SABnzbd API.
type APIError struct {
	Message string
}

func (e *APIError) Error() string {
	return e.Message
}

// doRequest performs an API request and decodes the JSON response.
func (c *Client) doRequest(ctx context.Context, params url.Values, result interface{}) error {
	params.Set("apikey", c.apiKey)
	params.Set("output", "json")

	reqURL := fmt.Sprintf("%s/api?%s", c.baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	// Check for API errors in response.
	var errorResp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Error != "" {
		return &APIError{Message: errorResp.Error}
	}

	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}

	return nil
}

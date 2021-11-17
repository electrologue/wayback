// Package json Wayback Availability JSON API.
// https://archive.org/help/wayback_api.php
package json

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

// BaseURL base URL of the API endpoint.
const BaseURL = "http://archive.org/wayback/available"

// Client is an API JSON client.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// New creates a new Client.
func New() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		baseURL:    BaseURL,
	}
}

// Available test to see if a given url is archived and currently accessible in the Wayback Machine.
func (c Client) Available(ctx context.Context, host, timestamp string) (*APIResponse, error) {
	endpoint, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}

	query := endpoint.Query()
	query.Set("url", host)
	if timestamp != "" {
		query.Set("timestamp", timestamp)
	}
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse APIResponse
	err = json.Unmarshal(data, &apiResponse)
	if err != nil {
		return nil, err
	}

	return &apiResponse, nil
}

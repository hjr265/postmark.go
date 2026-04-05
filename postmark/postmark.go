// Package postmark provides a client for the Postmark API.
package postmark

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	DefaultHost = "api.postmarkapp.com"
)

// Option is a functional option for configuring a Client.
type Option func(*Client)

// WithHost sets a custom API host. If not set, DefaultHost is used.
func WithHost(host string) Option {
	return func(c *Client) {
		c.host = host
	}
}

// Client is a Postmark API client.
type Client struct {
	apiKey string
	host   string
}

// New creates a new Postmark API client with the given API key.
func New(apiKey string, opts ...Option) *Client {
	c := &Client{apiKey: apiKey}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) endpoint(path string) *url.URL {
	url := &url.URL{}
	url.Scheme = "https"

	if c.host == "" {
		url.Host = DefaultHost
	} else {
		url.Host = c.host
	}

	url.Path = path

	return url
}

// do sends a POST request to the given API path, encoding body as JSON and
// decoding the JSON response into result.
func (c *Client) do(path string, body, result interface{}) error {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.endpoint(path).String(), &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Server-Token", c.apiKey)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(result)
}

// get sends a GET request to the given API path and decodes the JSON response
// into result.
func (c *Client) get(path string, result interface{}) error {
	req, err := http.NewRequest("GET", c.endpoint(path).String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Postmark-Server-Token", c.apiKey)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(result)
}

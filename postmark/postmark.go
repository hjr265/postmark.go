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

// Client is a Postmark API client.
type Client struct {
	ApiKey string
	Secure bool

	Host string // Host for the API endpoints, DefaultHost if ""
}

func (c *Client) endpoint(path string) *url.URL {
	url := &url.URL{}
	if c.Secure {
		url.Scheme = "https"
	} else {
		url.Scheme = "http"
	}

	if c.Host == "" {
		url.Host = DefaultHost
	} else {
		url.Host = c.Host
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
	req.Header.Set("X-Postmark-Server-Token", c.ApiKey)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(result)
}

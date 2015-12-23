// Package postmark provides a convenient wrapper for the Postmark API
package postmark

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
)

const (
	DefaultHost = "api.postmarkapp.com"
)

type Message struct {
	From          *mail.Address
	To            []*mail.Address
	Cc            []*mail.Address
	Bcc           []*mail.Address
	Subject       string
	Tag           string
	HtmlBody      io.Reader
	TextBody      io.Reader
	TemplateId    int
	TemplateModel map[string]interface{}
	ReplyTo       *mail.Address
	Headers       mail.Header
	Attachments   []Attachment
}

func (m *Message) MarshalJSON() ([]byte, error) {
	doc := &struct {
		From          string
		To            string
		Cc            string
		Bcc           string
		Subject       string `json:",omitempty"`
		Tag           string
		HtmlBody      string `json:",omitempty"`
		TextBody      string `json:",omitempty"`
		TemplateId    int    `json:",omitempty"`
		TemplateModel map[string]interface{}
		ReplyTo       string
		Headers       []map[string]string
		Attachments   []Attachment `json:"omitempty"`
	}{}

	doc.From = m.From.String()
	to := []string{}
	for _, addr := range m.To {
		to = append(to, addr.String())
	}
	doc.To = strings.Join(to, ", ")
	cc := []string{}
	for _, addr := range m.Cc {
		cc = append(cc, addr.String())
	}
	doc.Cc = strings.Join(cc, ", ")
	bcc := []string{}
	for _, addr := range m.Bcc {
		bcc = append(bcc, addr.String())
	}
	doc.Bcc = strings.Join(bcc, ", ")
	doc.Subject = m.Subject
	doc.Tag = m.Tag
	if m.HtmlBody != nil {
		htmlBody, err := ioutil.ReadAll(m.HtmlBody)
		if err != nil {
			return nil, err
		}
		doc.HtmlBody = string(htmlBody)
	}
	if m.TextBody != nil {
		textBody, err := ioutil.ReadAll(m.TextBody)
		if err != nil {
			return nil, err
		}
		doc.TextBody = string(textBody)
	}
	doc.TemplateId = m.TemplateId
	doc.TemplateModel = m.TemplateModel
	if m.ReplyTo != nil {
		doc.ReplyTo = m.ReplyTo.String()
	}
	headers := []map[string]string{}
	for k, vs := range m.Headers {
		for _, v := range vs {
			headers = append(headers, map[string]string{
				"Name":  k,
				"Value": v,
			})
		}
	}
	doc.Headers = headers
	doc.Attachments = m.Attachments

	return json.Marshal(doc)
}

type Attachment struct {
	Name        string
	Content     io.Reader
	ContentType string
}

func (a *Attachment) MarshalJSON() ([]byte, error) {
	doc := &struct {
		Name        string
		Content     string
		ContentType string
	}{}

	doc.Name = a.Name
	content, err := ioutil.ReadAll(a.Content)
	if err != nil {
		return nil, err
	}
	doc.Content = base64.StdEncoding.EncodeToString(content)
	doc.ContentType = a.ContentType

	return json.Marshal(doc)
}

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

// Send sends a single message
func (c *Client) Send(msg *Message) (*Result, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(msg)
	if err != nil {
		return nil, err
	}

	url := c.endpoint("email")
	if msg.TemplateId != 0 {
		url = c.endpoint("email/withTemplate")
	}

	req, err := http.NewRequest("POST", url.String(), &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Server-Token", c.ApiKey)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}

	res := &Result{}
	json.NewDecoder(resp.Body).Decode(res)
	return res, nil
}

// SendBatch sends multiple messages using the batch API
func (c *Client) SendBatch(msg []*Message) ([]*Result, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(msg)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.endpoint("email/batch").String(), &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Server-Token", c.ApiKey)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}

	res := []*Result{}
	json.NewDecoder(resp.Body).Decode(res)
	return res, nil
}

type Result struct {
	ErrorCode   int
	Message     string
	MessageID   string
	SubmittedAt string
	To          string
}

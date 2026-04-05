package postmark

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/mail"
	"strings"
)

// Message represents an email message.
//
// Deprecated: Use EmailRequest instead.
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
		Attachments   []Attachment `json:",omitempty"`
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
		htmlBody, err := io.ReadAll(m.HtmlBody)
		if err != nil {
			return nil, err
		}
		doc.HtmlBody = string(htmlBody)
	}
	if m.TextBody != nil {
		textBody, err := io.ReadAll(m.TextBody)
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

// Attachment represents a file attachment.
//
// Deprecated: Use EmailAttachment instead.
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
	content, err := io.ReadAll(a.Content)
	if err != nil {
		return nil, err
	}
	doc.Content = base64.StdEncoding.EncodeToString(content)
	doc.ContentType = a.ContentType

	return json.Marshal(doc)
}

// Result represents a response from the Postmark API.
//
// Deprecated: Use EmailResponse instead.
type Result struct {
	ErrorCode   int
	Message     string
	MessageID   string
	SubmittedAt string
	To          string
}

// Send sends a single message.
//
// Deprecated: Use SendEmail instead.
func (c *Client) Send(msg *Message) (*Result, error) {
	if msg.TemplateId != 0 {
		return c.sendWithTemplate(msg)
	}

	req, err := messageToEmailRequest(msg)
	if err != nil {
		return nil, err
	}

	resp, err := c.SendEmail(req)
	if err != nil {
		return nil, err
	}

	return emailResponseToResult(resp), nil
}

// SendBatch sends multiple messages using the batch API.
//
// Deprecated: Use SendEmailBatch instead.
func (c *Client) SendBatch(msgs []*Message) ([]*Result, error) {
	reqs := make([]*EmailRequest, len(msgs))
	for i, msg := range msgs {
		req, err := messageToEmailRequest(msg)
		if err != nil {
			return nil, err
		}
		reqs[i] = req
	}

	resps, err := c.SendEmailBatch(reqs)
	if err != nil {
		return nil, err
	}

	results := make([]*Result, len(resps))
	for i, resp := range resps {
		results[i] = emailResponseToResult(resp)
	}
	return results, nil
}

func (c *Client) sendWithTemplate(msg *Message) (*Result, error) {
	res := &Result{}
	err := c.do("email/withTemplate", msg, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func messageToEmailRequest(msg *Message) (*EmailRequest, error) {
	req := &EmailRequest{}

	req.From = msg.From.String()

	to := make([]string, len(msg.To))
	for i, addr := range msg.To {
		to[i] = addr.String()
	}
	req.To = strings.Join(to, ", ")

	cc := make([]string, len(msg.Cc))
	for i, addr := range msg.Cc {
		cc[i] = addr.String()
	}
	req.Cc = strings.Join(cc, ", ")

	bcc := make([]string, len(msg.Bcc))
	for i, addr := range msg.Bcc {
		bcc[i] = addr.String()
	}
	req.Bcc = strings.Join(bcc, ", ")

	req.Subject = msg.Subject
	req.Tag = msg.Tag

	if msg.HtmlBody != nil {
		htmlBody, err := io.ReadAll(msg.HtmlBody)
		if err != nil {
			return nil, err
		}
		req.HtmlBody = string(htmlBody)
	}

	if msg.TextBody != nil {
		textBody, err := io.ReadAll(msg.TextBody)
		if err != nil {
			return nil, err
		}
		req.TextBody = string(textBody)
	}

	if msg.ReplyTo != nil {
		req.ReplyTo = msg.ReplyTo.String()
	}

	for k, vs := range msg.Headers {
		for _, v := range vs {
			req.Headers = append(req.Headers, Header{Name: k, Value: v})
		}
	}

	for _, att := range msg.Attachments {
		content, err := io.ReadAll(att.Content)
		if err != nil {
			return nil, err
		}
		req.Attachments = append(req.Attachments, EmailAttachment{
			Name:        att.Name,
			Content:     base64.StdEncoding.EncodeToString(content),
			ContentType: att.ContentType,
		})
	}

	return req, nil
}

func emailResponseToResult(resp *EmailResponse) *Result {
	return &Result{
		ErrorCode:   resp.ErrorCode,
		Message:     resp.Message,
		MessageID:   resp.MessageID,
		SubmittedAt: resp.SubmittedAt,
		To:          resp.To,
	}
}

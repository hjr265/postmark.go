package postmark

// Header represents a custom email header.
type Header struct {
	Name  string
	Value string
}

// EmailAttachment represents a file attachment for an email.
type EmailAttachment struct {
	Name        string
	Content     string // base64-encoded
	ContentType string
}

// EmailRequest represents a request to send an email via the Postmark API.
type EmailRequest struct {
	From          string            `json:",omitempty"`
	To            string            `json:",omitempty"`
	Cc            string            `json:",omitempty"`
	Bcc           string            `json:",omitempty"`
	Subject       string            `json:",omitempty"`
	Tag           string            `json:",omitempty"`
	HtmlBody      string            `json:",omitempty"`
	TextBody      string            `json:",omitempty"`
	ReplyTo       string            `json:",omitempty"`
	Headers       []Header          `json:",omitempty"`
	TrackOpens    bool              `json:",omitempty"`
	TrackLinks    string            `json:",omitempty"`
	Metadata      map[string]string `json:",omitempty"`
	Attachments   []EmailAttachment `json:",omitempty"`
	MessageStream string            `json:",omitempty"`
}

// EmailResponse represents a response from the Postmark API after sending an
// email.
type EmailResponse struct {
	To          string
	SubmittedAt string
	MessageID   string
	ErrorCode   int
	Message     string
}

// SendEmail sends a single email via the Postmark API.
func (c *Client) SendEmail(req *EmailRequest) (*EmailResponse, error) {
	res := &EmailResponse{}
	err := c.do("email", req, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SendEmailBatch sends multiple emails in a single API call.
func (c *Client) SendEmailBatch(req []*EmailRequest) ([]*EmailResponse, error) {
	res := []*EmailResponse{}
	err := c.do("email/batch", req, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

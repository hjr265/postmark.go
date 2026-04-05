package postmark

// BulkMessage represents a recipient message within a bulk email request.
type BulkMessage struct {
	To            string            `json:",omitempty"`
	Cc            string            `json:",omitempty"`
	Bcc           string            `json:",omitempty"`
	TemplateModel map[string]string `json:",omitempty"`
	Metadata      map[string]string `json:",omitempty"`
	Headers       []Header          `json:",omitempty"`
}

// BulkEmailRequest represents a request to send bulk emails via the Postmark
// API.
type BulkEmailRequest struct {
	From          string            `json:",omitempty"`
	ReplyTo       string            `json:",omitempty"`
	Subject       string            `json:",omitempty"`
	HtmlBody      string            `json:",omitempty"`
	TextBody      string            `json:",omitempty"`
	TemplateId    int               `json:",omitempty"`
	TemplateAlias string            `json:",omitempty"`
	InlineCss     bool              `json:",omitempty"`
	Tag           string            `json:",omitempty"`
	Metadata      map[string]string `json:",omitempty"`
	MessageStream string            `json:",omitempty"`
	TrackOpens    bool              `json:",omitempty"`
	TrackLinks    string            `json:",omitempty"`
	Attachments   []EmailAttachment `json:",omitempty"`
	Headers       []Header          `json:",omitempty"`
	Messages      []BulkMessage     `json:",omitempty"`
}

// BulkEmailResponse represents a response from the Postmark API after
// submitting a bulk email request.
type BulkEmailResponse struct {
	ID          string
	Status      string
	SubmittedAt string
}

// BulkEmailStatusResponse represents the status of a bulk email request.
type BulkEmailStatusResponse struct {
	Id                  string
	Status              string
	SubmittedAt         string
	TotalMessages       int
	PercentageCompleted float64
	Subject             string
}

// SendBulkEmail submits a bulk email request via the Postmark API.
func (c *Client) SendBulkEmail(req *BulkEmailRequest) (*BulkEmailResponse, error) {
	res := &BulkEmailResponse{}
	err := c.do("email/bulk", req, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetBulkEmailStatus retrieves the status of a bulk email request.
func (c *Client) GetBulkEmailStatus(id string) (*BulkEmailStatusResponse, error) {
	res := &BulkEmailStatusResponse{}
	err := c.get("email/bulk/"+id, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

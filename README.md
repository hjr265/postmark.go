# Postmark.go

Postmark.go is a Go client library for the [Postmark API](https://postmarkapp.com/developer).

## Installation

    go get github.com/hjr265/postmark.go/postmark

## Usage

### Send a Single Email

```go
c := postmark.New("YOUR-API-KEY")

resp, err := c.SendEmail(&postmark.EmailRequest{
    From:     "sender@example.com",
    To:       "recipient@example.com",
    Subject:  "Hello",
    TextBody: "Hello from Postmark!",
})
```

### Send Batch Emails

```go
resp, err := c.SendEmailBatch([]*postmark.EmailRequest{
    {
        From:     "sender@example.com",
        To:       "recipient1@example.com",
        Subject:  "Hello",
        TextBody: "Hello from Postmark!",
    },
    {
        From:     "sender@example.com",
        To:       "recipient2@example.com",
        Subject:  "Hello",
        TextBody: "Hello from Postmark!",
    },
})
```

### Send Bulk Email

```go
resp, err := c.SendBulkEmail(&postmark.BulkEmailRequest{
    From:          "sender@example.com",
    Subject:       "Hello {{FirstName}}",
    TextBody:      "Hello {{FirstName}}!",
    MessageStream: "broadcast",
    Messages: []postmark.BulkMessage{
        {
            To:            "recipient1@example.com",
            TemplateModel: map[string]string{"FirstName": "Alice"},
        },
        {
            To:            "recipient2@example.com",
            TemplateModel: map[string]string{"FirstName": "Bob"},
        },
    },
})
```

### Check Bulk Email Status

```go
status, err := c.GetBulkEmailStatus("bulk-request-id")
```

### Custom Host

```go
c := postmark.New("YOUR-API-KEY", postmark.WithHost("custom.api.host"))
```

## Documentation

- [API Reference](https://pkg.go.dev/github.com/hjr265/postmark.go/postmark)

## Contributing

Contributions are welcome.

## License

Postmark.go is available under the [BSD (3-Clause) License](https://opensource.org/licenses/BSD-3-Clause).

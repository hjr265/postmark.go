package postmark_test

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/hjr265/postmark.go/postmark"
)

func ExampleClient() {
	c := postmark.Client{
		ApiKey: "YOUR-API-KEY",
		Secure: true,
	}

	res, err := c.Send(&postmark.Message{
		From: &mail.Address{
			Name:    "SENDER-NAME",
			Address: "SENDER-EMAIL",
		},
		To: []*mail.Address{
			{
				Name:    "RECIPIENT NAME",
				Address: "RECIPIENT EMAIL",
			},
		},
		Subject:  "SUBJECT",
		TextBody: strings.NewReader("MESSAGE-BODY-AS-TEXT"),
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", res)
}

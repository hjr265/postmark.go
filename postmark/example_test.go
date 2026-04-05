package postmark_test

import (
	"fmt"

	"github.com/hjr265/postmark.go/postmark"
)

func ExampleClient() {
	c := postmark.New("YOUR-API-KEY")

	res, err := c.SendEmail(&postmark.EmailRequest{
		From:     "SENDER-NAME <SENDER-EMAIL>",
		To:       "RECIPIENT NAME <RECIPIENT EMAIL>",
		Subject:  "SUBJECT",
		TextBody: "MESSAGE-BODY-AS-TEXT",
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", res)
}

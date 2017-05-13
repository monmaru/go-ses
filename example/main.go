package main

import (
	"fmt"
	"log"
	"os"

	"github.com/monmaru/go-ses"
)

var region, accessKeyID, secretAccessKey string

func init() {
	region = os.Getenv("AWS_REGION")
	accessKeyID = os.Getenv("AWS_ACCESS_KEY")
	secretAccessKey = os.Getenv("AWS_SECRET_KEY")
}

func main() {
	setting := ses.AWSSetting{
		Region:          region,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}
	c := ses.NewClient(setting, nil)

	from := ""
	to := ""
	params := &ses.EmailParams{
		From:     from,
		To:       to,
		Subject:  "Hello world!!",
		BodyText: "This is a test mail.",
	}

	out, err := c.SendEmail(params)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
}

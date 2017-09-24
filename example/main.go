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
	setting := ses.Config{
		Region:          region,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}
	c := ses.NewClient(setting, nil)

	from := ""
	to := []string{""}
	mail := ses.Mail{
		From:     from,
		To:       to,
		Subject:  "Hello world!!",
		BodyText: "This is a test mail.",
	}

	msgID, err := c.SendEmail(mail)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(msgID)
}

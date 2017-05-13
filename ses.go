package ses

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// AWSSetting ...
type AWSSetting struct {
	Region, AccessKeyID, SecretAccessKey string
}

// Client ...
type Client struct {
	setting    AWSSetting
	httpClient *http.Client
}

// NewClient ...
func NewClient(awsSetting AWSSetting, httpClient *http.Client) *Client {
	return &Client{
		setting:    awsSetting,
		httpClient: httpClient,
	}
}

func (c *Client) newSESSession() *ses.SES {
	config := &aws.Config{
		Region: aws.String(c.setting.Region),
		Credentials: credentials.NewStaticCredentials(
			c.setting.AccessKeyID,
			c.setting.SecretAccessKey,
			""),
	}

	if c.httpClient != nil {
		config.HTTPClient = c.httpClient
	}

	awsSession := session.New(config)
	return ses.New(awsSession)
}

// SendEmail ...
func (c *Client) SendEmail(params EmailParams) (string, error) {
	email := params.toSendEmailInput()
	session := c.newSESSession()
	out, err := session.SendEmail(email)
	if err != nil {
		return "", err
	}

	return out.GoString(), nil
}

// SendRawEmail ...
func (c *Client) SendRawEmail(raw []byte) (string, error) {
	rowEmailInput := &ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{
			Data: []byte(raw),
		},
	}

	session := c.newSESSession()
	out, err := session.SendRawEmail(rowEmailInput)
	if err != nil {
		return "", err
	}

	return out.GoString(), nil
}

// EmailParams ...
type EmailParams struct {
	from, to, subject, bodyText, bodyHTML string
}

func (e *EmailParams) toSendEmailInput() *ses.SendEmailInput {
	message := &ses.Message{
		Body: &ses.Body{
			Text: &ses.Content{
				Data: aws.String(e.bodyText),
			},
			Html: &ses.Content{
				Data: aws.String(e.bodyHTML),
			},
		},
		Subject: &ses.Content{
			Data: aws.String(e.subject),
		},
	}

	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(e.to)},
		},
		Message: message,
		Source:  aws.String(e.from),
		ReplyToAddresses: []*string{
			aws.String(e.from),
		},
	}
}

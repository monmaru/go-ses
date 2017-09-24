package ses

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// Config ...
type Config struct {
	Region, AccessKeyID, SecretAccessKey string
}

// Mail ...
type Mail struct {
	From     string   `json:"from"`
	To       []string `json:"to"`
	Subject  string   `json:"subject"`
	BodyText string   `json:"bodyText"`
	BodyHTML string   `json:"bodyHtml"`
}

// Option ...
type Option func(*Client)

// HTTPClient ...
func HTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// Client ...
type Client struct {
	config     Config
	httpClient *http.Client
}

// NewClient ...
func NewClient(config Config, opts ...Option) *Client {
	c := &Client{
		config:     config,
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) newSESSession() *ses.SES {
	config := &aws.Config{
		Region: aws.String(c.config.Region),
		Credentials: credentials.NewStaticCredentials(
			c.config.AccessKeyID,
			c.config.SecretAccessKey,
			""),
	}

	if c.httpClient != nil {
		config.HTTPClient = c.httpClient
	}

	awsSession := session.New(config)
	return ses.New(awsSession)
}

// SendEmail ...
func (c *Client) SendEmail(mail Mail) (string, error) {
	email := mail.buildEmailInput()
	out, err := c.newSESSession().SendEmail(email)
	if err != nil {
		return "", err
	}

	return *out.MessageId, nil
}

// SendRawEmail ...
func (c *Client) SendRawEmail(rawText string) (string, error) {
	rawEmailInput := &ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{
			Data: []byte(rawText),
		},
	}

	out, err := c.newSESSession().SendRawEmail(rawEmailInput)
	if err != nil {
		return "", err
	}

	return *out.MessageId, nil
}

func (m *Mail) buildEmailInput() *ses.SendEmailInput {
	message := &ses.Message{
		Body: &ses.Body{
			Text: &ses.Content{
				Data: aws.String(m.BodyText),
			},
		},
		Subject: &ses.Content{
			Data: aws.String(m.Subject),
		},
	}

	if len(m.BodyHTML) > 0 {
		message.Body.Html = &ses.Content{
			Data: aws.String(m.BodyHTML),
		}
	}

	var toAddresses []*string
	for _, t := range m.To {
		toAddresses = append(toAddresses, &t)
	}

	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: toAddresses,
		},
		Message: message,
		Source:  aws.String(m.From),
		ReplyToAddresses: []*string{
			aws.String(m.From),
		},
	}
}

package email

import (
	"net/http"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// FromAddress => Default Sendgrid from address
const FromAddress = "test@test.com"

// SendGridAPIUrl => Sennd API url
const SendGridAPIUrl = "https://api.sendgrid.com"

// SendGridAPIEndpoint => Send grid API endpoint
const SendGridAPIEndpoint = "/v3/mail/send"

// SendGridDispatcher , extending default dispatcher
type SendGridDispatcher struct {
	to      string
	from    string
	msg     string
	subject string
	APIKey  string
}

// NewSendGridDispatcher => returns a new send grid dispatcher instance
func NewSendGridDispatcher(to, subject, msg, APIKey string) *SendGridDispatcher {
	return &SendGridDispatcher{
		to:      to,
		from:    FromAddress,
		msg:     msg,
		subject: subject,
		APIKey:  APIKey,
	}
}

// Dispatch => Create payload and calls sendgrid API with given payload (Create & Send Email)
func (sd *SendGridDispatcher) Dispatch() (bool, error) {
	body := GetHTMLBody(sd.to, sd.from, sd.subject, sd.msg)
	return SendMail(body, sd.APIKey)
}

// GetHTMLBody => Create mail body from Sendgrid
func GetHTMLBody(toAddress, fromAddress, subject, contentHTML string) []byte {
	from := mail.NewEmail("Test User", fromAddress)
	to := mail.NewEmail("", toAddress)
	content := mail.NewContent("text/html", contentHTML)
	m := mail.NewV3MailInit(from, subject, to, content)

	return mail.GetRequestBody(m)
}

// SendMail => calls SendGrid API (Sends Mail)
func SendMail(body []byte, apiKey string) (bool, error) {
	request := sendgrid.GetRequest(apiKey, SendGridAPIEndpoint, SendGridAPIUrl)
	request.Method = http.MethodPost
	request.Body = body
	response, err := sendgrid.API(request)
	if err != nil {
		return false, err
	}

	// https://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html
	// Sendgrid status codes.
	return response.StatusCode == http.StatusAccepted ||
		response.StatusCode == http.StatusOK, nil
}

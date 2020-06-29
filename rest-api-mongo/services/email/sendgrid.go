package email

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"net/http"
)

const FromAddress = "test@test.com"
const SendGridApiUrl = "https://api.sendgrid.com"
const SendGridApiEndPoint = "/v3/mail/send"

func GetHtmlBody(toAddress, subject, contentHtml string) []byte {
	from := mail.NewEmail("Test User", FromAddress)
	to := mail.NewEmail("", toAddress)
	content := mail.NewContent("text/html", contentHtml)
	m := mail.NewV3MailInit(from, subject, to, content)

	return mail.GetRequestBody(m)
}

func SendMail(body []byte, apiKey string) (bool, error) {
	request := sendgrid.GetRequest(apiKey, SendGridApiEndPoint, SendGridApiUrl)
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

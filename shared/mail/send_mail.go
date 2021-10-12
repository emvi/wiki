package mail

import (
	"emviwiki/shared/config"
	"fmt"
	"github.com/aws/aws-sdk-go/service/sesv2"
	"github.com/emvi/logbuch"
	"github.com/jordan-wright/email"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"net/smtp"
	"net/textproto"
)

const (
	mailContentType = "text/html"
)

// Selects the mail provider by environment variable.
// If the SendGrid API key is set, SendGridMail is used, SMTP otherwise.
func SelectMailSender() Sender {
	if config.Get().Mail.SendGridAPIKey != "" {
		return SendGridMail
	} else if config.Get().Mail.AmazonSESAPIKey != "" {
		return SendAmazonSESMail
	}

	return SendMail
}

// Sender sends mails with given subject, HTML body from given email address to given email addresses.
// This is used to use a real or mock implementation to send mails.
type Sender func(string, string, string, ...string) error

// SendMail implements the Sender interface and sends real mails through SMTP.
// This is the default send method.
func SendMail(subject, msgHTML, from string, to ...string) error {
	from, to = selectToFromMail(from, to)
	m := &email.Email{
		To:      to,
		From:    from,
		Subject: subject,
		HTML:    []byte(msgHTML),
		Headers: textproto.MIMEHeader{},
	}
	auth := smtp.PlainAuth("", mailSmtpUser, mailSmtpPassword, mailSmtpServer)

	if err := m.Send(fmt.Sprintf("%s:%d", mailSmtpServer, mailSmtpPort), auth); err != nil {
		logbuch.Error("Error sending mail", logbuch.Fields{"err": err, "subject": subject, "from": from, "to": to})
		return err
	}

	return nil
}

// SendAmazonSESMail implements the Sender interface and sends real mails through AWS SES.
func SendAmazonSESMail(subject, msgHTML, from string, to ...string) error {
	from, to = selectToFromMail(from, to)
	receiver := make([]*string, 0, len(to))

	for _, toMail := range to {
		receiver = append(receiver, &toMail)
	}

	_, err := awsSESv2.SendEmail(&sesv2.SendEmailInput{
		FromEmailAddress: stringPtr(from),
		Destination: &sesv2.Destination{
			ToAddresses: receiver,
		},
		Content: &sesv2.EmailContent{
			Simple: &sesv2.Message{
				Subject: &sesv2.Content{
					Charset: stringPtr("UTF-8"),
					Data:    stringPtr(subject),
				},
				Body: &sesv2.Body{
					Html: &sesv2.Content{
						Charset: stringPtr("UTF-8"),
						Data:    stringPtr(msgHTML),
					},
				},
			},
		},
	})

	if err != nil {
		logbuch.Error("Error sending mail", logbuch.Fields{"err": err})
		return err
	}

	return nil
}

// SendGridMail implements the Sender interface and sends real mails through SendGrids API.
func SendGridMail(subject, msgHTML, from string, to ...string) error {
	from, to = selectToFromMail(from, to)
	client := sendgrid.NewSendClient(mailSendGrindAPIKey)
	fromMail := mail.NewEmail("", from)
	content := mail.NewContent(mailContentType, msgHTML)

	for _, toMail := range to {
		receiver := mail.NewEmail("", toMail)
		payload := mail.NewV3MailInit(fromMail, subject, receiver, content)

		if _, err := client.Send(payload); err != nil {
			logbuch.Error("Error sending mail", logbuch.Fields{"err": err})
			return err
		}
	}

	return nil
}

func selectToFromMail(from string, to []string) (string, []string) {
	// use configured default sender when no receiver is passed
	if len(to) == 0 {
		to = append(to, from)
		from = mailFrom
	}

	logbuch.Debug("Mail receiver/sender", logbuch.Fields{"receiver": to, "sender": from})
	return from, to
}

func stringPtr(str string) *string {
	return &str
}

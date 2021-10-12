package mail

import (
	"emviwiki/shared/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sesv2"
	"github.com/emvi/logbuch"
)

var (
	mailFrom, mailSmtpServer, mailSmtpUser, mailSmtpPassword, mailSendGrindAPIKey string
	mailSmtpPort                                                                  int
	awsSession                                                                    *session.Session
	awsSESv2                                                                      *sesv2.SESV2
)

func LoadConfig() {
	c := config.Get()
	mailFrom = c.Mail.Sender
	mailSmtpServer = c.Mail.SMTP.Server
	mailSmtpPort = c.Mail.SMTP.Port
	mailSmtpUser = c.Mail.SMTP.User
	mailSmtpPassword = c.Mail.SMTP.Password
	mailSendGrindAPIKey = c.Mail.SendGridAPIKey

	if config.Get().Mail.AmazonSESAPIKey != "" {
		createAWSSession()
	}
}

func createAWSSession() {
	logbuch.Info("Creating AWS SES client...")
	var err error
	awsSession, err = session.NewSession(&aws.Config{
		Region:      aws.String(config.Get().Mail.AmazonSESRegion),
		Credentials: credentials.NewStaticCredentials(config.Get().Mail.AmazonSESAPIKey, config.Get().Mail.AmazonSESAPISecret, ""),
	})

	if err != nil {
		logbuch.Fatal("Error creating AWS session", logbuch.Fields{"err": err})
	}

	awsSESv2 = sesv2.New(awsSession)
	logbuch.Info("AWS SES client created")
}

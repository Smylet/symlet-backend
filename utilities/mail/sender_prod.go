package mail

import (
	"bytes"
	"fmt"

	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"gopkg.in/gomail.v2"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type SESEmailSender struct {
	fromEmailAddress string
	session          *session.Session
	config           utils.Config
}

func NewSESEmailSender(fromEmailAddress string, session *session.Session, config utils.Config) EmailSender {
	return &SESEmailSender{
		fromEmailAddress: fromEmailAddress,
		config:           config,
		session:          session,
	}
}

func (sender *SESEmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	// Check if we're in development mode
	if sender.config.Environment == "development" {
		fmt.Println("Sending email in development mode")
		return SendEmailDev(subject, content, to, cc, bcc, attachFiles, sender.fromEmailAddress)

	}

	var err error

	// Create raw message
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", sender.fromEmailAddress, "")
	msg.SetHeader("To", to...)
	if len(cc) > 0 {
		msg.SetHeader("Cc", cc...)
	}
	if len(bcc) > 0 {
		msg.SetHeader("Bcc", bcc...)
	}
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", content)

	// Attach files if any
	for _, f := range attachFiles {
		msg.Attach(f)
	}

	// Create a buffer to store raw data
	var emailRaw bytes.Buffer
	msg.WriteTo(&emailRaw)

	// Create a new raw message
	rawMessage := ses.RawMessage{Data: emailRaw.Bytes()}
	destinations := append(to, append(cc, bcc...)...)

	// Create SES service client
	svc := ses.New(sender.session)

	// Send raw email
	_, err = svc.SendRawEmail(&ses.SendRawEmailInput{
		Source:       &sender.fromEmailAddress,
		Destinations: aws.StringSlice(destinations),
		RawMessage:   &rawMessage,
	})

	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}
	return nil
}

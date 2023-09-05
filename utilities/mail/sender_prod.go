package mail

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"gopkg.in/gomail.v2"
)

func SendEmailProd(session *session.Session, fromEmailAddress string, to []string, subject string, content string, attachedFiles []string) error {
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", fromEmailAddress, "")
	msg.SetHeader("To", to...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", content)

	// Attach files if any
	for _, f := range attachedFiles {
		msg.Attach(f)
	}

	var emailRaw bytes.Buffer
	_, err := msg.WriteTo(&emailRaw)
	if err != nil {
		return err
	}

	// Create a new raw message
	rawMessage := ses.RawMessage{Data: []byte(content)}

	// Create SES service client
	svc := ses.New(session)

	// Send raw email
	_, err = svc.SendRawEmail(&ses.SendRawEmailInput{
		Source:       &fromEmailAddress,
		Destinations: aws.StringSlice(to),
		RawMessage:   &rawMessage,
	})

	if err != nil {
		return err
	}

	return nil
}

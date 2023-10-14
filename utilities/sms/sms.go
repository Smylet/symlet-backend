package sms

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	// "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"

	// "github.com/Smylet/symlet-backend/utilities/utils"
)
type SMSSender struct{
	AWSsession *session.Session

}

func NewSMSSender(AWSsession *session.Session) *SMSSender{
	return &SMSSender{
		AWSsession: AWSsession,
	}
}

func (s SMSSender) SendSMS(phoneNumber , message string) error {
	
	smsSvc := sns.New(s.AWSsession, aws.NewConfig().WithLogLevel(aws.LogDebugWithHTTPBody))

	params := &sns.PublishInput{
		Message:     aws.String(message),
		PhoneNumber: aws.String(phoneNumber),
	}
	_, err := smsSvc.Publish(params)

	if awsErr, ok := err.(awserr.Error); ok {
		return fmt.Errorf("unable to send sms %w", awsErr)
	} else if err != nil {
		return fmt.Errorf("unable to send sms %w", err)
	}
	return nil
}

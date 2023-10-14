package common

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/Smylet/symlet-backend/utilities/utils"
)

func CreateAWSSession(config *utils.Config) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.AWSRegion),
		Credentials: credentials.NewStaticCredentials(
			config.AwsAccessKeyID,
			config.AwsSecretAccessKey,
			"",
		),
	},
	)
	if err != nil {
		return nil, fmt.Errorf("error occurred while creating AWS session: %w", err)
	}
	return sess, nil
}

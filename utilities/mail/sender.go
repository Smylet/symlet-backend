package mail

import (
	"bytes"
	"html/template"
	"sync"

	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/utilities/mail/templates"
	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/aws/aws-sdk-go/aws/session"
)

type EmailSender interface {
	SendEmail(
		[]PersonalizedData,
	) []error
}

type PersonalizedData struct {
	Subject       string
	Cc            []string
	Bcc           []string
	Content       string
	User          users.UserSerializer
	AttachedFiles []string
	Others        map[string]interface{}
}
type SESEmailSender struct {
	fromEmailAddress string
	session          *session.Session
	config           utils.Config
}

func NewSESEmailSender(fromEmailAddress string, session *session.Session, config utils.Config) EmailSender {
	return &SESEmailSender{
		fromEmailAddress: fromEmailAddress,
		session:          session,
		config:           config,
	}
}

func (sender *SESEmailSender) SendEmail(
	dataList []PersonalizedData,
) []error {
	// Create a channel to receive errors
	errorsChan := make(chan error, len(dataList))

	// Semaphore to limit the number of concurrent goroutines
	sem := make(chan struct{}, 10)

	// Wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	for _, data := range dataList {
		wg.Add(1)
		go func(data PersonalizedData) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			// Generate personalized email content
			content, err := GeneratePersonalizedEmail(data)
			if err != nil {
				errorsChan <- err
				return
			}

			// Send email
			if sender.config.Environment == "development" || sender.config.Environment == "test" {
				err = SendEmailDev(data.Subject, content, []string{data.User.Email}, sender.fromEmailAddress)
				if err != nil {
					errorsChan <- err
					return
				}

			} else {

				err = SendEmailProd(sender.session, sender.fromEmailAddress, []string{data.User.Email}, data.Subject, content, data.AttachedFiles)
				if err != nil {
					errorsChan <- err
					return
				}

			}
		}(data)
	}

	go func() {
		wg.Wait()
		close(errorsChan)
	}()

	// Collect errors
	var errors []error
	for err := range errorsChan {
		errors = append(errors, err)
	}

	return errors
}

// Function to generate personalized email content
func GeneratePersonalizedEmail(data PersonalizedData) (string, error) {
	// Parse the registration template
	t, err := template.New("registrationEmail").Parse(templates.RegistrationTemplate)
	if err != nil {
		return "", err
	}

	// Create a buffer to store the generated email content
	var emailContent bytes.Buffer

	// Execute the template with the provided personalized data
	err = t.Execute(&emailContent, data)
	if err != nil {
		return "", err
	}

	return emailContent.String(), nil
}

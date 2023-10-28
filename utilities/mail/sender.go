package mail

import (
	"bytes"
	"sync"
	"text/template"

	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/aws/aws-sdk-go/aws/session"
)

type EmailSender interface {
	SendEmail(
		[]Data,
	) []error
}

type Data struct {
	Subject       string
	To            []string
	From          string
	Cc            []string
	Bcc           []string
	Content       string
	Email         string
	AttachedFiles []string
	Url           string
	UserName      string
	EmailTemplate string
	TemplateName  string
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
	dataList []Data,
) []error {
	// Create a channel to receive errors
	errorsChan := make(chan error, len(dataList))

	// Semaphore to limit the number of concurrent goroutines
	sem := make(chan struct{}, 10)

	// Wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	for _, data := range dataList {
		wg.Add(1)
		go func(data Data) {
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
			switch sender.config.Environment {
			case "development", "test":
				// log the email to the terminal
				if err != nil {
					errorsChan <- err
					return
				}

			default:
				err = SendEmailProd(sender.session, sender.fromEmailAddress, data.To, data.Subject, content, data.AttachedFiles)
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
// func GeneratePersonalizedEmail(data Data) (string, error) {
// 	// Parse the registration template
// 	t, err := template.New(data.TemplateName).Parse(data.EmailTemplate)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Create a buffer to store the generated email content
// 	var emailContent bytes.Buffer

// 	// Execute the template with the provided personalized data
// 	err = t.Execute(&emailContent, data)
// 	if err != nil {
// 		return "", err
// 	}

// 	return emailContent.String(), nil
// }

func GeneratePersonalizedEmail(data Data) (string, error) {
	// Parse the HTML template
	tmpl, err := template.New("emailTemplate").Parse(data.EmailTemplate)
	if err != nil {
		return "", err
	}

	// Create a buffer to store the generated HTML
	var emailContent bytes.Buffer

	// Execute the template with email data
	err = tmpl.Execute(&emailContent, data)
	if err != nil {
		return "", err
	}

	to := ""
	// Compose the email
	if len(data.To) > 0 {
		// Now it's safe to access sliceOrArray[0]
		to = data.To[0]
	}

	msg := []byte("Subject: " + data.Subject + "\r\n" +
		"From: " + data.From + "\r\n" +
		"To: " + to + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=utf-8\r\n" +
		"\r\n" +
		emailContent.String())

	return string(msg), nil
}

package mail

import (
	"fmt"
	"net"
	"net/smtp"
	"time"

	smtpmock "github.com/mocktools/go-smtp-mock"
	"github.com/rs/zerolog/log"
)

// You might want to make this a global var or configure it outside of the function
var server = smtpmock.New(smtpmock.ConfigurationAttr{
	LogToStdout:       true,
	LogServerActivity: true,
})

func init() {
	// Starting the mock SMTP server on init, so it's ready for any function calls.
	if err := server.Start(); err != nil {
		fmt.Println("Failed to start mock SMTP server:", err)
	}
}

func SendEmailDev(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string, // Note: The current code doesn't handle attachments yet.
	fromEmailAddress string,
) error {
	hostAddress, portNumber := "127.0.0.1", server.PortNumber
	address := fmt.Sprintf("%s:%d", hostAddress, portNumber)
	timeout := time.Duration(2) * time.Second

	// Connecting to the mock SMTP server and sending an email
	connection, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return fmt.Errorf("failed to dial SMTP server: %w", err)
	}
	defer connection.Close()

	client, err := smtp.NewClient(connection, hostAddress)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	if err := client.Mail(fromEmailAddress); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to start data command: %w", err)
	}

	_, err = wc.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("failed to write content: %w", err)
	}

	// log the email content
	log.Info().Str("subject", subject).Str("content", content).Msg("email content")

	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	if err := client.Quit(); err != nil {
		return fmt.Errorf("failed SMTP Quit command: %w", err)
	}

	return nil
}

// You might want to add a cleanup function if your application needs to shutdown the server gracefully
func Cleanup() {
	if err := server.Stop(); err != nil {
		fmt.Println("Failed to stop mock SMTP server:", err)
	}
}

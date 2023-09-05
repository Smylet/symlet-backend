package mail

import (
	"testing"

	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := utils.LoadConfig("../../env")
	require.NoError(t, err)

	sess, err := common.CreateAWSSession(&config)
	require.NoError(t, err)

	sender := NewSESEmailSender(
		config.EmailSenderAddress,
		sess,
		config,
	)

	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from <a href="http://symlet.com">Symlet</a></p>
	`
	to := []string{"abdulrahmanolamilekan88@gmail.com"}
	attachFiles := []string{}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}

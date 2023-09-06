package mail

import (
	"testing"

	"github.com/Smylet/symlet-backend/api/users"
	"github.com/Smylet/symlet-backend/utilities/common"
	"github.com/Smylet/symlet-backend/utilities/utils"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := utils.LoadConfig("../..")
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

	errs := sender.SendEmail([]PersonalizedData{
		{
			User: users.UserSerializer{
				User: users.User{
					Email: "abdulrahmanolamilekan88@gmail.com",
				},
			},
			Subject:       subject,
			Content:       content,
			Cc:            []string{},
			Bcc:           []string{},
			AttachedFiles: []string{},
		},
	})

	require.Empty(t, errs)
}

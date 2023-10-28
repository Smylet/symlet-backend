package templates

const (
	// Email template for registration
	Registration = "registrationEmail"
	// Email template for reset password
	ResetPassword = "resetPasswordEmail"
	// Email template for change password
	ChangePassword = "changePasswordEmail"
	// Email template for email verification reminder
	EmailVerificationReminder = "emailVerificationReminderEmail"
)

const RegistrationTemplate = `
Hello {{.UserName}},

Thank you for registering with us! 

Please click the link below to verify your email:
<a href="{{.Url}}">click here</a>

Best Regards,
Smylet Team
`

const ResetPasswordTemplate = `
Hello {{.UserName}},
Please click the link below to reset your password:
<a href="{{.Url}}">click here</a>

Best Regards,
Smylet Team
`

const ChangePasswordTemplate = `
Dear {{.UserName}},
Your password has been changed successfully.

Best Regards,
Smylet Team
`

const EmailVerificationReminderTemplate = `
Dear {{.UserName}},

It seems you haven't verified your email address with us yet. Please verify your email to enjoy the full benefits of our service.

If you did not receive the verification email, please check your spam folder or contact our support team.

Thank you for your attention.

Best Regards,
Smylet Team
`

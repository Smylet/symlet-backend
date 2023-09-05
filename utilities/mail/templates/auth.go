package templates

const RegistrationTemplate = `
Hello {{.User.Username}},

Thank you for registering with us! 

Please click the link below to verify your email:
<a href="{{.Others.VerificationLink}}">click here</a>

Best Regards,
Smylet Team
`

const ResetPasswordTemplate = `
Hello {{.User.Username}},
Please click the link below to reset your password:
<a href="{{.Others.ResetPasswordLink}}">click here</a>

Best Regards,
Smylet Team
`

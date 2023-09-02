// validators.go
package users

type CreateUserReq struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CreateUserTxParams struct {
	CreateUserReq
	AfterCreate func(user User) error
}

type CreateUserTxResult struct {
	User User
}

type CreateVerifyEmailParams struct {
	UserID     uint
	Email      string `json:"email"`
	SecretCode string `json:"secret_code"`
}

type ConfirmVerifyEmailParams struct {
	UserID     uint   `form:"user_id" binding:"required"`
	VerEmailID uint   `form:"ver_email_id" binding:"required"`
	SecretCode string `form:"secret_code" binding:"required"`
}

type ValidationStatus struct {
	Valid   bool
	Message string
}
type UpdateVerifyEmailParams struct {
	Email      string
	SecretCode string
}

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
	Username   string `json:"username"`
	Email      string `json:"email"`
	SecretCode string `json:"secret_code"`
}
type UpdateVerifyEmailParams struct {
	Email      string
	SecretCode string
}

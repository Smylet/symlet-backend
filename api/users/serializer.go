package users

type UserSerializer struct {
	User
}

type ProfileSerializer struct {
	Profile
}

func (s *UserSerializer) Response() map[string]interface{} {
	response := map[string]interface{}{
		"id":         s.ID,
		"username":   s.Username,
		"email":      s.Email,
		"created_at": s.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
	}
	return response
}

func (s *ProfileSerializer) Response() map[string]interface{} {
	response := map[string]interface{}{
		"id":      s.ID,
		"bio":     s.Bio,
		"image":   s.Image,
		"user_id": s.UserID,
	}
	return response
}

type VerifyEmailRequest struct {
	Email     string `json:"email"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"` // This can be used to send the expiration time of the token, if needed.
}

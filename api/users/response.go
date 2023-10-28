package users

func (s *UserSerializer) Response(scenario string) map[string]interface{} {

	responses := map[string]map[string]interface{}{
		"login": {
			"session_id":               s.SessionID,
			"access_token":             s.AccessToken,
			"access_token_expires_at":  s.AccessTokenExpiresAt,
			"refresh_token":            s.RefreshToken,
			"refresh_token_expires_at": s.RefreshTokenExpiresAt,
		},
		"create": {
			"user_name": s.UserName,
			"email":     s.Email,
		},

		"get": {
			"user_name": s.UserName,
			"email":     s.Email,
		},
		"renewAccessToken": {
			"access_token":            s.AccessToken,
			"access_token_expires_at": s.AccessTokenExpiresAt,
		},
	}

	response := responses[scenario]

	// Check if s.User is nil before accessing its properties
	if s.User != nil {
		switch scenario {
		case "get":
			response["uid"] = s.User.UID
			response["is_email_confirmed"] = s.User.IsEmailConfirmed
			response["created_at"] = s.User.CreatedAt
			response["updated_at"] = s.User.UpdatedAt
		case "login":
			response["user_uid"] = s.User.UID
		}
	}

	// If the scenario is "get" and profile exists, add it to the response
	if scenario == "get" && s.User != nil && s.User.Profile.ID != 0 {
		response["profile"] = s.User.Profile
	}

	if s.TwoFACode != 0 {
		response["msg"] = "Please check your email for login code"
	}
	return response
}

func (s *UserSerializer) ListResponse(users []User) []map[string]interface{} {
	usersList := make([]map[string]interface{}, 0)

	for _, user := range users {
		userSerializer := UserSerializer{User: &user}
		userSerializer.Email = user.Email
		userSerializer.UserName = user.UserName

		userData := userSerializer.Response("get") // We're assuming "get" is the scenario for retrieving user data.
		usersList = append(usersList, userData)
	}

	return usersList
}

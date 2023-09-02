package users

import (
	"net"
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	// PasswordRegex is a regex for password validation.
	PasswordRegex = regexp.MustCompile(`^[a-zA-Z\d]{8,}$`)

	// UsernameRegex is a regex for username validation.
	UsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{4,}$`)
)

func ValidateRegisterUserReq(req CreateUserReq) ValidationStatus {
	if len(req.Email) < 3 || len(req.Email) > 254 {
		return ValidationStatus{
			Valid:   false,
			Message: "Email must be between 3 and 254 characters",
		}
	}

	if !emailRegex.MatchString(req.Email) {
		return ValidationStatus{
			Valid:   false,
			Message: "Invalid email",
		}
	}

	parts := strings.Split(req.Email, "@")
	mx, err := net.LookupMX(parts[1])

	if err != nil || len(mx) == 0 {
		return ValidationStatus{
			Valid:   false,
			Message: "Invalid email",
		}
	}

	if len(req.Password) < 8 {
		return ValidationStatus{
			Valid:   false,
			Message: "Password must be at least 8 characters",
		}
	}

	if !PasswordRegex.MatchString(req.Password) {
		return ValidationStatus{
			Valid:   false,
			Message: "Password must only contain letters and numbers",
		}
	}

	if len(req.Username) < 4 {
		return ValidationStatus{
			Valid:   false,
			Message: "Username must be at least 4 characters",
		}
	}

	if !UsernameRegex.MatchString(req.Username) {
		return ValidationStatus{
			Valid:   false,
			Message: "Username must only contain letters, numbers, and underscores",
		}
	}

	return ValidationStatus{
		Valid:   true,
		Message: "",
	}
}

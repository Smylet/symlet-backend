package users

import (
	"fmt"

	"github.com/go-playground/validator"
)

func (s *UserSerializer) Validate() error {
	// Instantiate the validator
	validate := validator.New()
	var errorMessage string
	// Validate the struct based on the scenario-specific tags
	if err := validate.Struct(s); err != nil {

		validationErrors := err.(validator.ValidationErrors)

		// Return validation errors to the client
		for _, err := range validationErrors {
			errorMessage += fmt.Sprintf("Field %s: Validation Error (%s) \n", err.Field(), err.Tag())
		}

	}

	// parts := strings.Split(s.Email, "@")
	// mx, err := net.LookupMX(parts[1])
	// if err != nil || len(mx) == 0 {
	// 	errorMessage += fmt.Sprintf("Email validation failed: %s is not a valid email address", s.Email)
	// 	return fmt.Errorf(errorMessage)
	// }

	return nil
}

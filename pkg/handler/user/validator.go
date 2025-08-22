package user

import (
	"errors"
	"regexp"
	"subscription-management/pkg/handler"
)

var (
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrInvalidPhoneNumber = errors.New("invalid phone number format")
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

// ValidateEmail checks if the provided email is in a valid format.
func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

// ValidatePhoneNumber checks if the provided phone number is in a valid format.
func ValidatePhoneNumber(phone string) error {
	if !phoneRegex.MatchString(phone) {
		return ErrInvalidPhoneNumber
	}
	return nil
}

func ValidateUserRequest(user handler.UserRequest) error {
	if user.Name == "" {
		return errors.New("name is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	if user.PhoneNumber == "" {
		return errors.New("phone number is required")
	}
	// Validate email and phone number formats
	if err := ValidateEmail(user.Email); err != nil {
		return err
	}
	if err := ValidatePhoneNumber(user.PhoneNumber); err != nil {
		return err
	}
	return nil
}

package handler

import (
	"errors"
	"regexp"
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

package user

import (
	"errors"
	"subscription-management/pkg/handler"
)

func ValidateUserRequest(user handler.UserRequest) error {
	if user.Name == "" {
		return errors.New("name is required")
	}
	if user.Email_id == "" {
		return errors.New("email id is required")
	}
	if user.Phone_number == "" {
		return errors.New("phone number is required")
	}
	// Validate email and phone number formats
	if err := handler.ValidateEmail(user.Email_id); err != nil {
		return err
	}
	if err := handler.ValidatePhoneNumber(user.Phone_number); err != nil {
		return err
	}
	return nil
}

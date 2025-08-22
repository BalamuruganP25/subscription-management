package customer

import (
	"errors"
	"subscription-management/pkg/handler"
)

func ValidateCustomerRequest(user handler.CreateCustomerRequest) error {
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

func ValidateSubscriptionRequest(sub handler.CreateSubscriptionRequest) error {
	if sub.CustomerID == "" {
		return errors.New("customer id is required")
	}
	if sub.PriceID == "" {
		return errors.New("price id is required")
	}
	return nil
}

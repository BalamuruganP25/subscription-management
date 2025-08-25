package tax

import (
	"errors"
	"strconv"
)

func ValidateTaxRequest(country, amount, state, city, postalCode string) error {
	if country == "" || amount == "" || state == "" || city == "" || postalCode == "" {
		return errors.New("missing required parameters")
	}

	if _, err := strconv.ParseFloat(amount, 64); err != nil {
		return errors.New("invalid amount")
	}

	return nil
}

func getCurrency(country string) string {
	currency := "usd"
	switch country {
	case "IN":
		currency = "inr"
	case "GB":
		currency = "gbp"
	case "EU", "FR", "DE", "IT", "NL", "ES":
		currency = "eur"
	}
	return currency
}

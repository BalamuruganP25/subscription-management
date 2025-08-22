package tax

import (
	"errors"
	"strconv"
)

func ValidateTaxRequest(country, amount, state string) error {
	if country == "" || amount == "" || state == "" {
		return errors.New("missing required parameters")
	}

	if _, err := strconv.ParseFloat(amount, 64); err != nil {
		return errors.New("invalid amount")
	}

	return nil
}


func getTaxRate(country, state string) float64 {
	if country == "US" {
		switch state {
		case "CA":
			return 7.25
		case "NY":
			return 8.875
		default:
			return 5.0
		}
	}
	return 10.0
}
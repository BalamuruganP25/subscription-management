package tax

import (
	"fmt"
	"net/http"
	"strconv"
	"subscription-management/pkg/handler"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/tax/calculation"

	"github.com/go-chi/chi/v5"
)

func GetTax(s *handler.ProcessConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		country := chi.URLParam(r, "country")
		state := chi.URLParam(r, "state")
		city := chi.URLParam(r, "city")
		amountStr := chi.URLParam(r, "amount")
		postalCode := chi.URLParam(r, "postal_code")

		fmt.Println("Country:", country, "State:", state, "City:", city, "Amount:", amountStr, "Postal Code:", postalCode)

		err := ValidateTaxRequest(country, amountStr, state, city, postalCode)
		if err != nil {
			handler.ErrorResponse(w, http.StatusBadRequest,
				handler.ErrResponse{
					Title:   "validation error",
					Details: err.Error(),
				},
			)
			return
		}

		amount, err := strconv.ParseInt(amountStr, 10, 64)
		if err != nil {
			handler.ErrorResponse(w, http.StatusBadRequest,
				handler.ErrResponse{
					Title:   "validation error",
					Details: "invalid amount",
				},
			)
			return
		}

		currency := getCurrency(country)
		stripe.Key = s.StripeKey
		params := &stripe.TaxCalculationParams{
			Currency: stripe.String(currency),
			CustomerDetails: &stripe.TaxCalculationCustomerDetailsParams{
				AddressSource: stripe.String("billing"),
				Address: &stripe.AddressParams{
					Country:    stripe.String(country),
					State:      stripe.String(state),
					City:       stripe.String(city),
					PostalCode: stripe.String(postalCode),
					Line1:      stripe.String("123 Placeholder Lane"),
				},
			},
			LineItems: []*stripe.TaxCalculationLineItemParams{
				{
					Amount:      stripe.Int64(amount),
					TaxBehavior: stripe.String("exclusive"),
					Reference:   stripe.String("prod_001"),
				},
			},
		}

		// Stripe Tax Calculation
		result, err := calculation.New(params)
		if err != nil {
			handler.ErrorResponse(w, http.StatusBadRequest,
				handler.ErrResponse{
					Title:   "failed to calculate tax",
					Details: err.Error(),
				},
			)
			return
		}

		// Send the tax information in the response
		handler.SendResponse(w, result, http.StatusOK)
	}
}

package tax

import (
	"fmt"
	"net/http"
	"strconv"
	"subscription-management/pkg/handler"

	"github.com/go-chi/chi/v5"
)

func GetTax(s *handler.ProcessConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		country := chi.URLParam(r, "country")
		state := chi.URLParam(r, "state")
		amountStr := chi.URLParam(r, "amount")

		fmt.Println("Country:", country, "State:", state, "Amount:", amountStr)

		err := ValidateTaxRequest(country, amountStr, state)
		if err != nil {
			handler.ErrorResponse(w, http.StatusBadRequest,
				handler.ErrResponse{
					Title:   "validation error",
					Details: err.Error(),
				},
			)
			return
		}

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			handler.ErrorResponse(w, http.StatusBadRequest,
				handler.ErrResponse{
					Title:   "validation error",
					Details: "invalid amount",
				},
			)
			return
		}

		taxRate := getTaxRate(country, state)
		taxAmount := amount * taxRate / 100

		resp := handler.TaxResponse{
			Country:   country,
			State:     state,
			TaxRate:   taxRate,
			TaxAmount: taxAmount,
			Amount:    amount,
		}

		// Send the created customer in the response
		handler.SendResponse(w, resp, http.StatusOK)
	}
}

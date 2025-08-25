package tax_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"subscription-management/pkg/handler"
	"subscription-management/pkg/handler/tax"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)



func TestGetTaxHandler(t *testing.T) {
	tests := []struct {
		name     string
		country  string
		state    string
		amount   string
		postalCode string
		city	 string
		wantCode int
	}{
		{
			name:     "Valid request",
			country:  "IN",
			state:    "KA",
			amount:   "100.0",
			city:	 "Bangalore",
			postalCode: "560001",
			wantCode: http.StatusOK,
		},
		{
			name:     "Invalid amount format",
			country:  "IN",
			state:    "KA",
			amount:   "abc",
			city:	 "Bangalore",
			postalCode: "560001",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Missing country",
			country:  "",
			state:    "KA",
			amount:   "100.0",
			city:	 "Bangalore",
			postalCode: "560001",
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request
			req := httptest.NewRequest("GET", "/", nil)

			// Inject chi route context with URL parameters
			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("country", tc.country)
			routeCtx.URLParams.Add("state", tc.state)
			routeCtx.URLParams.Add("amount", tc.amount)

			// Attach to request context
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

			// Create response recorder
			rec := httptest.NewRecorder()

			// Call the handler
			handlerFunc := tax.GetTax(&handler.ProcessConfig{})
			handlerFunc.ServeHTTP(rec, req)

			// Assertions
			assert.Equal(t, tc.wantCode, rec.Code)

			if rec.Code == http.StatusOK {
				var resp handler.TaxResponse
				err := json.NewDecoder(rec.Body).Decode(&resp)
				assert.NoError(t, err)
				assert.Equal(t, tc.country, resp.Country)
				assert.Equal(t, tc.state, resp.State)
			} else {
				var errResp handler.ErrResponse
				err := json.NewDecoder(rec.Body).Decode(&errResp)
				assert.NoError(t, err)
				assert.NotEmpty(t, errResp.Title)
			}
		})
	}
}

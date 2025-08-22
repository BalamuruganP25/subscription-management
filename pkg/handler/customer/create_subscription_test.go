package customer_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"subscription-management/pkg/handler"
	"subscription-management/pkg/handler/customer"
	"subscription-management/pkg/mocks"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/suite"
)

type createSubscriptionTestSuite struct {
	suite.Suite
	router   *chi.Mux
	recorder *httptest.ResponseRecorder
	// CurdRepo *mocks.CrudRepo
}

func (s *createSubscriptionTestSuite) SetupTest() {
	s.router = chi.NewRouter()
	s.recorder = httptest.NewRecorder()
	// s.CurdRepo = new(mocks.CrudRepo)
	config := handler.ProcessConfig{}
	s.router.Post("/v1/api/subscriptions", customer.CreateSubscription(&config))
}

func TestCreateSubscriptionRequest(t *testing.T) {
	testCases := []struct {
		name     string
		reqBody  string
		wantCode int
		mockFunc func(repo *mocks.CrudRepo)
	}{

		{
			name: "customer id missing",
			reqBody: `{"price_id":"mo1",
				"promo_code":"CHOID","customer_id":""}`,
			wantCode: http.StatusBadRequest,
			mockFunc: nil,
		},

		{
			name: "price_id is missing",
			reqBody: `{"price_id":"",
				"promo_code":"CHOID","subscription_id":"783",
				"subscription_status":"ACTIVE","customer_id":""}`,
			wantCode: http.StatusBadRequest,
			mockFunc: nil,
		},

		{
			name:     "invalid json",
			reqBody:  `{"name":}`,
			wantCode: http.StatusBadRequest,
			mockFunc: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := chi.NewRouter()
			recorder := httptest.NewRecorder()
			mockRepo := new(mocks.CrudRepo)
			config := handler.ProcessConfig{CurdRepo: mockRepo}
			router.Post("/v1/api/subscriptions", customer.CreateSubscription(&config))

			if tc.mockFunc != nil {
				tc.mockFunc(mockRepo)
			}

			req := httptest.NewRequest("POST", "/v1/api/subscriptions", strings.NewReader(tc.reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(recorder, req)
			if recorder.Code != tc.wantCode {
				t.Errorf("expected status %d, got %d", tc.wantCode, recorder.Code)
			}
		})
	}

	suite.Run(t, new(createSubscriptionTestSuite))
}

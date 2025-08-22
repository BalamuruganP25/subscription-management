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

type createCustomerTestSuite struct {
	suite.Suite
	router   *chi.Mux
	recorder *httptest.ResponseRecorder
	CurdRepo *mocks.CrudRepo
}

func (s *createCustomerTestSuite) SetupTest() {
	s.router = chi.NewRouter()
	s.recorder = httptest.NewRecorder()
	s.CurdRepo = new(mocks.CrudRepo)
	config := handler.ProcessConfig{
		CurdRepo: s.CurdRepo,
	}
	s.router.Post("/v1/api/customer", customer.CreateCustomer(&config))
}

func TestCreateCustomerRequest(t *testing.T) {
	testCases := []struct {
		name     string
		reqBody  string
		wantCode int
		mockFunc func(repo *mocks.CrudRepo)
	}{
		{
			name:     "missing fields",
			reqBody:  `{}`,
			wantCode: http.StatusBadRequest,
			mockFunc: nil,
		},
		{
			name:     "invalid email",
			reqBody:  `{"name":"Bob","email_id":"invalid-email","phone_number":"1234567890"}`,
			wantCode: http.StatusBadRequest,
			mockFunc: nil,
		},
		{
			name:     "invalid phone",
			reqBody:  `{"name":"Bob","email_id":"bob@example.com","phone_number":"invalid-phone"}`,
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
			router.Post("/v1/api/customer", customer.CreateCustomer(&config))

			if tc.mockFunc != nil {
				tc.mockFunc(mockRepo)
			}

			req := httptest.NewRequest("POST", "/v1/api/customer", strings.NewReader(tc.reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(recorder, req)
			if recorder.Code != tc.wantCode {
				t.Errorf("expected status %d, got %d", tc.wantCode, recorder.Code)
			}
		})
	}

	suite.Run(t, new(createCustomerTestSuite))
}

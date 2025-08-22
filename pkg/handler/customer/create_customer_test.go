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
	"github.com/stretchr/testify/mock"
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
	s.router.Post("/v1/customer", customer.CreateCustomer(&config))
}

func (s *createCustomerTestSuite) executeCreateCustomerTestSuiteRequest(reqBody string) {
	req := httptest.NewRequest("POST", "/v1/customer", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(s.recorder, req)
}

func TestCreateCustomerRequest(t *testing.T) {
	testCases := []struct {
		name     string
		reqBody  string
		wantCode int
		mockFunc func(repo *mocks.CrudRepo)
	}{
		{
			name:     "valid user",
			reqBody:  `{"name":"Alice","email_id":"alice@example.com","phone_number":"1234567890"}`,
			wantCode: http.StatusCreated,
			mockFunc: func(repo *mocks.CrudRepo) {
				repo.On("CreateCustomer",
					mock.Anything, mock.Anything, "Alice", "alice@example.com", "1234567890",
				).Return(nil)
			},
		},
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
			router.Post("/v1/customer", customer.CreateCustomer(&config))

			if tc.mockFunc != nil {
				tc.mockFunc(mockRepo)
			}

			req := httptest.NewRequest("POST", "/v1/customer", strings.NewReader(tc.reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(recorder, req)
			if recorder.Code != tc.wantCode {
				t.Errorf("expected status %d, got %d", tc.wantCode, recorder.Code)
			}
		})
	}

	suite.Run(t, new(createCustomerTestSuite))
}

func (s *createCustomerTestSuite) TestCreateCustomerSuccess() {
	s.CurdRepo.On("CreateCustomer", mock.Anything, mock.Anything, "Alice", "alice@example.com", "1234567890").Return(nil)
	s.executeCreateCustomerTestSuiteRequest(`{"name":"Alice","email_id":"alice@example.com","phone_number":"1234567890"}`)
	s.Require().Equal(http.StatusCreated, s.recorder.Code)
}

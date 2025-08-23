package user_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"subscription-management/pkg/handler"
	"subscription-management/pkg/handler/user"
	"subscription-management/pkg/mocks"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type createUserTestSuite struct {
	suite.Suite
	router   *chi.Mux
	recorder *httptest.ResponseRecorder
	CurdRepo *mocks.CrudRepo
}

func (s *createUserTestSuite) SetupTest() {
	s.router = chi.NewRouter()
	s.recorder = httptest.NewRecorder()
	s.CurdRepo = new(mocks.CrudRepo)
	config := handler.ProcessConfig{
		CurdRepo: s.CurdRepo,
	}
	s.router.Post("/v1/api/users", user.CreateUser(&config))
}

func (s *createUserTestSuite) executeCreateUserTestSuiteRequest(reqBody string) {
	req := httptest.NewRequest("POST", "/v1/api/users", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(s.recorder, req)
}

func TestCreateUserRequest(t *testing.T) {
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
				repo.On("CreateUser",
					mock.Anything, "Alice", "alice@example.com", "1234567890",
				).Return("1", nil)
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
			router.Post("/v1/api/users", user.CreateUser(&config))

			if tc.mockFunc != nil {
				tc.mockFunc(mockRepo)
			}

			req := httptest.NewRequest("POST", "/v1/api/users", strings.NewReader(tc.reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(recorder, req)
			if recorder.Code != tc.wantCode {
				t.Errorf("expected status %d, got %d", tc.wantCode, recorder.Code)
			}
		})
	}

	suite.Run(t, new(createUserTestSuite))
}

func (s *createUserTestSuite) TestCreateUserSuccess() {
	s.CurdRepo.On("CreateUser", mock.Anything, "Alice", "alice@example.com", "1234567890").Return("1", nil)
	s.executeCreateUserTestSuiteRequest(`{"name":"Alice","email_id":"alice@example.com","phone_number":"1234567890"}`)
	s.Require().Equal(http.StatusCreated, s.recorder.Code)
}

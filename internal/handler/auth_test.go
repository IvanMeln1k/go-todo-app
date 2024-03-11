package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/IvanMeln1k/go-todo-app/internal/domain"
	"github.com/IvanMeln1k/go-todo-app/internal/service"
	mock_service "github.com/IvanMeln1k/go-todo-app/internal/service/mocks"
	"github.com/IvanMeln1k/go-todo-app/pkg/validate"
	"github.com/go-playground/validator"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user domain.User)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           domain.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "ok",
			inputBody: `{"name":"string","username":"string","password":"string"}`,
			inputUser: domain.User{
				Name:     "string",
				Username: "string",
				Password: "string",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user domain.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "{\"id\":1}\n",
		},
		{
			name:                "empty fields",
			inputBody:           `{}`,
			inputUser:           domain.User{},
			mockBehavior:        func(s *mock_service.MockAuthorization, user domain.User) {},
			expectedStatusCode:  400,
			expectedRequestBody: "{\"message\":\"invalid body\"}\n",
		},
		{
			name:      "username already in use",
			inputBody: `{"name":"string","username":"string","password":"string"}`,
			inputUser: domain.User{
				Name:     "string",
				Username: "string",
				Password: "string",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user domain.User) {
				s.EXPECT().CreateUser(user).Return(0, service.ErrUsernameAlreadyInUse)
			},
			expectedStatusCode:  409,
			expectedRequestBody: "{\"message\":\"Username already in use\"}\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			e := echo.New()
			e.POST("/signUp", handler.signUp)
			e.Validator = &validate.CustomValidator{Validator: validator.New()}

			req := httptest.NewRequest(http.MethodPost, "/signUp",
				strings.NewReader(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, testCase.expectedStatusCode, rec.Code)
			assert.Equal(t, testCase.expectedRequestBody, rec.Body.String())
		})
	}
}

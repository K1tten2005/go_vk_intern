package http

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/K1tten2005/go_vk_intern/internal/models"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/auth"
	"github.com/K1tten2005/go_vk_intern/internal/pkg/auth/mocks"
	"github.com/satori/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSignIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockAuthUsecase(ctrl)

	tests := []struct {
		name             string
		reqBody          string
		mockBehavior     func()
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:             "Invalid JSON format",
			reqBody:          `{"login": "somerandomlogin"`,
			mockBehavior:     func() {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "incorrect request",
		},
		{
			name:    "Invalid credentials",
			reqBody: `{"login": "test2025", "password": "Test2026!"}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().SignIn(gomock.Any(), gomock.Any()).Return(models.UserResp{}, "", auth.ErrInvalidCredentials)
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "wrong login or password",
		},
		{
			name:    "User not found",
			reqBody: `{"login": "test1999", "password": "Test2029!"}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().SignIn(gomock.Any(), gomock.Any()).Return(models.UserResp{}, "", auth.ErrUserNotFound)
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "user not found",
		},
		{
			name:    "Successful signin",
			reqBody: `{"login": "test2025", "password": "Test2025!"}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().SignIn(gomock.Any(), gomock.Any()).Return(models.UserResp{}, "valid_token", nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "",
		},
		{
			name:    "Unknown error",
			reqBody: `{"login": "test2025", "password": "Test2025!"}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().SignIn(gomock.Any(), gomock.Any()).Return(models.UserResp{}, "", fmt.Errorf("unknown error"))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "unknown error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewBufferString(tt.reqBody))
			req.Header.Set("Content-Type", "application/json")

			tt.mockBehavior()

			rr := httptest.NewRecorder()
			handler := &AuthHandler{uc: mockUsecase}

			handler.SignIn(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			body := rr.Body.String()
			assert.Contains(t, body, tt.expectedResponse)
		})
	}
}

func TestSignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	id := uuid.NewV4()

	mockUsecase := mocks.NewMockAuthUsecase(ctrl)

	tests := []struct {
		name             string
		reqBody          string
		mockBehavior     func()
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:             "Invalid JSON format",
			reqBody:          `{"login": "somerandomlogin"`,
			mockBehavior:     func() {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "incorrect request",
		},
		{
			name:    "Invalid login",
			reqBody: `{"login": "привет", "password": "Test2025!"}`,
			mockBehavior: func() {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: auth.ErrInvalidLogin.Error(),
		},
		{
			name:    "Invalid password",
			reqBody: `{"login": "test2025", "password": "ну привет"}`,
			mockBehavior: func() {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: auth.ErrInvalidPassword.Error(),
		},
		{
			name:    "Successful signup",
			reqBody: `{"login": "test2025", "password": "Test2025!"}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(models.UserResp{Id: id, Login: "test2025"}, "valid_token", nil)
			},
			expectedStatus:   http.StatusCreated,
			expectedResponse: "",
		},
		{
			name:    "Unknown error",
			reqBody: `{"login": "test2025", "password": "Test2025!"}`,
			mockBehavior: func() {
				mockUsecase.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(models.UserResp{}, "", fmt.Errorf("unknown error"))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "unknown error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString(tt.reqBody))
			req.Header.Set("Content-Type", "application/json")

			tt.mockBehavior()

			rr := httptest.NewRecorder()
			handler := &AuthHandler{uc: mockUsecase}

			handler.SignUp(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			body := rr.Body.String()
			assert.Contains(t, body, tt.expectedResponse)
		})
	}
}

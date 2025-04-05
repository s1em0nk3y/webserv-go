package app

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterNewUser(t *testing.T) {
	appWithValidStorage := App{&userStorageMocker{err: nil}}
	appWithInvalidStorage := App{&userStorageMocker{err: errors.New("err")}}
	testcases := []struct {
		name           string
		expectedStatus int
		app            *App
		bodyData       io.Reader
	}{
		{
			name:           "Valid use, valid storage, returns 200",
			expectedStatus: 200,
			bodyData:       strings.NewReader(`{"user":"testuser","password":"password"}`),
			app:            &appWithValidStorage,
		},
		{
			name:           "Valid use, invalid storage, returns 500",
			expectedStatus: 500,
			bodyData:       strings.NewReader(`{"user":"testuser","password":"password"}`),
			app:            &appWithInvalidStorage,
		},
		{
			name:           "Invalid data (non-json) in request",
			expectedStatus: 400,
			bodyData:       strings.NewReader(`{no`),
			app:            &appWithValidStorage,
		},
		{
			name:           "Invalid data (json) in request",
			expectedStatus: 400,
			bodyData:       strings.NewReader(`{}`),
			app:            &appWithValidStorage,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, "/", testcase.bodyData)
			response := httptest.NewRecorder()
			testcase.app.RegisterNewUser(response, request)
			if response.Code != testcase.expectedStatus {
				t.Errorf("got status [%d] wanted [%d]", response.Code, testcase.expectedStatus)
			}

		})
	}
}

type userStorageMocker struct {
	err error
}

func (s *userStorageMocker) CreateUser(user, passwordHash string) error {
	return s.err
}

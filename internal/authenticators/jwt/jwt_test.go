package jwt

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

func TestAuthenticateJWT(t *testing.T) {
	signKey := "test"
	validUser := "testuser"
	authenticator := NewAuthenticator("HS256", signKey, &userStorageMocker{validUser: validUser})
	testcases := []struct {
		notSelfPath    bool
		name           string
		expectedStatus int
		body           string
		signingMethod  jwt.SigningMethod
		claims         claims
		noHeader       bool
	}{
		{
			name:           "Returns 'passed' and 200",
			expectedStatus: http.StatusOK,
			body:           "passed",
			signingMethod:  jwt.SigningMethodHS256,
			claims:         claims{Username: validUser},
		},
		{
			name:           "Corrupt signing method, returns 401",
			expectedStatus: 401,
			signingMethod:  jwt.SigningMethodHS512,
			claims:         claims{Username: validUser},
		},
		{
			name:           "Expired token, returns 401",
			expectedStatus: 401,
			signingMethod:  jwt.SigningMethodHS256,
			claims: claims{Username: validUser, RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now()),
			}},
		},
		{
			name:           "Missing Authorization header, returns 401",
			expectedStatus: 401,
			signingMethod:  jwt.SigningMethodHS256,
			claims:         claims{Username: validUser},
			noHeader:       true,
		},
		{
			name:           "User is not found in storage, returns 401",
			expectedStatus: 401,
			signingMethod:  jwt.SigningMethodHS256,
			claims:         claims{Username: "not-found-user"},
		},
		{
			name:           "User tries to get other user's data, returns 403",
			expectedStatus: 403,
			signingMethod:  jwt.SigningMethodHS256,
			claims:         claims{Username: "not-found-user"},
			notSelfPath:    true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, "/", nil)
			response := httptest.NewRecorder()
			signedToken, err := jwt.NewWithClaims(testcase.signingMethod, testcase.claims).SignedString([]byte(signKey))
			if err != nil {
				t.Fatal(err)
			}
			if !testcase.noHeader {
				request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", signedToken))
			}
			if !testcase.notSelfPath {
				request = mux.SetURLVars(request, map[string]string{"id": testcase.claims.Username})
			} else {
				request = mux.SetURLVars(request, map[string]string{"id": fmt.Sprintf("non%s", testcase.claims.Username)})
			}
			handler := authenticator.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("passed"))
			}))
			handler.ServeHTTP(response, request)
			got := response.Body.String()
			if testcase.body != "" && got != testcase.body {
				t.Errorf("got [%s] wanted [%s]", got, testcase.body)
			}
			if response.Code != testcase.expectedStatus {
				t.Errorf("got status [%d] wanted [%d]", response.Code, testcase.expectedStatus)
			}
		})
	}
}

func TestRegisterNewUser(t *testing.T) {
	signKey := "test"
	appWithValidStorage := NewAuthenticator("HS256", signKey, &userStorageMocker{err: nil})
	appWithInvalidStorage := NewAuthenticator("HS256", signKey, &userStorageMocker{err: errors.New("err")})
	testcases := []struct {
		name           string
		expectedStatus int
		app            *AuthJWT
		bodyData       io.Reader
	}{
		{
			name:           "Valid use, valid storage, returns 200",
			expectedStatus: 200,
			bodyData:       strings.NewReader(`{"user":"testuser","password":"password"}`),
			app:            appWithValidStorage,
		},
		{
			name:           "Valid use, invalid storage, returns 500",
			expectedStatus: 500,
			bodyData:       strings.NewReader(`{"user":"testuser","password":"password"}`),
			app:            appWithInvalidStorage,
		},
		{
			name:           "Invalid data (non-json) in request",
			expectedStatus: 400,
			bodyData:       strings.NewReader(`{no`),
			app:            appWithValidStorage,
		},
		{
			name:           "Invalid data (json) in request",
			expectedStatus: 400,
			bodyData:       strings.NewReader(`{}`),
			app:            appWithValidStorage,
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

func TestLoginUser(t *testing.T) {
	signKey := "test"
	testcases := []struct {
		name           string
		app            *AuthJWT
		bodyData       io.Reader
		expectedStatus int
	}{
		{
			name:           "Valid, returns 200 and token",
			app:            NewAuthenticator("HS256", signKey, &userStorageMocker{ok: true, err: nil}),
			bodyData:       strings.NewReader(`{"user":"testuser","password":"password"}`),
			expectedStatus: 200,
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, "/", testcase.bodyData)
			response := httptest.NewRecorder()
			testcase.app.LoginUser(response, request)
			if response.Code != testcase.expectedStatus {
				t.Errorf("got status [%d] wanted [%d]", response.Code, testcase.expectedStatus)
			}
		})
	}
}

type userStorageMocker struct {
	ok        bool
	err       error
	validUser string
}

func (s *userStorageMocker) CreateUser(user, passwordHash string) error {
	return s.err
}

func (s *userStorageMocker) ValidateUsername(username string) error {
	if s.validUser == username {
		return nil
	}
	return errors.New("err")
}
func (s *userStorageMocker) CheckCredents(user, passwordHash string) (bool, error) {
	return s.ok, s.err
}

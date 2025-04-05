package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

func TestAuthenticateJWT(t *testing.T) {
	signKey := "test"
	validUser := "testuser"
	authenticator := NewAuthenticator("HS256", signKey, &userStorageMocker{validUser})
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
			}
			handler := authenticator.AuthenticateJWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

type userStorageMocker struct {
	validUser string
}

func (s *userStorageMocker) ValidateUsername(username string) error {
	if s.validUser == username {
		return nil
	}
	return errors.New("err")
}

package jwt

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type UserStorage interface {
	ValidateUsername()
}

type AuthJWT struct {
	signKey    []byte
	signMethod string
}

func NewAuthenticator(signMethod string, signKey string) *AuthJWT {

	return &AuthJWT{signKey: []byte(signKey), signMethod: signMethod}
}

func (a *AuthJWT) AuthenticateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if len(tokenStr) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing 'Authorization' Header"))
			return
		}
		token, err := jwt.ParseWithClaims(
			tokenStr,
			&claims{},
			func(t *jwt.Token) (interface{}, error) {
				if t.Method.Alg() != a.signMethod {
					return nil, fmt.Errorf("unexpected algorithm %s, wanted %s", t.Method.Alg(), a.signMethod)
				}
				return a.signKey, nil
			},
		)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("cant verify token"))
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("token is not valid"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

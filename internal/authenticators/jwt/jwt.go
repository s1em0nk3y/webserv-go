package jwt

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type UserStorage interface {
	ValidateUsername(username string) error
}

type AuthJWT struct {
	signKey    []byte
	signMethod string
	storage    UserStorage
}

func NewAuthenticator(signMethod string, signKey string, storage UserStorage) *AuthJWT {
	return &AuthJWT{signKey: []byte(signKey), signMethod: signMethod, storage: storage}
}

func (a *AuthJWT) AuthenticateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if len(tokenStr) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing 'Authorization' Header"))
			return
		}
		claims := &claims{}
		if _, err := jwt.ParseWithClaims(
			tokenStr,
			claims,
			func(t *jwt.Token) (interface{}, error) {
				if t.Method.Alg() != a.signMethod {
					return nil, fmt.Errorf("unexpected algorithm %s, wanted %s", t.Method.Alg(), a.signMethod)
				}
				return a.signKey, nil
			},
		); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("cant verify token"))
			return
		}
		if userIDFromRoute := mux.Vars(r)["id"]; claims.Username != userIDFromRoute {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("you cannot view other user's paths"))
		}
		if a.storage.ValidateUsername(claims.Username) != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("user [%s] not found", claims.Username)))
			return
		}
		next.ServeHTTP(w, r)
	})
}

package jwt

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type AuthJWT struct {
	signKey    []byte
	signMethod string
	storage    UserStorage
}

func NewAuthenticator(signMethod string, signKey string, storage UserStorage) *AuthJWT {
	return &AuthJWT{signKey: []byte(signKey), signMethod: signMethod, storage: storage}
}

func (a *AuthJWT) RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	registerUser := &registerUser{}
	if err := json.NewDecoder(r.Body).Decode(registerUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to decode user and password"))
		return
	}

	if registerUser.Password == "" || registerUser.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("username or password not provided"))
	}
	hash := sha256.Sum256([]byte(registerUser.Password))
	if err := a.storage.CreateUser(registerUser.Username, hex.EncodeToString(hash[:])); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to create user due to"))
		return
	}
}

func (a *AuthJWT) LoginUser(w http.ResponseWriter, r *http.Request) {
	loginUser := &loginUser{}
	if err := json.NewDecoder(r.Body).Decode(loginUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to decode user and password"))
		return
	}
	if loginUser.Password == "" || loginUser.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("username or password not provided"))
		return
	}
	hash := sha256.Sum256([]byte(loginUser.Password))
	ok, err := a.storage.CheckCredents(loginUser.Username, hex.EncodeToString(hash[:]))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("some internal error occured"))
		return
	}

	if !ok {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("cant logon"))
		return
	}

	signedString, err := jwt.NewWithClaims(
		jwt.GetSigningMethod(a.signMethod),
		claims{
			Username: loginUser.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now()),
			},
		},
	).SignedString(a.signKey)
	if err != nil {
		w.Write([]byte("some internal error occured"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(signedString))
}

func (a *AuthJWT) Authenticate(next http.Handler) http.Handler {
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
		if userIDFromRoute := mux.Vars(r)["id"]; userIDFromRoute != "" && claims.Username != userIDFromRoute {
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

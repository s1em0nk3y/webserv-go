package app

import (
	"encoding/json"
	"net/http"
)

// CREATE: /users/register
func (a *App) RegisterNewUser(w http.ResponseWriter, r *http.Request) {
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

	if err := a.storage.CreateUser(registerUser.Username, registerUser.Password); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to create user"))
		return
	}
}

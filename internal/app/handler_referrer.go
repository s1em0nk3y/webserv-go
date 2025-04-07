package app

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type ReferralStorage interface {
	AddNewReferral(user string, referral string) error
}

// POST /users/:id/referrer; should contain other user name
func (a *App) AddReferral(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["id"]
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(bytes) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("referral name not provided"))
		return
	}
	if err = a.AddNewReferral(username, string(bytes)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("cannot add referral"))
		return
	}
}

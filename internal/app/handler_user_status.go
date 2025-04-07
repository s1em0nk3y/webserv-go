package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type UserStatusGetter interface {
	GetUserStatus(username string) (*UserStatus, error)
}

type UserStatus struct {
	Username           string  `json:"username"`
	CompletedTaskCount uint    `json:"completed_tasks"`
	Score              float64 `json:"score"`
}

func (a *App) GetUserStatus(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["id"]
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userStatus, err := a.UserStatusGetter.GetUserStatus(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(userStatus)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

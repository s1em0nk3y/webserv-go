package app

import (
	"encoding/json"
	"net/http"
)

type LeaderBoardGetter interface {
	GetLeaderBoard() (LeaderBoard, error)
}

type LeaderBoard []struct {
	User  string  `json:"user"`
	Score float64 `json:"score"`
}

func (a *App) GetLeaderBoard(w http.ResponseWriter, r *http.Request) {
	leaderBoard, err := a.LeaderBoardGetter.GetLeaderBoard()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(leaderBoard)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

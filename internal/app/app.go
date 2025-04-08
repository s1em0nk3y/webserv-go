package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/s1em0nk3y/webserv-go/internal/authenticators/jwt"
)

type App struct {
	ReferralStorage
	TaskCompleter
	LeaderBoardGetter
	UserStatusGetter
	Authenticator *jwt.AuthJWT
}

func (a *App) Run(port uint) {
	router := mux.NewRouter()
	router.HandleFunc("/register", a.Authenticator.RegisterNewUser)
	router.HandleFunc("/login", a.Authenticator.LoginUser)
	authRouter := router.PathPrefix("/users").Subrouter()
	authRouter.Use(a.Authenticator.Authenticate)
	authRouter.HandleFunc("/{id:[A-Za-z0-9]+}/referrer/", a.addReferral)
	authRouter.HandleFunc("/{id:[A-Za-z0-9]+}/task/complete", a.completeTask)
	authRouter.HandleFunc("/leaderboard", a.getLeaderBoard)
	authRouter.HandleFunc("/{id:[A-Za-z0-9]+}/status", a.getUserStatus)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		log.Fatal(err)
	}
}

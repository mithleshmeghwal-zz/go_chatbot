package main

import (
	"log"
	"net/http"

	"chatbot/router"
	"github.com/gorilla/mux"
)

type App struct {
	Routes router.Router
	R *mux.Router
}

func (a *App)CreateRoutes() {

	conn , err := ConnectDatabase("redis", "6379")
	if err != nil {
		log.Fatal(err)
	}
	routes := router.Router {
		Conn: conn,
	}
	a.Routes = routes
	a.R = mux.NewRouter()

	a.R.HandleFunc("/webhook", a.Routes.VerificationEndpoint).Methods("GET")
	a.R.HandleFunc("/webhook", a.Routes.MessagesEndpoint).Methods("POST")
	
}

func (a *App) Run() {
	defer a.Routes.Conn.Close()
	if err := http.ListenAndServe(":8080", a.R); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"net/http"

	"chatbot/router"
	"github.com/gorilla/mux"
)

type App struct {
	Router router.Router
}

// func (a *App)VerificationEndpoint(w http.ResponseWriter, r *http.Request) {
// 	log.Println("Verifying Endpoint")
// 	challenge := r.URL.Query().Get("hub.challenge")
// 	mode := r.URL.Query().Get("hub.mode")
// 	token := r.URL.Query().Get("hub.verify_token")
// 	log.Println("TEST", challenge, mode, token)
// 	log.Println("ENV TOKEN", os.Getenv("VERIFY_TOKEN"))
// 	if mode != "" && token == os.Getenv("VERIFY_TOKEN") {
// 		w.WriteHeader(200)
// 		w.Write([]byte(challenge))
// 	} else {
// 		w.WriteHeader(404)
// 		w.Write([]byte("Error, wrong validation token"))
// 	}
// }


// func (a *App)MessagesEndpoint(w http.ResponseWriter, r *http.Request) {
// 	log.Println("Message Event")
// 	var callback Callback
// 	json.NewDecoder(r.Body).Decode(&callback)
// 	log.Println("MESSAGE RECIEVED ", callback)
// 	if callback.Object == "page" {
// 		for _, entry := range callback.Entry {
// 			for _, event := range entry.Messaging {
// 				a.ProcessMessage(event)
// 			}
// 		}
// 		w.WriteHeader(200)
// 		w.Write([]byte("Got your message"))
// 	} else {
// 		w.WriteHeader(404)
// 		w.Write([]byte("Message not supported"))
// 	}
// }

func (a *App)CreateRoutes() {

	conn , err := ConnectDatabase("redis", "6379")
	if err != nil {
		log.Fatal(err)
	}
	routes := router.Router {
		Conn: conn,
	}
	a.Router = routes
	defer a.Router.Conn.Close()
	r := mux.NewRouter()
	r.HandleFunc("/webhook", a.Router.VerificationEndpoint).Methods("GET")
	r.HandleFunc("/webhook", a.Router.MessagesEndpoint).Methods("POST")
	if err := http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		log.Fatal(err)
	}
}

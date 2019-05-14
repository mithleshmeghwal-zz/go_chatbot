package router

import (
	"os"
	"log"
	"net/http"
)

func (a *Router) VerificationEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Verifying Endpoint")
	challenge := r.URL.Query().Get("hub.challenge")
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	log.Println("TEST", challenge, mode, token)
	log.Println("ENV TOKEN", os.Getenv("VERIFY_TOKEN"))
	if mode != "" && token == os.Getenv("VERIFY_TOKEN") {
		w.WriteHeader(200)
		w.Write([]byte(challenge))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Error, wrong validation token"))
	}
}
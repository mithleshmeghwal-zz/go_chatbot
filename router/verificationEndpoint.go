package router

import (
	"os"
	"log"
	"net/http"
)

func (a *Router) VerificationEndpoint(w http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("hub.challenge")
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	if mode != "" && token == os.Getenv("VERIFY_TOKEN") {
		log.Println("Token Verified")
		w.WriteHeader(200)
		w.Write([]byte(challenge))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Error, wrong validation token"))
	}
}
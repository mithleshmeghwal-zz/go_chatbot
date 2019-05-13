package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Docker")
	})
	http.HandleFunc("/ding", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "DongcBong from Docker")
	})

	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

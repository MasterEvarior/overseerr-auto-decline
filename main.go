package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", webhookHandler)
	http.ListenAndServe(":8080", nil)
}

type payload struct {
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Printf("%s is not a valid HTTP method for this webhook", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

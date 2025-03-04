package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	apiKey := getEnvVar("API_KEY")
	url := getEnvVar("URL")

	client := NewClient(url, apiKey)

	body, err := client.GetRequest("75")
	if err != nil {
		log.Fatal(err)
	}
	log.Print(body)

	//http.HandleFunc("/", webhookHandler)
	//http.ListenAndServe(":8080", nil)
}

type payload struct {
}

func getMedia() []string {
	return strings.Split(getEnvVar("MEDIA"), ",")
}

func getEnvVarWithDefault(name string) {

}

func getEnvVar(name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalf("Environment variable '%s' was not defined", name)
	}
	return value
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Printf("%s is not a valid HTTP method for this webhook", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

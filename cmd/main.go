package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/MasterEvarior/overseerr-auto-decline/cmd/client"
	"github.com/MasterEvarior/overseerr-auto-decline/cmd/handler"
)

func main() {
	apiKey := getEnvVar("API_KEY")
	url := getEnvVar("URL")
	_, deleteRequest := os.LookupEnv("DELETE_REQUESTS")
	mediaIDs := getMedia("MEDIA")

	log.Printf("The media with the following IDs will be processed: %v", mediaIDs)
	if deleteRequest {
		log.Println("Requests will be deleted after they have been declined, if you wish otherwise unset the 'DELETE_REQUESTS' environment variable")
	}

	h := handler.Handler{
		OverseerrClient: client.NewClient(url, apiKey),
		DeleteRequests:  deleteRequest,
		BannedMediaIDs:  mediaIDs,
	}

	http.HandleFunc("/", h.WebhookHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("Could not start the server because of the following issue: %v", err)
	}
}

func getMedia(name string) []string {
	return strings.Split(getEnvVar(name), ",")
}

func getEnvVar(name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalf("Environment variable '%s' was not defined", name)
	}
	return value
}

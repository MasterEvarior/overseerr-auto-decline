package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
)

var apiKey string
var url string
var deleteRequest bool
var mediaIds []string

func main() {
	apiKey = getEnvVar("API_KEY")
	url = getEnvVar("URL")
	_, deleteRequest = os.LookupEnv("DELETE_REQUESTS")
	mediaIds = getMedia("MEDIA")

	log.Printf("The media with the following IDs will be processed: %v", mediaIds)
	if deleteRequest {
		log.Println("Requests will be deleted after they have been declined, if you wish otherwise unset the 'DELETE_REQUESTS' environment variable")
	}

	http.HandleFunc("/", webhookHandler)
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

type WebhookPayload struct {
	RequestId string `json:"request_id"`
	TmDbId    string `json:"tmdbid"`
	TvDbId    string `json:"tvdbid"`
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Printf("%s is not a valid HTTP method for this webhook", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Could not read body of this request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var payload WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("Could not unmarshal body of this request, body was '%s', error is: %v", string(body[:]), err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Successfully received payload on webhook with the following data: %+v", payload)

	if !slices.Contains(mediaIds, payload.TmDbId) && !slices.Contains(mediaIds, payload.TvDbId) {
		log.Printf("'%s' or '%s' not found inside the configured media IDs, doing nothing", payload.TmDbId, payload.TvDbId)
		return
	}

	overseerrClient := NewClient(url, apiKey)
	err = overseerrClient.DeclineRequest(payload.RequestId)
	if err != nil {
		log.Printf("Could not decline request with the id '%s' because of the following error: %v", payload.RequestId, err)
	}

	if deleteRequest {
		err = overseerrClient.DeleteRequest(payload.RequestId)
		if err != nil {
			log.Printf("Could not delete request with the id '%s' because of the following error: %v", payload.RequestId, err)
		}
	}

	log.Print("Finished with this request")
}

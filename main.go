package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
)

var apiKey string
var url string
var deleteRequest bool
var mediaIds []uint

func main() {
	apiKey = getEnvVar("API_KEY")
	url = getEnvVar("URL")
	_, deleteRequest = os.LookupEnv("DELETE_REQUESTS")
	mediaIds = getMedia("MEDIA")

	log.Println("The media with the following IDs will be processed: ", mediaIds)
	if deleteRequest {
		log.Println("Request will be deleted after they have been declined, if you wish otherwise unset the 'DELETE_REQUESTS' environment variable")
	}

	http.HandleFunc("/", webhookHandler)
	http.ListenAndServe(":8080", nil)
}

func getMedia(name string) []uint {
	strings := strings.Split(getEnvVar(name), ",")
	ints := make([]uint, len(strings))

	for i, s := range strings {
		num, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			log.Fatalf("Could not convert media id '%s' to a valid number", s)
		}
		ints[i] = uint(num)
	}
	return ints
}

func getEnvVar(name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalf("Environment variable '%s' was not defined", name)
	}
	return value
}

type WebhookPayload struct {
	tmDbId uint `json:"tmdbid"`
	tvDbId uint `json:"tvdbid"`
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Printf("%s is not a valid HTTP method for this webhook", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var payload WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("Successfully received payload on webhook with the following data:", payload)

	if !slices.Contains(mediaIds, payload.tmDbId) && !slices.Contains(mediaIds, payload.tvDbId) {
		log.Println("'%d' or '%d' not found inside the configured media IDs, doing nothing", payload.tmDbId, payload.tvDbId)
		return
	}

	mediaIdToDecline := payload.tmDbId
	if payload.tmDbId == 0 {
		mediaIdToDecline = payload.tvDbId
	}

	overseerrClient := NewClient(url, apiKey)
	overseerrClient.DeclineRequest(mediaIdToDecline)
	if deleteRequest {
		overseerrClient.DeleteRequest(mediaIdToDecline)
	}
}

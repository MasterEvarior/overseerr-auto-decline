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

	log.Printf("The media with the following IDs will be processed: %v", mediaIds)
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
	RequestId uint `json:"request_id"`
	TmDbId    uint `json:"tmdbid"`
	TvDbId    uint `json:"tvdbid"`
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
		log.Printf("Could not unmarshal body of this request, body was: %s", string(body[:]))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Successfully received payload on webhook with the following data: %+v", payload)

	if !slices.Contains(mediaIds, payload.TmDbId) && !slices.Contains(mediaIds, payload.TvDbId) {
		log.Printf("'%d' or '%d' not found inside the configured media IDs, doing nothing", payload.TmDbId, payload.TvDbId)
		return
	}

	overseerrClient := NewClient(url, apiKey)
	err = overseerrClient.DeclineRequest(payload.RequestId)
	if err != nil {
		log.Printf("Could not decline request with the id '%d' because of the following error: %v", payload.RequestId, err)
	}

	if deleteRequest {
		err = overseerrClient.DeleteRequest(payload.RequestId)
		if err != nil {
			log.Printf("Could not delete request with the id '%d' because of the following error: %v", payload.RequestId, err)
		}
	}

	log.Print("Finished with this request")
}

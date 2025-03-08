package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
)

type handler struct {
	overseer       *OverseerClient
	deleteRequests bool
	bannedMediaIDs []string
}

func main() {
	apiKey := getEnvVar("API_KEY")
	url := getEnvVar("URL")
	_, deleteRequest := os.LookupEnv("DELETE_REQUESTS")
	mediaIDs := getMedia("MEDIA")

	log.Printf("The media with the following IDs will be processed: %v", mediaIDs)
	if deleteRequest {
		log.Println("Requests will be deleted after they have been declined, if you wish otherwise unset the 'DELETE_REQUESTS' environment variable")
	}

	h := handler{
		overseer:       NewClient(url, apiKey),
		deleteRequests: deleteRequest,
		bannedMediaIDs: mediaIDs,
	}

	http.HandleFunc("/", h.webhookHandler)
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
	RequestID string `json:"request_id"`
	TmDbId    string `json:"tmdbid"`
	TvDbId    string `json:"tvdbid"`
}

func (h *handler) webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Printf("%s is not a valid HTTP method for this webhook", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var payload WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Could not unmarshal body of this request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Successfully received payload on webhook with the following data: %+v", payload)
	if !slices.Contains(h.bannedMediaIDs, payload.TmDbId) && !slices.Contains(h.bannedMediaIDs, payload.TvDbId) {
		log.Printf("%q or %q not found inside the configured media IDs, doing nothing", payload.TmDbId, payload.TvDbId)
		return
	}

	err := h.overseer.DeclineRequest(payload.RequestID)
	if err != nil {
		log.Printf("Could not decline request with the id '%s' because of the following error: %v", payload.RequestID, err)
	}

	if h.deleteRequests {
		err = h.overseer.DeleteRequest(payload.RequestID)
		if err != nil {
			log.Printf("Could not delete request with the id '%s' because of the following error: %v", payload.RequestID, err)
		}
	}

	log.Print("Finished with this request")
}

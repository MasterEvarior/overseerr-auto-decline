package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"

	"github.com/MasterEvarior/overseerr-auto-decline/cmd/client"
)

type Handler struct {
	OverseerrClient client.OverseerClient
	DeleteRequests  bool
	BannedMediaIDs  []string
}

type WebhookPayload struct {
	RequestID string `json:"request_id"`
	TmDbId    string `json:"tmdbid"`
	TvDbId    string `json:"tvdbid"`
}

func (h *Handler) WebhookHandler(w http.ResponseWriter, r *http.Request) {
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

	if payload.RequestID == "" {
		http.Error(w, "Missing required field: request_id", http.StatusBadRequest)
		return
	}

	log.Printf("Successfully received payload on webhook with the following data: %+v", payload)
	if !slices.Contains(h.BannedMediaIDs, payload.TmDbId) && !slices.Contains(h.BannedMediaIDs, payload.TvDbId) {
		log.Printf("%q or %q not found inside the configured media IDs, doing nothing", payload.TmDbId, payload.TvDbId)
		return
	}

	err := h.OverseerrClient.DeclineRequest(payload.RequestID)
	if err != nil {
		log.Printf("Could not decline request with the id '%s' because of the following error: %v", payload.RequestID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Print("Successfully declined the request")

	if h.DeleteRequests {
		err = h.OverseerrClient.DeleteRequest(payload.RequestID)
		if err != nil {
			log.Printf("Could not delete request with the id '%s' because of the following error: %v", payload.RequestID, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Print("Successfully deleted the request")
	}

	log.Print("Finished with this request")
	w.WriteHeader(http.StatusNoContent)
}

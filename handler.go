package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
)

type Handler struct {
	overseer       *OverseerClient
	deleteRequests bool
	bannedMediaIDs []string
}

type WebhookPayload struct {
	RequestID string `json:"request_id"`
	TmDbId    string `json:"tmdbid"`
	TvDbId    string `json:"tvdbid"`
}

func (h *Handler) webhookHandler(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Internal Server Error", 500)
		return
	}
	log.Print("Successfully declined the request")

	if h.deleteRequests {
		err = h.overseer.DeleteRequest(payload.RequestID)
		if err != nil {
			log.Printf("Could not delete request with the id '%s' because of the following error: %v", payload.RequestID, err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		log.Print("Successfully deleted the request")
	}

	log.Print("Finished with this request")
	w.WriteHeader(http.StatusNoContent)
}

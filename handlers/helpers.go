package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/vaForge/TaskManagerAPI/validation"
)

// WriteJSON sends a JSON response with given status code.

func writeJSON(w http.ResponseWriter, status int, payload any) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if payload == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to encode JSON response: %v", err)
	}
}

// WriteError Sends a consistent JSON error message
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}

func requireJSONContentType(w http.ResponseWriter, r *http.Request) bool {

	if !validation.IsJSONContentType(r) {
		writeError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return false
	}
	return true
}

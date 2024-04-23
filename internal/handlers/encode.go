package handlers

import (
	"encoding/json"
	"net/http"
)

func writeJson[T any](w http.ResponseWriter, payload T) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

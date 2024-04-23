package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

func writeJson[T any](w http.ResponseWriter, payload T) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func parseBody[T any](body io.Reader, parsedResult T) error {

	if err := json.NewDecoder(body).Decode(parsedResult); err != nil {
		return err
	}

	return nil
}

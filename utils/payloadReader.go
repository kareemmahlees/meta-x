package utils

import (
	"encoding/json"
	"io"

	"github.com/charmbracelet/log"
)

func ReadBody(body io.ReadCloser)map[string]any{
	var parsedPayload map[string]any
	decoder := json.NewDecoder(body)
    err := decoder.Decode(&parsedPayload)
    if err != nil {
        log.Error(err)
    }
	return parsedPayload
}
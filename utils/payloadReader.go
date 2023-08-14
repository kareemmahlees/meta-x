package utils

import (
	"encoding/json"
	"io"

	"github.com/charmbracelet/log"
)

func ReadPayload(payload io.ReadCloser)map[string]any{
	var parsedPayload map[string]any
	decoder := json.NewDecoder(payload)
    err := decoder.Decode(&parsedPayload)
    if err != nil {
        log.Error(err)
    }
	return parsedPayload
}
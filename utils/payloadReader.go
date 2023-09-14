package utils

import (
	"encoding/json"
	"io"

	"github.com/charmbracelet/log"
)

type body interface {
	[]map[string]any | map[string]any
}

func ReadBody[K body](body io.ReadCloser) K {
	var parsedPayload K
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&parsedPayload)
	if err != nil {
		log.Error(err)
	}
	return parsedPayload
}

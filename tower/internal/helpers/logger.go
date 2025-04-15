package helpers

import (
	"encoding/json"
	"fmt"
	"log"
)

// Entry defines the structure of a JSON log entry
type Entry struct {
	URI      string      `json:"uri"`
	Response interface{} `json:"response"`
}

// PrintPretty logs the Entry as pretty-printed JSON
func (e Entry) PrintPretty() {
	// Try to decode raw JSON if Response is []byte
	if raw, ok := e.Response.([]byte); ok {
		var decoded interface{}
		if err := json.Unmarshal(raw, &decoded); err == nil {
			e.Response = decoded
		}
	}

	jsonBytes, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		log.Printf(`{"level":"error","message":"Failed to marshal log entry","uri":"%s","error":"%v"}`, e.URI, err)
		return
	}

	fmt.Println(string(jsonBytes))
}

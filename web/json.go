package web

import (
	"encoding/json"
	"net/http"
)

type envelope map[string]interface{}

// Based on:
// Let's Go Further - Chapter 3.4 (Alex Edwards)
func writeJSON(w http.ResponseWriter, status int, data envelope) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// append a newline for nicer terminal output
	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

package main

import (
	"encoding/json"
	"net/http"
)

func responseWithError(w http.ResponseWriter, code int, message string) {
	responseWithJSON(w, code, map[string]string{"error": message})
}

func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

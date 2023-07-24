package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ChirpRequestBody struct {
	Body string `json:"body"`
}

type ChirpValidationResponse struct {
	Valid bool   `json:"valid"`
	Error string `json:"error"`
}

func handlerValidation(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	requestBody := ChirpRequestBody{}
	err := decoder.Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Chirp is invalid
	if len(requestBody.Body) > 140 {
		responseBody := ChirpValidationResponse{
			Valid: false,
			Error: "Body must be less than 140 letters",
		}
		dat, err := json.Marshal(responseBody)
		if err != nil {
			log.Printf("Error marshalling response: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(dat)

		return
	}

	// Chirp is valid
	responseBody := ChirpValidationResponse{
		Valid: true,
		Error: "",
	}
	dat, err := json.Marshal(responseBody)
	if err != nil {
		log.Printf("Error marshalling response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ChirpRequestBody struct {
	Body string `json:"body"`
}

type ChirpValidationResponse struct {
	CleanedBody string `json:"cleanedBody"`
	Error       string `json:"error"`
}

func handlerValidation(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	requestBody := ChirpRequestBody{}
	err := decoder.Decode(&requestBody)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Chirp is invalid
	if len(requestBody.Body) > 140 {
		responseWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	// Chirp is valid
	cleanedText := replaceProfaneWords(requestBody.Body)

	responseBody := ChirpValidationResponse{
		CleanedBody: cleanedText,
		Error:       "",
	}
	responseWithJSON(w, http.StatusOK, responseBody)
}

func replaceProfaneWords(text string) string {
	// Define the list of profane words to replace
	profaneWords := map[string]string{
		"kerfuffle": "****",
		"sharbert":  "****",
		"fornax":    "****",
	}

	words := strings.Split(text, " ")

	for index, word := range words {
		lowercaseWord := strings.ToLower(word)
		if replacement, ok := profaneWords[lowercaseWord]; ok {
			words[index] = replacement
		}
	}

	cleanedText := strings.Join(words, " ")
	return cleanedText
}

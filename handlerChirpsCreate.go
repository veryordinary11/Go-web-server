package main

import (
	"encoding/json"
	"net/http"

	"github.com/veryordinary11/Go-web-server/database"
)

func handlerChirpsCreate(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			responseWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		// Decode the request body into a Chirp
		decoder := json.NewDecoder(r.Body)
		requestBody := database.Chirp{}
		err := decoder.Decode(&requestBody)
		if err != nil {
			responseWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Validate the chirp
		if len(requestBody.Body) > 140 {
			responseWithError(w, http.StatusBadRequest, "Chirp is too long")
			return
		}

		// Chirp is valid
		cleanedText := replaceProfaneWords(requestBody.Body)

		// Add the chirp to the database
		chirp, err := db.CreateChirp(cleanedText)
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, "Failed to create chirp")
			return
		}

		// Respond with the chirp
		responseWithJSON(w, http.StatusOK, chirp)
	}
}

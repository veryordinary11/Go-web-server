package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/veryordinary11/Go-web-server/database"
)

func handlerChirpsGet(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the ID from the URL
		chirpIDStr := chi.URLParam(r, "id")
		chirpID, err := strconv.Atoi(chirpIDStr)
		if err != nil {
			responseWithError(w, http.StatusBadRequest, "Invalid chirp ID")
			return
		}

		// Get the chirp by ID from the database
		chirp, err := db.GetChirpByID(chirpID)
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, "Failed to get chirp")
			return
		}

		// Respond with the chirp
		responseWithJSON(w, http.StatusOK, chirp)
	}
}

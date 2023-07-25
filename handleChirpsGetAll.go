package main

import (
	"net/http"

	"github.com/veryordinary11/Go-web-server/database"
)

func handlerChirpsGetAll(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get all chirps from the database
		chirps, err := db.GetChirps()
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, "Failed to get chirps")
			return
		}

		// Respond with the chirps
		responseWithJSON(w, http.StatusOK, chirps)
	}
}

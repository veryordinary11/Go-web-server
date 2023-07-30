package main

import (
	"net/http"
	"strconv"

	"github.com/veryordinary11/Go-web-server/database"
)

func handlerChirpsGetAll(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Grab authorId from URL
		authorId := r.URL.Query().Get("author_id")
		// Get all chirps from the database
		chirps, err := db.GetChirps()
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, "Failed to get chirps")
			return
		}

		// Filter chirps by authorId
		if authorId == "" {
			responseWithJSON(w, http.StatusOK, chirps)
			return
		}

		filteredChirps := []database.Chirp{}
		authorIdInt, err := strconv.Atoi(authorId)
		if err != nil {
			responseWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		for _, chirp := range chirps {
			if chirp.AuthorID == authorIdInt {
				filteredChirps = append(filteredChirps, chirp)
			}
		}

		// Respond with the chirps
		responseWithJSON(w, http.StatusOK, filteredChirps)
	}
}

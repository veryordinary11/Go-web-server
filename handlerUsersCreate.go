package main

import (
	"encoding/json"
	"net/http"

	"github.com/veryordinary11/Go-web-server/database"
)

// POST /api/users
func handlerUsersCreate(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			responseWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		// Decode the request body into a User
		decoder := json.NewDecoder(r.Body)
		requestBody := database.User{}
		err := decoder.Decode(&requestBody)
		if err != nil {
			responseWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Email and password are valid
		user, err := db.CreateUser(requestBody.Email, requestBody.Password)
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Respond with the user
		responseWithJSON(w, http.StatusOK, user)
	}
}

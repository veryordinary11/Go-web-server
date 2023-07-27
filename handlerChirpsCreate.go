package main

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/veryordinary11/Go-web-server/database"
)

func handlerChirpsCreate(apiCfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			responseWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		// Validate the JWT token and get the user ID from the claims
		authToken := ExtractTokenFromHeader(*r)
		token, err := jwt.ParseWithClaims(authToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(apiCfg.jwtSecret), nil
		})
		if err != nil || !token.Valid {
			responseWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		claims, ok := token.Claims.(*jwt.MapClaims)
		if !ok {
			responseWithError(w, http.StatusUnauthorized, "Invalid JWT token")
			return
		}

		userId := (*claims)["sub"].(string)

		// Decode the request body into a Chirp
		decoder := json.NewDecoder(r.Body)
		requestBody := database.Chirp{}
		err = decoder.Decode(&requestBody)
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
		chirp, err := apiCfg.DB.CreateChirp(userId, cleanedText)
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, "Failed to create chirp")
			return
		}

		// Respond with the chirp
		responseWithJSON(w, http.StatusOK, chirp)
	}
}

package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
)

func handlerChirpsDelete(apiCfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the request header
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

		// Get the chirp ID from the URL
		chirpId := chi.URLParam(r, "id")

		// Delete the chirp from the database
		if err := apiCfg.DB.DeleteChirp(userId, chirpId); err != nil {
			responseWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Respond with a 204 status code
		w.WriteHeader(http.StatusNoContent)
	}
}

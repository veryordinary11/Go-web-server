package main

import (
	"net/http"

	"github.com/golang-jwt/jwt"
)

func handlerRevokeToken(apiCfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the refresh token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			responseWithError(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		// Strip off the "Bearer " prefix from the header
		refreshToken := extractAuthToken(authHeader)
		if refreshToken == "" {
			responseWithError(w, http.StatusUnauthorized, "Authorization header not found")
			return
		}

		// Validate the refresh token
		token, err := jwt.ParseWithClaims(refreshToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(apiCfg.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			responseWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Check if the token is a refresh token
		if claims, ok := token.Claims.(*jwt.MapClaims); ok {

			if (*claims)["iss"] != "chirpy-refresh" {
				responseWithError(w, http.StatusUnauthorized, "Invalid JWT token: Invalid type")
				return
			}
		} else {
			responseWithError(w, http.StatusUnauthorized, "Invalid JWT token: Type not found")
			return
		}

		// Revoke the refresh token in the database
		if err := apiCfg.DB.RevokedRefreshToken(refreshToken); err != nil {
			responseWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Respond with a 204 status code
		w.WriteHeader(http.StatusNoContent)
	}
}

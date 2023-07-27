package main

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

func handlerRefreshToken(apiCfg *apiConfig) http.HandlerFunc {
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

		// Check if the token is not revoked in the database
		if isRevoked, err := apiCfg.DB.IsRefreshTokenRevoked(refreshToken); err != nil {
			responseWithError(w, http.StatusInternalServerError, err.Error())
			return
		} else if isRevoked {
			responseWithError(w, http.StatusUnauthorized, "Invalid JWT token: Token is revoked")
			return
		}

		// Create a new access token
		userId := getUserIdFromToken(token)
		accessToken, err := CreateAccessToken(apiCfg.jwtSecret, userId)
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Send the access token back to the client
		responseWithJSON(w, http.StatusOK, map[string]string{
			"token": accessToken,
		})
	}
}

func getUserIdFromToken(token *jwt.Token) string {
	if claims, ok := token.Claims.(*jwt.MapClaims); ok {
		return (*claims)["sub"].(string)
	}
	return ""
}

func CreateAccessToken(jwtSecret string, userId string) (string, error) {
	expiresAt := time.Now().UTC().Add(time.Second * time.Duration(accessTokenExpireInSecondDefault)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"iss": "chirpy-access",
		"exp": expiresAt,
		"iat": time.Now().UTC().Unix(),
	})

	// Sign the token with the secret key
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

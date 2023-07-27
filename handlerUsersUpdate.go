package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/veryordinary11/Go-web-server/database"
)

func handlerUsersUpdate(apiCfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the JWT token from the Authorization header
		authToken := ExtractTokenFromHeader(*r)

		// Validate the JWT token
		token, err := jwt.ParseWithClaims(authToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(apiCfg.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			responseWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Check the token
		if claims, ok := token.Claims.(*jwt.MapClaims); ok {

			// Check the expiration time
			if exp, ok := (*claims)["exp"].(float64); ok {
				if time.Now().UTC().Unix() > int64(exp) {
					responseWithError(w, http.StatusUnauthorized, "JWT token is expired")
					return
				}
			} else {
				responseWithError(w, http.StatusUnauthorized, "Invalid JWT token: Expiration time not found")
				return
			}

			// Check the issuer
			if iss, ok := (*claims)["iss"].(string); ok {
				if iss != "chirpy-access" {
					responseWithError(w, http.StatusUnauthorized, "Invalid JWT token: Invalid issuer")
					return
				}
			} else {
				responseWithError(w, http.StatusUnauthorized, "Invalid JWT token: Issuer not found")
				return
			}

		} else {
			fmt.Println("Invalid token claims:", token.Claims)
			responseWithError(w, http.StatusUnauthorized, "Invalid JWT token")
			return
		}

		// Get the user ID from the claims
		userID := (*token.Claims.(*jwt.MapClaims))["sub"].(string)

		// Parse the request body to get the updated user information
		var requestBody struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err = json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			responseWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Update the user in the database
		updatedUser, err := apiCfg.DB.UpdateUser(userID, requestBody.Email, requestBody.Password)
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Respond with the updated user resource without password
		responseWithJSON(w, http.StatusOK, database.UserWithoutPassword{
			ID:    updatedUser.ID,
			Email: updatedUser.Email,
		})
	}
}

// Helper function to extract the JWT token from the Authorization header
func extractAuthToken(authHeader string) string {
	const bearerPrefix = "Bearer "
	if len(authHeader) > len(bearerPrefix) && authHeader[:len(bearerPrefix)] == bearerPrefix {
		return authHeader[len(bearerPrefix):]
	}

	return ""
}

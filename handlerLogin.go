package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/veryordinary11/Go-web-server/database"
	"golang.org/x/crypto/bcrypt"
)

func handlerLogin(db *database.DB, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			responseWithError(w, http.StatusMethodNotAllowed, "Method not matched, must be POST")
			return
		}

		// Decode the request body into a UserLoginRequest
		decoder := json.NewDecoder(r.Body)
		requestBody := database.UserLoginRequest{}
		err := decoder.Decode(&requestBody)
		if err != nil {
			responseWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Lookup the user by email in the database
		user, err := db.GetUserByEmail(requestBody.Email)
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Check if the password is correct
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
		if err != nil {
			responseWithError(w, http.StatusUnauthorized, "Password is incorrect")
			return
		}

		// Generate the JWT token
		var tokenExpireInSecond int

		// Check if the expiresInSecond is valid(0~86400)
		if requestBody.ExpiresInSecond != nil && *requestBody.ExpiresInSecond > 0 && *requestBody.ExpiresInSecond <= 86400 {
			tokenExpireInSecond = *requestBody.ExpiresInSecond
		} else {
			// Default expiration time is 24 hours
			tokenExpireInSecond = 86400
		}

		expiresAt := time.Now().UTC().Add(time.Second * time.Duration(tokenExpireInSecond)).Unix()

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": strconv.Itoa(user.ID),
			"iss": "chirpy",
			"exp": expiresAt,
			"iat": time.Now().UTC().Unix(),
		})

		// Sign the token with the secret key
		signedToken, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Respond with the userWithoutPassword
		responseWithJSON(w, http.StatusOK, struct {
			ID    int    `json:"id"`
			Email string `json:"email"`
			Token string `json:"token"`
		}{
			ID:    user.ID,
			Email: user.Email,
			Token: signedToken,
		})
	}
}

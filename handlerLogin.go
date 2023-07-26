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

		// Generate the JWT access token
		const accessTokenExpireInSecondDefault = 60 * 60 // 1 hour
		expiresAt := time.Now().UTC().Add(time.Second * time.Duration(accessTokenExpireInSecondDefault)).Unix()

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": strconv.Itoa(user.ID),
			"iss": "chirpy-access",
			"exp": expiresAt,
			"iat": time.Now().UTC().Unix(),
		})

		// Sign the token with the secret key
		signedToken, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Generate the JWT refresh token
		const refreshTokenExpireInSecondDefault = 60 * 60 * 24 * 60 // 60 days
		refreshTokenExpiresAt := time.Now().UTC().Add(time.Second * time.Duration(refreshTokenExpireInSecondDefault)).Unix()

		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": strconv.Itoa(user.ID),
			"iss": "chirpy-refresh",
			"exp": refreshTokenExpiresAt,
			"iat": time.Now().UTC().Unix(),
		})

		// Sign the refresh token with the secret key
		signedRefreshToken, err := refreshToken.SignedString([]byte(jwtSecret))
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Respond with the userWithoutPassword
		responseWithJSON(w, http.StatusOK, struct {
			ID           int    `json:"id"`
			Email        string `json:"email"`
			Token        string `json:"token"`
			RefreshToken string `json:"refreshToken"`
		}{
			ID:           user.ID,
			Email:        user.Email,
			Token:        signedToken,
			RefreshToken: signedRefreshToken,
		})
	}
}

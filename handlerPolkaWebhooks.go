package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type polkaWebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserId int `json:"user_id"`
	} `json:"data"`
}

const APIKEY_PREFIX = "Apikey "

func handlerPolkaWebhooks(apiCfg *apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			responseWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		// Verify the request is from Polka by checking the Authorization header
		authorization := r.Header.Get("Authorization")
		if authorization != APIKEY_PREFIX+apiCfg.polkaKey {
			responseWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Decode the request body
		decoder := json.NewDecoder(r.Body)
		requestBody := polkaWebhookRequest{}
		err := decoder.Decode(&requestBody)
		if err != nil {
			responseWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Update the user's IsChirpyRed status
		userId := requestBody.Data.UserId
		user, err := apiCfg.DB.GetUserByID(strconv.Itoa(userId))
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		err = apiCfg.DB.UpdateUserIsChirpyRed(user)
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Return a success response
		responseWithJSON(w, http.StatusOK, map[string]string{"status": "success"})
	}
}

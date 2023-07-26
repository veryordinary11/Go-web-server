package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func createAPIRouter(apiCfg *apiConfig) http.Handler {
	r := chi.NewRouter()

	// Test API
	r.Get("/healthz", handlerReadiness)
	r.Get("/metrics", apiCfg.handlerMetrics)
	r.Post("/validate_chirp", handlerValidation)

	// Chirps API
	r.Get("/chirps", handlerChirpsGetAll(apiCfg.DB))
	r.Get("/chirps/{id}", handlerChirpsGet(apiCfg.DB))
	r.Post("/chirps", handlerChirpsCreate(apiCfg.DB))

	// Users API
	r.Post("/users", handlerUsersCreate(apiCfg.DB))
	r.Post("/login", handlerLogin(apiCfg.DB, apiCfg.jwtSecret))
	r.Put("/users", handlerUsersUpdate(apiCfg))

	// Token API
	r.Post("/refresh", handlerRefreshToken(apiCfg.jwtSecret))
	r.Post("/revoke", handlerRevokeToken(apiCfg.jwtSecret))

	return r
}

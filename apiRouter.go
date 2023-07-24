package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func createAPIRouter(apiCfg *apiConfig) http.Handler {
	r := chi.NewRouter()
	r.Get("/healthz", handlerReadiness)
	r.Get("/metrics", apiCfg.handlerMetrics)
	r.Post("/validate_chirp", handlerValidation)
	return r
}

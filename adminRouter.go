package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func createAdminRouter(apiCfg *apiConfig) http.Handler {
	r := chi.NewRouter()
	r.Get("/metrics", apiCfg.handlerMetrics)
	return r
}

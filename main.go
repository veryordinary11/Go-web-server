package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := &apiConfig{
		fileserverHits: 0,
	}

	apiRouter := createAPIRouter(apiCfg)

	adminRouter := createAdminRouter(apiCfg)

	r := chi.NewRouter()

	r.Mount("/api", apiRouter)

	r.Mount("/admin", adminRouter)

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)

	corsRouter := middlewareCors(r)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsRouter,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

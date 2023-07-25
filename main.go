package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/veryordinary11/Go-web-server/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	// Create a new database connection
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatalf("Error creating database: %v", err)
	}

	apiCfg := &apiConfig{
		fileserverHits: 0,
		DB:             db,
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

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

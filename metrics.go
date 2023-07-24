package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	templatePath := "admin_metrics.html"

	tmpl, err := template.New("admin_metrics").ParseFiles(templatePath)
	if err != nil {
		fmt.Println(err)
		defer http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Hits int
	}{
		Hits: cfg.fileserverHits,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		defer http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

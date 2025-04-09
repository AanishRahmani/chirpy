package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		c.fileServerHits.Add(1)
		next.ServeHTTP(writer, request)
	})
}

func (c *apiConfig) handleHits(w http.ResponseWriter, _ *http.Request) {

	hits := c.fileServerHits.Load()

	templatePath := filepath.Join("html", "admin.html")

	templ, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "failed to parse template", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = templ.Execute(w, struct {
		Hits int32
	}{
		Hits: hits,
	})

	if err != nil {
		http.Error(w, "failed to remder template", http.StatusInternalServerError)
	}

	//doesnt load through html
	//ans := fmt.Sprintf("Hits: %d", hits)
	//_, err := w.Write([]byte(ans))
	//w.WriteHeader(http.StatusOK)
	//w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	//if err != nil {
	//	return
	//}
}

func (c *apiConfig) resetHits() {
	c.fileServerHits.Store(0)
}

func (c *apiConfig) resetHitsHandler(w http.ResponseWriter, _ *http.Request) {
	c.resetHits()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

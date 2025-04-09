package main

import (
	"fmt"
	"net/http"
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
	ans := fmt.Sprintf("Hits: %d", hits)
	_, err := w.Write([]byte(ans))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if err != nil {
		return
	}
}

func (c *apiConfig) resetHits() {
	c.fileServerHits.Store(0)
}

func (c *apiConfig) resetHitsHandler(w http.ResponseWriter, _ *http.Request) {
	c.resetHits()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

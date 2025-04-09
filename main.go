package main

import (
	"fmt"
	"net/http"
)

func main() {

	fmt.Println("Starting server on 8080")
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	cfg := &apiConfig{}

	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir("./html")))
	mux.Handle("/app", cfg.middlewareMetricsInc(fileServer)) //mkae sure to give it a directory

	imageServer := http.StripPrefix("/app/assets/", http.FileServer(http.Dir("images")))
	mux.Handle("/app/assets/", cfg.middlewareMetricsInc(imageServer))

	mux.Handle("GET /admin/healthz", http.HandlerFunc(isReady))

	mux.HandleFunc("GET /admin/metrics", cfg.handleHits)

	mux.HandleFunc("POST /admin/reset", cfg.resetHitsHandler)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

}

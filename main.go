package main

import (
	"log"
	"net/http"
)

const (
	filePathRoot = "."
	serverPort   = "8080"
)

type apiConfig struct {
	fileServerHits int
}

func main() {
	log.Printf("starting server on port %v\n", serverPort)

	cfg := &apiConfig{}
	mux := http.NewServeMux()
	mux.Handle("/app/*", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("GET /api/reset", cfg.handlerReset)

	server := http.Server{
		Addr:    "localhost:" + serverPort,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

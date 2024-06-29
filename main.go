package main

import (
	"fmt"
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

func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.fileServerHits++
		next.ServeHTTP(w, r)
	})
}

func (a *apiConfig) handlerReset(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(200)
	a.fileServerHits = 0
}

func (a *apiConfig) handlerMetrics(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(200)
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err := writer.Write([]byte(fmt.Sprintf("Hits: %v", a.fileServerHits)))
	if err != nil {
		log.Printf("metrics error: %s", err.Error())
	}
}

func main() {
	log.Printf("starting server on port %v\n", serverPort)

	cfg := &apiConfig{}
	mux := http.NewServeMux()
	mux.Handle("/app/*", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot)))))
	mux.HandleFunc("GET /healthz", handlerReadiness)
	mux.HandleFunc("GET /metrics", cfg.handlerMetrics)
	mux.HandleFunc("GET /reset", cfg.handlerReset)

	server := http.Server{
		Addr:    "localhost:" + serverPort,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

func handlerReadiness(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	_, err := writer.Write([]byte("OK"))
	if err != nil {
		log.Printf("writer error: %s", err.Error())
	}
}

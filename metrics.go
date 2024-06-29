package main

import (
	"fmt"
	"log"
	"net/http"
)

func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.fileServerHits++
		next.ServeHTTP(w, r)
	})
}

func (a *apiConfig) handlerMetrics(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(200)
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err := writer.Write([]byte(fmt.Sprintf("Hits: %v", a.fileServerHits)))
	if err != nil {
		log.Printf("metrics error: %s", err.Error())
	}
}

func (a *apiConfig) handlerReset(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(200)
	a.fileServerHits = 0
}

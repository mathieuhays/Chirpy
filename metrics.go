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
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := fmt.Sprintf(`<html>
<body>
<h1>Welcome, Chirpy Admin</h1>
<p>Chirpy has been visited %d times!</p>
</body>
</html>`, a.fileServerHits)
	_, err := writer.Write([]byte(html))
	if err != nil {
		log.Printf("metrics error: %s", err.Error())
	}
}

func (a *apiConfig) handlerReset(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(200)
	a.fileServerHits = 0
}

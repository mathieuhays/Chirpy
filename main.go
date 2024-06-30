package main

import (
	"encoding/json"
	"github.com/mathieuhays/Chirpy/internal/database"
	"log"
	"net/http"
)

const (
	filePathRoot = "."
	serverPort   = "8080"
)

type apiConfig struct {
	fileServerHits int
	database       *database.DB
}

func main() {
	log.Printf("starting server on port %v\n", serverPort)

	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Fatal(err)
	}

	cfg := &apiConfig{
		fileServerHits: 0,
		database:       db,
	}
	mux := http.NewServeMux()
	mux.Handle("/app/*", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("GET /api/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", cfg.handlerPostChirps)
	mux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirp)
	mux.HandleFunc("POST /api/users", cfg.handlerPostUsers)

	server := http.Server{
		Addr:    "localhost:" + serverPort,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, err error, statusCode int) {
	writeJSON(w, statusCode, errorResponse{Error: err.Error()})
}

func writeJSON(w http.ResponseWriter, statusCode int, object interface{}) {
	data, mErr := json.Marshal(object)
	if mErr != nil {
		log.Printf("error marshalling error object: %v", mErr.Error())
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

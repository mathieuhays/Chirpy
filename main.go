package main

import (
	"log"
	"net/http"
)

const (
	serverPort = "8080"
)

func main() {
	log.Printf("starting server on port %v\n", serverPort)

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		writer.WriteHeader(200)
		_, err := writer.Write([]byte("OK"))
		if err != nil {
			log.Printf("writer error: %s", err.Error())
		}
	})

	server := http.Server{
		Addr:    "localhost:" + serverPort,
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("server error: %s\n", err.Error())
	}

	log.Println("shutting down")
}

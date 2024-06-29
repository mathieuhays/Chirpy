package main

import (
	"log"
	"net/http"
)

const (
	filePathRoot = "."
	serverPort   = "8080"
)

func main() {
	log.Printf("starting server on port %v\n", serverPort)

	mux := http.NewServeMux()
	mux.Handle("/app/*", http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot))))
	mux.HandleFunc("/healthz", handlerReadiness)

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

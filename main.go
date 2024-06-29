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
	//mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
	//	if request.URL.Path != "/" {
	//		http.NotFound(writer, request)
	//		return
	//	}
	//
	//	fmt.Fprintf(writer, "Welcome\n")
	//})

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

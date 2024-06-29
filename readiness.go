package main

import (
	"log"
	"net/http"
)

func handlerReadiness(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	_, err := writer.Write([]byte("OK"))
	if err != nil {
		log.Printf("writer error: %s", err.Error())
	}
}

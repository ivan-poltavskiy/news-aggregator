package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/fetch-articles", WebHandler)
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}

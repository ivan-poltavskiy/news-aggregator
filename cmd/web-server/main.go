package main

import (
	"log"
	"net/http"
)

func main() {
	const PORT = ":8080"
	http.HandleFunc("/articles/fetch", FetchArticleHandler)
	http.HandleFunc("/sources/add", AddSourceHandler)
	http.HandleFunc("/sources/delete", DeleteSourceByNameHandler)
	log.Println("Starting server on " + PORT)

	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}

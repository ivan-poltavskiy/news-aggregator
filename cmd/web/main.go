package main

import (
	"log"
	"net/http"
	"news-aggregator/cmd/web/handlers"
)

func main() {
	const PORT = ":8080"
	http.HandleFunc("GET /articles/fetch", handlers.FetchArticleHandler)
	http.HandleFunc("POST /sources/add", handlers.AddSourceHandler)
	http.HandleFunc("DELETE /sources/delete", handlers.DeleteSourceByNameHandler)
	log.Println("Starting server on " + PORT)

	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}

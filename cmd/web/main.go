package web

import (
	"log"
	"net/http"
	"news-aggregator/cmd/web/handlers"
)

func main() {
	const PORT = ":8080"
	http.HandleFunc("/articles/fetch", handlers.FetchArticleHandler)
	http.HandleFunc("/sources/add", handlers.AddSourceHandler)
	http.HandleFunc("/sources/delete", handlers.DeleteSourceByNameHandler)
	log.Println("Starting server on " + PORT)

	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}

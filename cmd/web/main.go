package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"news-aggregator/cmd/web/handlers"
	"news-aggregator/constant"
)

func main() {

	http.HandleFunc("GET /news", handlers.FetchNewsHandler)
	http.HandleFunc("POST /sources", handlers.AddSourceHandler)
	http.HandleFunc("DELETE /sources", handlers.DeleteSourceByNameHandler)
	logrus.Info("Starting server on " + constant.PORT)

	err := http.ListenAndServeTLS(constant.PORT, constant.CertFile, constant.KeyFile, nil)
	if err != nil {
		logrus.Fatalf("Could not start server: %s\n", err.Error())
	}
}

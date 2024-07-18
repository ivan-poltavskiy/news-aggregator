package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"news-aggregator/aggregator"
	"news-aggregator/cmd/web/handlers"
	"news-aggregator/collector"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
	"news-aggregator/storage"
)

func main() {
	sourceStorage := storage.NewJsonSourceStorage(source.PathToFile(constant.PathToStorage))
	newsCollector := collector.New(sourceStorage)
	newsAggregator := aggregator.New(newsCollector)

	http.HandleFunc("GET /news", func(w http.ResponseWriter, r *http.Request) {
		handlers.FetchNewsHandler(w, r, sourceStorage, newsAggregator)
	})
	http.HandleFunc("POST /sources", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddSourceHandler(w, r, sourceStorage)
	})
	http.HandleFunc("DELETE /sources", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteSourceByNameHandler(w, r, sourceStorage)
	})
	logrus.Info("Starting server on " + constant.PORT)

	err := http.ListenAndServeTLS(constant.PORT, constant.CertFile, constant.KeyFile, nil)
	if err != nil {
		logrus.Fatalf("Could not start server: %s\n", err.Error())
	}
}

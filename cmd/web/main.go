package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"net/http"
	"news-aggregator/aggregator"
	"news-aggregator/cmd/web/handlers"
	"news-aggregator/cmd/web/service"
	"news-aggregator/collector"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
	"news-aggregator/storage"
)

func main() {

	port := flag.String("port", constant.PORT, "port to listen on")
	pathToCertificate := flag.String("path to certificate", constant.CertFile, "Certificate file path")
	pathToKey := flag.String("path to key", constant.KeyFile, "Key file path")
	flag.Parse()

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
	logrus.Info("Starting server on " + *port)

	err := http.ListenAndServeTLS(*port, *pathToCertificate, *pathToKey, nil)
	if err != nil {
		logrus.Fatalf("Could not start server: %s\n", err.Error())
	}

	go service.PeriodicallyUpdateNews(sourceStorage)
	select {}
}

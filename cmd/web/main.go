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
	"time"
)

func main() {

	port := flag.String("port", constant.PORT, "port to listen on")
	pathToCertificate := flag.String("pathToCertificate", constant.PathToCertFile, "Certificate file path")
	pathToKey := flag.String("pathToKey", constant.PathToKeyFile, "Key file path")
	newsUpdatePeriod := flag.Int("newsUpdatePeriod", constant.NewsUpdatePeriodIOnMinutes, "Period of time in minutes for periodically news updating")
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

	go func() {
		err := http.ListenAndServeTLS(*port, *pathToCertificate, *pathToKey, nil)
		if err != nil {
			logrus.Fatalf("Could not start server: %s\n", err.Error())
		}
	}()

	logrus.Info("Starting periodic news update every ", *newsUpdatePeriod, " minutes")

	go service.PeriodicallyUpdateNews(sourceStorage, time.Duration(*newsUpdatePeriod)*time.Minute)

	select {}

}

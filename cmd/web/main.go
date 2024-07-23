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
	newsStorage "news-aggregator/storage/news"
	sourceStorage "news-aggregator/storage/source"
	"time"
)

func main() {

	port := flag.String("port", constant.PORT, "port to listen on")
	pathToCertificate := flag.String("pathToCertificate", constant.PathToCertFile, "Certificate file path")
	pathToKey := flag.String("pathToKey", constant.PathToKeyFile, "Key file path")
	newsUpdatePeriod := flag.Int("newsUpdatePeriod", constant.NewsUpdatePeriodIOnMinutes, "Period of time in minutes for periodically news updating")
	flag.Parse()

	newsJsonStorage := newsStorage.NewJsonStorage(source.PathToFile(constant.PathToResources))
	sourceJsonStorage := sourceStorage.NewJsonSourceStorage(source.PathToFile(constant.PathToStorage))
	newsCollector := collector.New(sourceJsonStorage)
	newsAggregator := aggregator.New(newsCollector)

	http.HandleFunc("GET /news", func(w http.ResponseWriter, r *http.Request) {
		handlers.FetchNewsHandler(w, r, newsAggregator)
	})
	http.HandleFunc("POST /sources", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddSourceHandler(w, r, sourceJsonStorage, newsJsonStorage)
	})
	http.HandleFunc("DELETE /sources", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteSourceByNameHandler(w, r, sourceJsonStorage)
	})
	logrus.Info("Starting server on " + *port)

	go func() {
		err := http.ListenAndServeTLS(*port, *pathToCertificate, *pathToKey, nil)
		if err != nil {
			logrus.Fatalf("Could not start server: %s\n", err.Error())
		}
	}()

	logrus.Info("Starting periodic news update every ", *newsUpdatePeriod, " minutes")

	go service.PeriodicallyUpdateNews(sourceJsonStorage, time.Duration(*newsUpdatePeriod)*time.Minute, newsJsonStorage)

	select {}

}

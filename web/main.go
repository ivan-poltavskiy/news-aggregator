package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"net/http"
	"news-aggregator/aggregator"
	"news-aggregator/collector"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
	"news-aggregator/storage"
	newsStorage "news-aggregator/storage/news"
	sourceStorage "news-aggregator/storage/source"
	"news-aggregator/web/handlers"
	"news-aggregator/web/news"
	"time"
)

func main() {

	port := flag.String("port", constant.PORT, "port to listen on")
	pathToCertificate := flag.String("certificate-path", constant.PathToCertFile, "Certificate file path")
	pathToKey := flag.String("key-path", constant.PathToKeyFile, "Key file path")
	newsUpdatePeriod := flag.Int("news-update-period", constant.NewsUpdatePeriodIOnMinutes, "Period of time in minutes for periodically news updating")
	flag.Parse()

	newsJsonStorage, err := newsStorage.NewJsonNewsStorage(source.PathToFile(constant.PathToResources))
	if err != nil {
		logrus.Fatal(err)
	}
	sourceJsonStorage, err := sourceStorage.NewJsonSourceStorage(source.PathToFile(constant.PathToStorage))
	if err != nil {
		logrus.Fatal(err)
	}

	newStorage := storage.NewStorage(newsJsonStorage, sourceJsonStorage)

	newsCollector := collector.New(newStorage)
	newsAggregator := aggregator.New(newsCollector)

	http.HandleFunc("GET /news", func(w http.ResponseWriter, r *http.Request) {
		handlers.FetchNewsHandler(w, r, newsAggregator)
	})
	http.HandleFunc("POST /sources", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddSourceHandler(w, r, newStorage)
	})
	http.HandleFunc("DELETE /sources", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteSourceByNameHandler(w, r, newStorage)
	})
	logrus.Info("Starting server on: " + *port)

	go func() {
		err := http.ListenAndServeTLS(":"+*port, *pathToCertificate, *pathToKey, nil)
		if err != nil {
			logrus.Fatalf("Could not start server: %s\n", err.Error())
		}
	}()

	logrus.Info("Starting periodic news update every ", *newsUpdatePeriod, " minutes")
	service := news.NewNewsService(newStorage)
	go service.PeriodicallyUpdateNews(time.Duration(*newsUpdatePeriod) * time.Minute)

	select {}

}

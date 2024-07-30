package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"net/http"
	"news-aggregator/aggregator"
	"news-aggregator/client"
	"news-aggregator/collector"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
	"news-aggregator/storage"
	newsStorage "news-aggregator/storage/news"
	sourceStorage "news-aggregator/storage/source"
	handler2 "news-aggregator/web/handler"
	"news-aggregator/web/news"
	"time"
)

func main() {

	port := flag.String("port", constant.PORT, "port to listen on")
	pathToCertificate := flag.String("certificate-path", constant.PathToCertFile, "Certificate file path")
	pathToKey := flag.String("key-path", constant.PathToKeyFile, "Key file path")
	newsUpdatePeriod := flag.Int("news-update-period", constant.NewsUpdatePeriodIOnMinutes, "Period of time in minutes for periodically news updating")
	flag.Parse()

	newsJsonStorage, err := newsStorage.NewJsonStorage(source.PathToFile(constant.PathToResources))
	if err != nil {
		logrus.Fatal(err)
	}
	sourceJsonStorage, err := sourceStorage.NewJsonStorage(source.PathToFile(constant.PathToStorage))
	if err != nil {
		logrus.Fatal(err)
	}

	resourcesStorage := storage.NewStorage(newsJsonStorage, sourceJsonStorage)

	newsCollector := collector.New(resourcesStorage)
	newsAggregator := aggregator.New(newsCollector)

	handler := handler2.NewHandler(resourcesStorage)

	http.HandleFunc("GET /news", func(w http.ResponseWriter, r *http.Request) {
		handler.GetNewsHandler().FetchNewsHandler(w, client.NewWebClient(*r, w, newsAggregator))
	})
	http.HandleFunc("POST /sources", func(w http.ResponseWriter, r *http.Request) {
		handler.GetSourceHandler().AddSourceHandler(w, r)
	})
	http.HandleFunc("DELETE /sources", func(w http.ResponseWriter, r *http.Request) {
		handler.GetSourceHandler().DeleteSourceByNameHandler(w, r)
	})
	logrus.Info("Starting server on: " + *port)

	logrus.Info("Starting periodic news update every ", *newsUpdatePeriod, " minutes")
	service := news.NewService(resourcesStorage)
	go service.PeriodicallyUpdateNews(time.Duration(*newsUpdatePeriod) * time.Minute)

	err = http.ListenAndServeTLS(":"+*port, *pathToCertificate, *pathToKey, nil)
	if err != nil {
		logrus.Fatalf("Could not start server: %s\n", err.Error())
	}

}

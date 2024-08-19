package main

import (
	"crypto/tls"
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
	"path/filepath"
)

func main() {

	port := flag.String("port", constant.PORT, "port to listen on")
	secretPath := flag.String("secret-path", "/etc/tls-secret", "Path to TLS Secret")
	flag.Parse()

	certPath := filepath.Join(*secretPath, "tls.crt")
	keyPath := filepath.Join(*secretPath, "tls.key")

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		logrus.Fatalf("Failed to load key pair: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	server := &http.Server{
		Addr:      ":" + *port,
		TLSConfig: tlsConfig,
	}

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

	handler := NewHandler(resourcesStorage)

	http.HandleFunc("GET /news", func(w http.ResponseWriter, r *http.Request) {
		handler.GetNewsHandler().FetchNewsHandler(w, client.NewWebClient(*r, w, newsAggregator))
	})
	http.HandleFunc("POST /sources", func(w http.ResponseWriter, r *http.Request) {
		handler.GetSourceHandler().AddSourceHandler(w, r)
	})
	http.HandleFunc("DELETE /sources", func(w http.ResponseWriter, r *http.Request) {
		handler.GetSourceHandler().DeleteSourceByNameHandler(w, r)
	})
	http.HandleFunc("PUT /sources", func(w http.ResponseWriter, r *http.Request) {
		handler.GetSourceHandler().UpdateSourceByName(w, r)
	})
	http.HandleFunc("GET /allSources", func(w http.ResponseWriter, r *http.Request) {
		handler.GetSourceHandler().GetAllSources(w)
	})
	logrus.Info("Starting server on: " + *port)

	logrus.Infof("Starting server on port %s", *port)
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		logrus.Fatalf("Could not start server: %v", err)
	}

}

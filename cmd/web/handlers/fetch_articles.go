package handlers

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"news-aggregator/aggregator"
	"news-aggregator/client"
	"news-aggregator/collector"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
	"news-aggregator/storage"
)

// FetchNewsHandler handles requests for fetching news.
func FetchNewsHandler(w http.ResponseWriter, r *http.Request, storage storage.Storage) {
	sources, err := source.LoadExistingSourcesFromStorage(constant.PathToStorage)
	if err != nil {
		http.Error(w, "Failed to load sources: "+err.Error(), http.StatusInternalServerError)
		return
	}
	newsCollector := collector.New(sources)
	newsAggregator := aggregator.New(newsCollector)

	webClient := client.NewWebClient(*r, w, newsAggregator, storage)
	news, err := webClient.FetchNews()
	if err != nil {
		logrus.Error("Failed to fetch news ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	webClient.Print(news)
}

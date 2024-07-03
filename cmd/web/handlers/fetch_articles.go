package handlers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"news-aggregator/aggregator"
	"news-aggregator/client"
	"news-aggregator/collector"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
)

// FetchArticleHandler handles HTTP requests for fetching articles.
func FetchArticleHandler(w http.ResponseWriter, r *http.Request) {
	sources, err := source.LoadExistingSourcesFromStorage(constant.PathToStorage)
	if err != nil {
		http.Error(w, "Failed to load sources: "+err.Error(), http.StatusInternalServerError)
		return
	}
	articleCollector := collector.New(sources)
	newsAggregator := aggregator.New(articleCollector)

	webClient := client.NewWebClient(*r, newsAggregator)
	articles, err := webClient.FetchArticles()
	if err != nil {
		logrus.Error("Failed to fetch articles ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(articles)
	if err != nil {
		logrus.Error("Failed to encode json ", err)
		return
	}
}

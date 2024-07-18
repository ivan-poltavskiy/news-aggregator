package handlers

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"news-aggregator/client"
	"news-aggregator/storage"
)

// FetchNewsHandler handles requests for fetching news.
func FetchNewsHandler(w http.ResponseWriter, r *http.Request, storage storage.Storage, newsAggregator client.Aggregator) {

	webClient := client.NewWebClient(*r, w, newsAggregator, storage)
	news, err := webClient.FetchNews()
	if err != nil {
		logrus.Error("Failed to fetch news ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	webClient.Print(news)
}

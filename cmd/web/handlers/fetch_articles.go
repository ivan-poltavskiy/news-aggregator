package handlers

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"news-aggregator/client"
)

// FetchNewsHandler handles requests for fetching news.
func FetchNewsHandler(w http.ResponseWriter, r *http.Request, newsAggregator client.Aggregator) {

	webClient := client.NewWebClient(*r, w, newsAggregator)
	news, err := webClient.FetchNews()
	if err != nil {
		logrus.Error("Failed to fetch news ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	webClient.Print(news)
}

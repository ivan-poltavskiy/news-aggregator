package main

import (
	"encoding/json"
	"net/http"
	"news-aggregator/aggregator"
	"news-aggregator/client"
	"news-aggregator/collector"
	"news-aggregator/entity/source"
)

// WebHandler handles HTTP requests for fetching articles.
func WebHandler(w http.ResponseWriter, r *http.Request) {

	sources := []source.Source{
		{Name: "bbc", PathToFile: "./resources/bbc-world-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "nbc", PathToFile: "./resources/nbc-news.json", SourceType: "JSON"},
		{Name: "abc", PathToFile: "./resources/abcnews-international-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "washington", PathToFile: "./resources/washingtontimes-world-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "usatoday", PathToFile: "./resources/usatoday-world-news.html", SourceType: "UsaToday"},
	}
	articleCollector := collector.New(sources)
	newsAggregator := aggregator.New(articleCollector)

	webClient := client.NewWebClient(*r, newsAggregator)

	articles, err := webClient.FetchArticles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"news-aggregator/aggregator"
	"news-aggregator/client"
	"news-aggregator/collector"
	"news-aggregator/entity/source"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// FetchArticleHandler handles HTTP requests for fetching articles.
func FetchArticleHandler(w http.ResponseWriter, r *http.Request) {

	sources := []source.Source{
		{Name: "bbc", PathToFile: "./resources/bbc-world-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "nbc", PathToFile: "./resources/nbc-news.json", SourceType: "JSON"},
		{Name: "abc", PathToFile: "./resources/abcnews-international-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "washington", PathToFile: "./resources/washingtontimes-world-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "usatoday", PathToFile: "./resources/usatoday-world-news.html", SourceType: "UsaToday"},
		{Name: "nyt", PathToFile: "./feeds/2024-06-24/feed.xml", SourceType: "RSS"},
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
	err = json.NewEncoder(w).Encode(articles)
	if err != nil {
		return
	}
}

func AddSourceHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	err, rssURL := getRssFeedLink(w, url)
	if rssURL == "" {
		return
	}

	rssResp, err := http.Get(rssURL)
	if err != nil || rssResp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to download RSS feed", http.StatusInternalServerError)
		return
	}
	defer rssResp.Body.Close()

	currentDate := time.Now().Format("2006-01-02")

	dirPath := filepath.Join("feeds", currentDate)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		return
	}

	fileName := extractDomainName(url) + ".xml"

	filePath := filepath.Join(dirPath, fileName)
	outputFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, rssResp.Body)
	if err != nil {
		http.Error(w, "Failed to save RSS feed", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "RSS feed downloaded and saved to %s", filePath)
}

func getRssFeedLink(w http.ResponseWriter, url string) (error, string) {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to download page", http.StatusInternalServerError)
		return err, ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read page content", http.StatusInternalServerError)
		return err, ""
	}

	re := regexp.MustCompile(`(?i)<link[^>]+type="application/rss\+xml"[^>]+href="([^"]+)"`)
	matches := re.FindStringSubmatch(string(body))

	if len(matches) < 2 {
		http.Error(w, "RSS link not found", http.StatusBadRequest)
		return nil, ""
	}

	rssURL := matches[1]
	return nil, rssURL
}

func extractDomainName(url string) string {
	re := regexp.MustCompile(`https?://(www\.)?([^/]+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 3 {
		return "unknown"
	}
	domain := matches[2]
	domain = strings.Split(domain, ".")[0]
	return domain
}

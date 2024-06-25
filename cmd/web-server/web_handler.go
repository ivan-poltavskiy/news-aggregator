package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

	sources, err := LoadSourcesFromFile("./storage/sources-storage.json")
	if err != nil {
		http.Error(w, "Failed to load sources: "+err.Error(), http.StatusInternalServerError)
		return
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

	sourceEntity := source.Source{
		Name:       source.Name(extractDomainName(url)),
		PathToFile: source.PathToFile(filePath),
		SourceType: source.RSS,
	}

	if !sourceExists(sourceEntity.Name) {
		addSourceToFile(sourceEntity)
		fmt.Fprintf(w, "RSS feed downloaded and saved to %s and source added", filePath)
	} else {
		fmt.Fprintf(w, "RSS feed downloaded and saved to %s but source already exists", filePath)
	}
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

func sourceExists(name source.Name) bool {
	sources := readSourcesFromFile()
	for _, s := range sources {
		if s.Name == name {
			return true
		}
	}
	return false
}

func readSourcesFromFile() []source.Source {
	file, err := os.Open("./storage/sources-storage.json")
	if err != nil {
		if os.IsNotExist(err) {
			return []source.Source{}
		}
		fmt.Println("Error opening sources file:", err)
		return nil
	}
	defer file.Close()

	var sources []source.Source
	if err := json.NewDecoder(file).Decode(&sources); err != nil {
		fmt.Println("Error decoding sources file:", err)
		return nil
	}
	return sources
}

func addSourceToFile(newSource source.Source) {
	sources := readSourcesFromFile()
	sources = append(sources, newSource)

	file, err := os.Create("./storage/sources-storage.json")
	if err != nil {
		fmt.Println("Error creating sources file:", err)
		return
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(sources); err != nil {
		fmt.Println("Error encoding sources to file:", err)
	}
}

// LoadSourcesFromFile loads sources from a JSON file
func LoadSourcesFromFile(filename string) ([]source.Source, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	value, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var sources []source.Source
	err = json.Unmarshal(value, &sources)
	if err != nil {
		return nil, err
	}

	return sources, nil
}

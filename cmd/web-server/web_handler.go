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

type addSourceRequest struct {
	URL string `json:"url"`
}

type deleteSourceRequest struct {
	Name string `json:"name"`
}

// FetchArticleHandler handles HTTP requests for fetching articles.
func FetchArticleHandler(w http.ResponseWriter, r *http.Request) {

	sources, err := source.LoadExistingSourcesFromStorage("./storage/sources-storage.json")
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
	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var requestBody addSourceRequest
	err = json.Unmarshal(body, &requestBody)
	if err != nil || requestBody.URL == "" {
		http.Error(w, "Invalid request body or URL parameter is missing", http.StatusBadRequest)
		return
	}

	url := requestBody.URL

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

func DeleteSourceByNameHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var request deleteSourceRequest
	err = json.Unmarshal(body, &request)
	if err != nil || request.Name == "" {
		http.Error(w, "Invalid request body or name parameter is missing", http.StatusBadRequest)
		return
	}

	sources := readSourcesFromFile()
	var updatedSources []source.Source
	found := false
	for _, currentSource := range sources {
		if strings.ToLower(string(currentSource.Name)) != strings.ToLower(request.Name) {
			updatedSources = append(updatedSources, currentSource)
		} else {
			found = true
			// Удаление файла
			err := os.Remove(string(currentSource.PathToFile))
			if err != nil {
				http.Error(w, "Failed to delete source file", http.StatusInternalServerError)
				return
			}
		}
	}

	if !found {
		http.Error(w, "Source not found", http.StatusNotFound)
		return
	}

	err = writeSourcesToFile(updatedSources)
	if err != nil {
		http.Error(w, "Failed to write updated sources to file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Source deleted successfully"))
}

func writeSourcesToFile(sources []source.Source) error {
	file, err := os.Create("./storage/sources-storage.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(&sources)
	if err != nil {
		return err
	}

	return nil
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

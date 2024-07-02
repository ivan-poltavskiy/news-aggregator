package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"news-aggregator/entity/article"
	"news-aggregator/entity/source"
	"news-aggregator/parser"
	"os"
	"path/filepath"
)

// addSourceRequest is a structure containing the fields required to add a new source.
type addSourceRequest struct {
	URL string `json:"url"`
}

// AddSourceHandler is a handler for adding the new source to the storage.
func AddSourceHandler(w http.ResponseWriter, r *http.Request) {

	var requestBody addSourceRequest

	if err := getUrlFromRequest(r, &requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rssURL, err := getRssFeedLink(w, requestBody.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filePath, err := downloadRssFeed(rssURL, requestBody.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sourceEntity := source.Source{
		Name:       source.Name(ExtractDomainName(requestBody.URL)),
		PathToFile: source.PathToFile(filePath),
		SourceType: source.STORAGE,
	}

	err, jsonPath := parseAndSaveArticles(sourceEntity, requestBody.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sourceEntity.PathToFile = source.PathToFile(jsonPath)

	if !IsSourceExists(sourceEntity.Name) {
		AddSourceToStorage(sourceEntity)
		if _, err := fmt.Fprintf(w, "RSS feed downloaded and saved to %s and source added", filePath); err != nil {
			log.Printf("Failed to write response: %v", err)
		}
	} else {
		if _, err := fmt.Fprintf(w, "RSS feed downloaded and saved to %s but source already exists", filePath); err != nil {
			log.Printf("Failed to write response: %v", err)
		}
	}

	if _, err := fmt.Fprintf(w, "RSS feed parsed and articles saved successfully"); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

// get the url of the source from the request
func getUrlFromRequest(r *http.Request, requestBody *addSourceRequest) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print("Error closing file: ", err)
		}
	}(r.Body)

	if err := json.Unmarshal(body, requestBody); err != nil || requestBody.URL == "" {
		return fmt.Errorf("invalid request body or URL parameter is missing")
	}

	return nil
}

// get link of rrs feed
func getRssFeedLink(w http.ResponseWriter, url string) (string, error) {
	err, rssURL := GetRssFeedLink(w, url)
	if err != nil || rssURL == "" {
		return "", err
	}
	return rssURL, nil
}

func downloadRssFeed(rssURL, sourceURL string) (string, error) {
	rssResp, err := http.Get(rssURL)
	if err != nil || rssResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download RSS feed")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print("Error closing file: ", err)
		}
	}(rssResp.Body)

	sourceName := ExtractDomainName(sourceURL)
	directoryPath := filepath.Join("resources", sourceName)
	if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory")
	}

	filePath := filepath.Join(directoryPath, sourceName+".xml")
	outputFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file")
	}
	defer func(outputFile *os.File) {
		err := outputFile.Close()
		if err != nil {
			log.Print("Error closing file: ", err)
		}
	}(outputFile)

	if _, err := io.Copy(outputFile, rssResp.Body); err != nil {
		return "", fmt.Errorf("failed to save RSS feed")
	}

	return filePath, nil
}

// parse rss feed and save the articles from this feed to the storage
func parseAndSaveArticles(sourceEntity source.Source, sourceURL string) (error, string) {
	articles, err := parser.Rss{}.Parse(sourceEntity.PathToFile, sourceEntity.Name)
	if err != nil {
		return fmt.Errorf("failed to parse RSS feed"), ""
	}

	jsonFilePath := filepath.Join("resources", ExtractDomainName(sourceURL), ExtractDomainName(sourceURL)+".json")

	var existingArticles []article.Article
	if _, err := os.Stat(jsonFilePath); err == nil {
		jsonFile, err := os.Open(jsonFilePath)
		if err != nil {
			return fmt.Errorf("failed to open existing JSON file"), ""
		}
		defer func(jsonFile *os.File) {
			err := jsonFile.Close()
			if err != nil {
				log.Print("Error closing file: ", err)
			}
		}(jsonFile)

		if err := json.NewDecoder(jsonFile).Decode(&existingArticles); err != nil {
			return fmt.Errorf("failed to decode existing articles from JSON file"), ""
		}
	}

	// Create a map of existing articles for quick lookup
	existingTitles := make(map[string]struct{})
	for _, existingArticle := range existingArticles {
		existingTitles[existingArticle.Title.String()] = struct{}{}
	}

	// Filter out duplicate articles
	var newArticles []article.Article
	for _, newArticle := range articles {
		if _, exists := existingTitles[newArticle.Title.String()]; !exists {
			newArticles = append(newArticles, newArticle)
		}
	}

	// If no new articles to add, skip the file update
	if len(newArticles) == 0 {
		return nil, jsonFilePath
	}

	existingArticles = append(existingArticles, newArticles...)

	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		return fmt.Errorf("failed to create JSON file"), ""
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.Print("Error closing file: ", err)
		}
	}(jsonFile)

	if err := json.NewEncoder(jsonFile).Encode(existingArticles); err != nil {
		return fmt.Errorf("failed to encode articles to JSON file"), ""
	}

	return nil, jsonFilePath
}

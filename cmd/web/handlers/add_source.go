package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
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
		httpResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.Info("AddSourceHandler: The URL from the request to add the source was successfully retrieved: ", requestBody.URL)

	rssURL, err := getRssFeedLink(w, requestBody.URL)
	if err != nil {
		httpResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Info("AddSourceHandler: The URL of feed was successfully retrieved: ", rssURL)

	domainName := ExtractDomainName(requestBody.URL)
	filePath, err := downloadRssFeed(rssURL, domainName)
	if err != nil {
		httpResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sourceEntity := source.Source{
		Name:       source.Name(domainName),
		PathToFile: source.PathToFile(filePath),
		SourceType: source.STORAGE,
	}

	err, jsonPath := parseAndSaveArticles(sourceEntity, domainName)
	if err != nil {
		httpResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sourceEntity.PathToFile = source.PathToFile(jsonPath)

	if !IsSourceExists(sourceEntity.Name) {
		AddSourceToStorage(sourceEntity)
		logrus.Info("RSS feed downloaded and saved to the " + filePath + " and source added")
	} else {
		logrus.Info("RSS feed downloaded and saved to the " + filePath + " but source already exists")
	}

	httpResponse(w, "Articles saved successfully. Name of source: "+string(sourceEntity.Name), http.StatusOK)
}

// get the URL of the source from the request
func getUrlFromRequest(r *http.Request, requestBody *addSourceRequest) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Error("Failed to read request body: ", err)
		return fmt.Errorf("failed to read request body")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Error("Error closing request body: ", err)
		}
	}(r.Body)

	if err := json.Unmarshal(body, requestBody); err != nil || requestBody.URL == "" {
		logrus.Error("Invalid request body or URL parameter is missing")
		return fmt.Errorf("invalid request body or URL parameter is missing")
	}

	logrus.Info("getUrlFromRequest: Successfully parsed request body")
	return nil
}

// get link of RSS feed
func getRssFeedLink(w http.ResponseWriter, url string) (string, error) {
	err, rssURL := GetRssFeedLink(w, url)
	if err != nil || rssURL == "" {
		logrus.Error("Failed to get RSS feed link: ", err)
		return "", err
	}
	logrus.Info("Rss link parsed successfully: ", rssURL)
	return rssURL, nil
}

func downloadRssFeed(rssURL, domainName string) (string, error) {
	rssResp, err := http.Get(rssURL)
	if err != nil || rssResp.StatusCode != http.StatusOK {
		logrus.Error("Failed to download RSS feed: ", err)
		return "", fmt.Errorf("failed to download RSS feed")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Error("Error closing RSS response body: ", err)
		}
	}(rssResp.Body)

	directoryPath := filepath.Join("resources", domainName)
	if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
		logrus.Error("Failed to create directory: ", err)
		return "", fmt.Errorf("failed to create directory")
	}

	filePath := filepath.Join(directoryPath, domainName+".xml")
	outputFile, err := os.Create(filePath)
	if err != nil {
		logrus.Error("Failed to create a file to save the RSS feed to: ", filePath)
		return "", fmt.Errorf("failed to create file")
	}
	defer func(outputFile *os.File) {
		err := outputFile.Close()
		if err != nil {
			logrus.Error("Error closing RSS feed file: ", err)
		}
	}(outputFile)

	if _, err := io.Copy(outputFile, rssResp.Body); err != nil {
		logrus.Error("Could not download RSS feed: ", err)
		return "", fmt.Errorf("failed to save RSS feed")
	}

	logrus.Info("downloadRssFeed: RSS feed successfully downloaded and saved to: ", filePath)
	return filePath, nil
}

// parse RSS feed and save the articles from this feed to the storage
func parseAndSaveArticles(sourceEntity source.Source, domainName string) (error, string) {
	articles, err := parser.Rss{}.Parse(sourceEntity.PathToFile, sourceEntity.Name)
	if err != nil {
		logrus.Error("Failed to parse RSS feed: ", err)
		return fmt.Errorf("failed to parse RSS feed"), ""
	}

	jsonFilePath := filepath.Join("resources", domainName, domainName+".json")

	var existingArticles []article.Article
	if _, err := os.Stat(jsonFilePath); err == nil {
		jsonFile, err := os.Open(jsonFilePath)
		if err != nil {
			logrus.Error("Failed to open existing JSON file: ", err)
			return fmt.Errorf("failed to open existing JSON file"), ""
		}
		defer func(jsonFile *os.File) {
			err := jsonFile.Close()
			if err != nil {
				logrus.Error("Error closing existing JSON file: ", err)
			}
		}(jsonFile)

		if err := json.NewDecoder(jsonFile).Decode(&existingArticles); err != nil {
			logrus.Error("Failed to decode existing articles from JSON file: ", err)
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
		logrus.Info("No new articles to add")
		return nil, jsonFilePath
	}

	existingArticles = append(existingArticles, newArticles...)

	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		logrus.Error("Failed to create JSON file: ", err)
		return fmt.Errorf("failed to create JSON file"), ""
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			logrus.Error("Error closing JSON file: ", err)
		}
	}(jsonFile)

	if err := json.NewEncoder(jsonFile).Encode(existingArticles); err != nil {
		logrus.Error("Failed to encode articles to JSON file: ", err)
		return fmt.Errorf("failed to encode articles to JSON file"), ""
	}

	logrus.Info("parseAndSaveArticles: Articles successfully parsed and saved to: ", jsonFilePath)
	return nil, jsonFilePath
}

// Helper function to write HTTP responses
func httpResponse(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	if _, err := w.Write([]byte(message)); err != nil {
		logrus.Error("Failed to write response: ", err)
	}
}

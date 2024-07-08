package service

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"news-aggregator/constant"
	"news-aggregator/entity/article"
	"news-aggregator/entity/source"
	"news-aggregator/parser"
	"os"
	"path/filepath"
	"regexp"
)

// AddSource processes the source URL and returns the source entity
func AddSource(url string) (source.Name, error) {
	if url == "" {
		return "", fmt.Errorf("passed url is empty")
	}
	rssURL, err := getRssFeedLink(url)
	if err != nil {
		return "", err
	}
	logrus.Info("AddSource: The URL of feed was successfully retrieved: ", rssURL)

	domainName := ExtractDomainName(url)

	filePath, err := downloadRssFeed(rssURL, domainName)
	if err != nil {
		return "", err
	}

	sourceEntity := source.Source{
		Name:       source.Name(domainName),
		PathToFile: source.PathToFile(filePath),
		SourceType: source.STORAGE,
	}

	err, jsonPath := parseAndSaveArticles(sourceEntity, domainName)
	if err != nil {
		return "", err
	}
	sourceEntity.PathToFile = source.PathToFile(jsonPath)

	if !IsSourceExists(sourceEntity.Name) {
		AddSourceToStorage(sourceEntity)
		logrus.Info("Source added")
	} else {
		logrus.Info("Source already exists")
	}
	return sourceEntity.Name, nil
}

func getRssFeedLink(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		logrus.Error("getRssFeedLink: RSS URL not found ", err)
		return "", fmt.Errorf("rss url not found: %s", url)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Error("getRssFeedLink: Error closing response body ", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("getRssFeedLink: Failed to read page content ", err)
		return "", err
	}

	re := regexp.MustCompile(`(?i)<link[^>]+type="application/rss\+xml"[^>]+href="([^"]+)"`)
	matches := re.FindStringSubmatch(string(body))

	if len(matches) < 2 {
		logrus.Warn("getRssFeedLink: RSS link not found")
		return "", nil
	}

	rssURL := matches[1]
	logrus.Info("getRssFeedLink: RSS link found: ", rssURL)
	return rssURL, nil
}

// downloadRssFeed downloads the RSS feed and returns the file path
func downloadRssFeed(rssURL, domainName string) (string, error) {
	rssResp, err := http.Get(rssURL)
	if err != nil || rssResp.StatusCode != http.StatusOK {
		logrus.Error("Failed to download RSS feed: ", err)
		return "", fmt.Errorf("failed to download RSS feed")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Error("Error closing response body ", err)
			return
		}
	}(rssResp.Body)

	directoryPath := filepath.Join(constant.PathToResources, domainName)
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
			logrus.Error("Failed to close a file to save the RSS feed to: ", filePath)
			return
		}
	}(outputFile)

	if _, err := io.Copy(outputFile, rssResp.Body); err != nil {
		logrus.Error("Could not download RSS feed: ", err)
		return "", fmt.Errorf("failed to save RSS feed")
	}

	logrus.Info("downloadRssFeed: RSS feed successfully downloaded and saved to: ", filePath)
	return filePath, nil
}

// parseAndSaveArticles parses RSS feed and saves the articles to the storage
func parseAndSaveArticles(sourceEntity source.Source, domainName string) (error, string) {
	articles, err := parseRssFeed(sourceEntity)
	if err != nil {
		return err, ""
	}

	jsonFilePath := filepath.ToSlash(filepath.Join(constant.PathToResources, domainName, domainName+".json"))

	existingArticles, err := readExistingArticles(jsonFilePath)
	if err != nil {
		return err, ""
	}

	newArticles := filterNewArticles(articles, existingArticles)
	if len(newArticles) == 0 {
		logrus.Info("No new articles to add")
		return nil, jsonFilePath
	}

	existingArticles = append(existingArticles, newArticles...)

	if err := saveArticles(jsonFilePath, existingArticles); err != nil {
		return err, ""
	}

	logrus.Info("parseAndSaveArticles: Articles successfully parsed and saved to: ", jsonFilePath)
	return nil, jsonFilePath
}

func parseRssFeed(sourceEntity source.Source) ([]article.Article, error) {
	articles, err := parser.Rss{}.Parse(sourceEntity.PathToFile, sourceEntity.Name)
	if err != nil {
		logrus.Error("Failed to parse RSS feed: ", err)
		return nil, fmt.Errorf("failed to parse RSS feed")
	}
	return articles, nil
}

func readExistingArticles(jsonFilePath string) ([]article.Article, error) {
	var existingArticles []article.Article

	if _, err := os.Stat(jsonFilePath); err == nil {
		jsonFile, err := os.Open(jsonFilePath)
		if err != nil {
			logrus.Error("Failed to open existing JSON file: ", err)
			return nil, fmt.Errorf("failed to open existing JSON file")
		}
		defer func(jsonFile *os.File) {
			err := jsonFile.Close()
			if err != nil {
				logrus.Error("Failed to close the existing JSON file: ", err)
			}
		}(jsonFile)

		if err := json.NewDecoder(jsonFile).Decode(&existingArticles); err != nil {
			logrus.Error("Failed to decode existing articles from JSON file: ", err)
			return nil, fmt.Errorf("failed to decode existing articles from JSON file")
		}
	}

	return existingArticles, nil
}
func filterNewArticles(articles []article.Article, existingArticles []article.Article) []article.Article {
	existingTitles := make(map[string]struct{})
	for _, existingArticle := range existingArticles {
		existingTitles[existingArticle.Title.String()] = struct{}{}
	}

	var newArticles []article.Article
	for _, newArticle := range articles {
		if _, exists := existingTitles[newArticle.Title.String()]; !exists {
			newArticles = append(newArticles, newArticle)
		}
	}

	return newArticles
}

func saveArticles(jsonFilePath string, articles []article.Article) error {
	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		logrus.Error("Failed to create JSON file: ", err)
		return fmt.Errorf("failed to create JSON file")
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			logrus.Error("Failed to close the JSON file: ", err)
		}
	}(jsonFile)

	if err := json.NewEncoder(jsonFile).Encode(articles); err != nil {
		logrus.Error("Failed to encode articles to JSON file: ", err)
		return fmt.Errorf("failed to encode articles to JSON file")
	}

	return nil
}

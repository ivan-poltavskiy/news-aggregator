package service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"news-aggregator/constant"
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	"news-aggregator/parser"
	newsStorage "news-aggregator/storage/news"
	sourceStorage "news-aggregator/storage/source"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

// SaveSource processes the source URL and returns the source entity
func SaveSource(url string, sourceStorage sourceStorage.Storage, newsStorage newsStorage.NewsStorage) (source.Name, error) {

	if url == "" {
		return "", fmt.Errorf("passed url is empty")
	}

	rssURL, err := getRssFeedLink(url)
	if err != nil {
		return "", err
	}
	logrus.Info("Save: The URL of feed was successfully retrieved: ", rssURL)

	domainName := ExtractDomainName(url)

	filePath, err := downloadRssFeed(rssURL, domainName)
	if err != nil {
		return "", err
	}

	sourceEntity := source.Source{
		Name:       source.Name(domainName),
		PathToFile: source.PathToFile(filePath),
		SourceType: source.STORAGE,
		Link:       source.Link(url),
	}

	err, jsonPath := parseAndSaveNews(sourceEntity, newsStorage)
	if err != nil {
		return "", err
	}
	sourceEntity.PathToFile = source.PathToFile(jsonPath)

	if !IsSourceExists(sourceEntity.Name, sourceStorage) {
		err = sourceStorage.SaveSource(sourceEntity)
		if err != nil {
			return "", err
		}
		logrus.Info("Source added")
	} else {
		logrus.Info("Source already exists")
	}
	return sourceEntity.Name, nil
}

// PeriodicallyUpdateNews updates news for all sources.
func PeriodicallyUpdateNews(sourceStorage sourceStorage.Storage, newsUpdatePeriod time.Duration, newsStorage newsStorage.NewsStorage) {
	ticker := time.NewTicker(newsUpdatePeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logrus.Info("Starting periodic update of news")
			sources, err := sourceStorage.GetSources()
			if err != nil {
				logrus.Error("Failed to retrieve sources: ", err)
				continue
			}

			var wg sync.WaitGroup
			for _, src := range sources {
				wg.Add(1)
				go func(src source.Source) {
					defer wg.Done()
					err := updateSourceNews(src, newsStorage)
					if err != nil {
						logrus.Error("Failed to update news for source: ", src.Name, err)
					}
				}(src)
			}
			wg.Wait()
			logrus.Info("Periodic update of news completed")
		}
	}
}

// getRssFeedLink takes link of rss feed from the input site
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
	rssResponse, err := http.Get(rssURL)
	if err != nil || rssResponse.StatusCode != http.StatusOK {
		logrus.Error("Failed to download RSS feed: ", err)
		return "", fmt.Errorf("failed to download RSS feed")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Error("Error closing response body ", err)
			return
		}
	}(rssResponse.Body)

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

	if _, err := io.Copy(outputFile, rssResponse.Body); err != nil {
		logrus.Error("Could not download RSS feed: ", err)
		return "", fmt.Errorf("failed to save RSS feed")
	}

	logrus.Info("downloadRssFeed: RSS feed successfully downloaded and saved to: ", filePath)
	return filePath, nil
}

// parseAndSaveNews parses RSS feed and saves the news to the storage
func parseAndSaveNews(sourceEntity source.Source, newsStorage newsStorage.NewsStorage) (error, string) {
	parsedNews, err := parseRssFeed(sourceEntity)
	if err != nil {
		return err, ""
	}

	jsonFilePath := filepath.ToSlash(filepath.Join(constant.PathToResources, string(sourceEntity.Name), string(sourceEntity.Name)+".json"))

	existingNews, err := newsStorage.GetNews(jsonFilePath)
	if err != nil {
		return err, ""
	}

	newArticles := newsUnification(parsedNews, existingNews)
	if len(newArticles) == 0 {
		logrus.Info("No new parsedNews to add")
		return nil, jsonFilePath
	}

	existingNews = append(existingNews, newArticles...)

	if err := newsStorage.SaveNews(jsonFilePath, existingNews); err != nil {
		return err, ""
	}

	logrus.Info("parseAndSaveNews: Articles successfully parsed and saved to: ", jsonFilePath)
	return nil, jsonFilePath
}

// parseRssFeed parses rss feed from the input site and return the news from it
func parseRssFeed(sourceEntity source.Source) ([]news.News, error) {
	parsedNews, err := parser.Rss{}.Parse(sourceEntity.PathToFile, sourceEntity.Name)
	if err != nil {
		logrus.Error("Failed to parse RSS feed: ", err)
		return nil, fmt.Errorf("failed to parse RSS feed")
	}
	return parsedNews, nil
}

// newsUnification checks whether there are articles from the new feed in the existing news, and if so, removes them
func newsUnification(articles []news.News, existingArticles []news.News) []news.News {
	existingTitles := make(map[string]struct{})
	for _, existingArticle := range existingArticles {
		existingTitles[existingArticle.Title.String()] = struct{}{}
	}

	var newArticles []news.News
	for _, newArticle := range articles {
		if _, exists := existingTitles[newArticle.Title.String()]; !exists {
			newArticles = append(newArticles, newArticle)
		}
	}

	return newArticles
}

// updateSourceNews updating the news of the input source
func updateSourceNews(inputSource source.Source, newsStorage newsStorage.NewsStorage) error {
	domainName := ExtractDomainName(string(inputSource.Link))
	rssURL, err := getRssFeedLink(string(inputSource.Link))
	if err != nil {
		return err
	}

	filePath, err := downloadRssFeed(rssURL, domainName)
	if err != nil {
		return err
	}

	inputSource.PathToFile = source.PathToFile(filePath)

	err, jsonPath := parseAndSaveNews(inputSource, newsStorage)
	if err != nil {
		return err
	}
	inputSource.PathToFile = source.PathToFile(jsonPath)

	return nil
}

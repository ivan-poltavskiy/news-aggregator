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

	err, jsonPath := parseAndSaveNews(sourceEntity, domainName, newsStorage)
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

// parseAndSaveNews parses RSS feed and saves the news to the storage
func parseAndSaveNews(sourceEntity source.Source, domainName string, newsStorage newsStorage.NewsStorage) (error, string) {
	articles, err := parseRssFeed(sourceEntity)
	if err != nil {
		return err, ""
	}

	jsonFilePath := filepath.ToSlash(filepath.Join(constant.PathToResources, domainName, domainName+".json"))

	existingArticles, err := newsStorage.GetNews(jsonFilePath)
	if err != nil {
		return err, ""
	}

	newArticles := newsUnification(articles, existingArticles)
	if len(newArticles) == 0 {
		logrus.Info("No new articles to add")
		return nil, jsonFilePath
	}

	existingArticles = append(existingArticles, newArticles...)

	if err := newsStorage.SaveNews(jsonFilePath, existingArticles); err != nil {
		return err, ""
	}

	logrus.Info("parseAndSaveNews: Articles successfully parsed and saved to: ", jsonFilePath)
	return nil, jsonFilePath
}

func parseRssFeed(sourceEntity source.Source) ([]news.News, error) {
	articles, err := parser.Rss{}.Parse(sourceEntity.PathToFile, sourceEntity.Name)
	if err != nil {
		logrus.Error("Failed to parse RSS feed: ", err)
		return nil, fmt.Errorf("failed to parse RSS feed")
	}
	return articles, nil
}

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

func updateSourceNews(src source.Source, newsStorage newsStorage.NewsStorage) error {
	domainName := ExtractDomainName(string(src.Link))
	rssURL, err := getRssFeedLink(string(src.Link))
	if err != nil {
		return err
	}

	filePath, err := downloadRssFeed(rssURL, domainName)
	if err != nil {
		return err
	}

	src.PathToFile = source.PathToFile(filePath)

	err, jsonPath := parseAndSaveNews(src, domainName, newsStorage)
	if err != nil {
		return err
	}
	src.PathToFile = source.PathToFile(jsonPath)

	return nil
}

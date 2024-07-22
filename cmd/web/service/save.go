package service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	"news-aggregator/parser"
	newsStorage "news-aggregator/storage/news"
	sourceStorage "news-aggregator/storage/source"
	"os"
	"regexp"
	"strings"
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

	parsedNews, err := parseRssFeed(rssURL, domainName)
	if err != nil {
		return "", err
	}

	sourceEntity := source.Source{
		Name:       source.Name(domainName),
		SourceType: source.STORAGE,
		Link:       source.Link(url),
	}

	sourceEntity, err = saveNews(sourceEntity, newsStorage, sourceStorage, parsedNews)
	if err != nil {
		return "", err
	}

	if !sourceStorage.IsSourceExists(sourceEntity.Name) {
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
					err := updateSourceNews(src, newsStorage, sourceStorage)
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

// parseRssFeed downloads the RSS feed and returns the parsed news
func parseRssFeed(rssURL, domainName string) ([]news.News, error) {
	rssResponse, err := http.Get(rssURL)
	if err != nil || rssResponse.StatusCode != http.StatusOK {
		logrus.Error("Failed to download RSS feed: ", err)
		return nil, fmt.Errorf("failed to download RSS feed")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Error(err)
			return
		}
	}(rssResponse.Body)

	tempFile, err := os.CreateTemp("", "*.xml")
	if err != nil {
		logrus.Error("Failed to create temporary file: ", err)
		return nil, fmt.Errorf("failed to create temporary file")
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			logrus.Error(err)
			return
		}
	}(tempFile.Name())

	if _, err := io.Copy(tempFile, rssResponse.Body); err != nil {
		logrus.Error("Failed to save RSS feed to temporary file: ", err)
		return nil, fmt.Errorf("failed to save RSS feed")
	}

	parsedNews, err := parser.Rss{}.Parse(source.PathToFile(tempFile.Name()), source.Name(domainName))
	if err != nil {
		logrus.Error("Failed to parse RSS feed: ", err)
		return nil, fmt.Errorf("failed to parse RSS feed")
	}

	return parsedNews, nil
}

// saveNews saves the news to the storage
func saveNews(sourceEntity source.Source, newsStorage newsStorage.NewsStorage, sourceStorage sourceStorage.Storage, parsedNews []news.News) (source.Source, error) {

	existingNews, err := newsStorage.GetNewsBySourceName(sourceEntity.Name, sourceStorage)
	if err != nil {
		return source.Source{}, err
	}

	newArticles := newsUnification(parsedNews, existingNews)
	if len(newArticles) == 0 {
		logrus.Info("No new parsed news to add")
		return sourceEntity, nil
	}

	existingNews = append(existingNews, newArticles...)

	sourceEntity, err = newsStorage.SaveNews(sourceEntity, existingNews)

	return sourceEntity, nil
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
func updateSourceNews(inputSource source.Source, newsStorage newsStorage.NewsStorage, sourceStorage sourceStorage.Storage) error {
	domainName := ExtractDomainName(string(inputSource.Link))
	rssURL, err := getRssFeedLink(string(inputSource.Link))
	if err != nil {
		return err
	}

	currentNews, err := parseRssFeed(rssURL, domainName)
	if err != nil {
		return err
	}

	_, err = saveNews(inputSource, newsStorage, sourceStorage, currentNews)
	if err != nil {
		return err
	}

	return nil
}

// ExtractDomainName parse the url to get the resource domain
func ExtractDomainName(url string) string {
	re := regexp.MustCompile(`https?://(www\.)?([^/]+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 3 {
		logrus.Warn("ExtractDomainName: Failed to extract domain name from URL: ", url)
		return "unknown"
	}
	domain := matches[2]
	domain = strings.Split(domain, ".")[0]
	logrus.Info("ExtractDomainName: Extracted domain name: ", domain)
	return domain
}

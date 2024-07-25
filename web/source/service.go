package source

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"news-aggregator/entity/source"
	"news-aggregator/storage"
	"news-aggregator/web/feed"
	"news-aggregator/web/news"
)

type Service struct {
	storage storage.Storage
}

// NewService creates new instance of the Service
func NewService(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

// DeleteSourceByName removes the source from storage by name.
func (service *Service) DeleteSourceByName(name source.Name) error {
	err := service.storage.DeleteSourceByName(name)
	if err != nil {
		logrus.Error("Error deleting source:", err)
		return err
	}
	return nil
}

// SaveSource processes the source URL and returns the source entity
func (service *Service) SaveSource(url string) (source.Name, error) {

	if url == "" {
		return "", fmt.Errorf("passed url is empty")
	}

	rssURL, err := feed.GetRssFeedLink(url)
	if err != nil {
		return "", err
	}
	logrus.Info("Save: The URL of feed was successfully retrieved: ", rssURL)

	domainName := feed.ExtractDomainName(url)

	parsedNews, err := feed.ParseRssFeed(rssURL, domainName)
	if err != nil {
		return "", err
	}

	sourceEntity := source.Source{
		Name:       source.Name(domainName),
		SourceType: source.STORAGE,
		Link:       source.Link(url),
	}
	newsService := news.NewService(service.storage)
	sourceEntity, err = newsService.SaveNews(sourceEntity, parsedNews)
	if err != nil {
		return "", err
	}

	if !service.storage.IsSourceExists(sourceEntity.Name) {
		err = service.storage.SaveSource(sourceEntity)
		if err != nil {
			return "", err
		}
		logrus.Info("Source added")
	} else {
		logrus.Info("Source already exists")
	}
	return sourceEntity.Name, nil
}

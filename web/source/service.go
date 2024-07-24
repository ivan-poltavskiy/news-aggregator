package source

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"news-aggregator/entity/source"
	"news-aggregator/storage"
	"news-aggregator/web/feed"
	"news-aggregator/web/news"
)

type SourcesService struct {
	storage storage.Storage
}

// NewSourceService creates new instance of the SourcesService
func NewSourceService(storage storage.Storage) *SourcesService {
	return &SourcesService{
		storage: storage,
	}
}

// DeleteSourceByName removes the source from storage by name.
func (service *SourcesService) DeleteSourceByName(name string) error {
	err := service.storage.DeleteSourceByName(name)
	if err != nil {
		logrus.Error("Error deleting source:", err)
		return err
	}
	return nil
}

// SaveSource processes the source URL and returns the source entity
func (service *SourcesService) SaveSource(url string) (source.Name, error) {

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
	newsService := news.NewNewsService(service.storage)
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

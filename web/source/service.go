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
func (service *Service) SaveSource(request AddSourceRequest) (source.Name, error) {

	if request.URL == "" || request.Name == "" {
		return "", fmt.Errorf("passed url or name are empty")
	}

	rssURL, err := feed.GetRssFeedLink(request.URL)
	if err != nil {
		return "", err
	}
	logrus.Info("Save: The URL of feed was successfully retrieved: ", rssURL)

	parsedNews, err := feed.ParseRssFeed(rssURL, request.Name)
	if err != nil {
		return "", err
	}

	sourceEntity := source.Source{
		Name:       source.Name(request.Name),
		SourceType: source.STORAGE,
		Link:       source.Link(request.URL),
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

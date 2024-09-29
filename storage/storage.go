package storage

import (
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
)

// Storage is the abstract repository for saving some resources in the app
//
//go:generate mockgen -source=storage.go -destination=mock_aggregator/mock_storage.go -package=client news-aggregator/storage Storage
type Storage interface {
	News
	Source
}

type jsonStorage struct {
	News
	Source
}

// NewStorage returns the new instance of the Storage interface
func NewStorage(newsStorage News, sourceStorage Source) Storage {
	return &jsonStorage{
		News:   newsStorage,
		Source: sourceStorage,
	}
}

// News is an implementation of the Storage for managing the news
type News interface {
	// SaveNews saves the news of provided source and return the entity of this source
	SaveNews(providedSource source.Source, news []news.News) (source.Source, error)
	// GetNews returns the slice of news from the provided path
	GetNews(path string) ([]news.News, error)
	// GetNewsBySourceName returns the slice of news by source name from the provided source storage
	GetNewsBySourceName(sourceName source.Name, sourceStorage Source) ([]news.News, error)
}

// Source is an implementation of the Storage for managing the sources
type Source interface {
	// SaveSource saves the provided source to the storage
	SaveSource(source source.Source) error
	// DeleteSourceByName removes the source by provided source's name
	DeleteSourceByName(source.Name) error
	// GetSources returns the all sources from the storage
	GetSources() ([]source.Source, error)
	// IsSourceExists check
	IsSourceExists(source.Name) bool
	GetSourceByName(source.Name) (source.Source, error)
	UpdateSource(updatedSource source.Source, currentName string) error
}

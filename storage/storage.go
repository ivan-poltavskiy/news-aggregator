package storage

import (
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
)

// Storage is the abstract repository for saving some resources in the app
//
//go:generate mockgen -source=storage.go -destination=mock_aggregator/mock_storage.go -package=mock_aggregator news-aggregator/storage Storage
type Storage interface {
	NewsStorage
	SourceStorage
}

type NewsStorage interface {
	SaveNews(currentSource source.Source, news []news.News) (source.Source, error)
	GetNews(path string) ([]news.News, error)
	GetNewsBySourceName(sourceName source.Name, sourceStorage SourceStorage) ([]news.News, error)
}

type SourceStorage interface {
	SaveSource(source source.Source) error
	DeleteSourceByName(name string) error
	GetSources() ([]source.Source, error)
	IsSourceExists(source.Name) bool
	GetSourceByName(source.Name) (source.Source, error)
}

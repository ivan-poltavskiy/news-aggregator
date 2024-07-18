package storage

import "news-aggregator/entity/source"

//go:generate mockgen -source=C:\Users\polta\GolandProjects\news-aggregator\storage\storage.go -destination=C:\Users\polta\GolandProjects\news-aggregator\mock_aggregator\mock_storage.go -package=mock_aggregator news-aggregator/storage Storage

// Storage is the type of repository for saving sources
type Storage interface {
	SaveSource(source source.Source) error
	DeleteSourceByName(name string) error
	GetSources() ([]source.Source, error)
}

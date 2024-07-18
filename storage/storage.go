package storage

import "news-aggregator/entity/source"

//go:generate mockgen -source=C:\Users\dange\GolandProjects\news-aggregator\storage\storage.go -destination=mock_aggregator\mock_storage.go -package=mock_aggregator news-aggregator/cmd Storage

// Storage is the type of repository for saving sources
type Storage interface {
	SaveSource(source source.Source) error
	DeleteSource()
	GetSources() ([]source.Source, error)
}

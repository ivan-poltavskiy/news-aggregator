package storage

import "news-aggregator/entity/source"

// Storage is the type of repository for saving sources
type Storage interface {
	SaveSource()
	DeleteSource()
	GetSources() ([]source.Source, error)
}

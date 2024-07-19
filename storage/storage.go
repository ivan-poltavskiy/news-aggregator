package storage

import "news-aggregator/entity/source"

//go:generate mockgen -source=C:\Users\polta\GolandProjects\news-aggregator\storage\storage.go -destination=C:\Users\polta\GolandProjects\news-aggregator\mock_aggregator\mock_storage.go -package=mock_aggregator news-aggregator/storage Storage

// Storage is the type of repository for saving and managing sources.
type Storage interface {
	// SaveSource saves the provided source to the storage.
	//Returns the error, if the save fails.
	SaveSource(source source.Source) error
	// DeleteSourceByName removes the source by his name from the storage.
	// Returns the error, if the deleting fails.
	DeleteSourceByName(name string) error
	// GetSources returns the slice of the sources which are provided in the storage.
	// Returns the empty slice an error if the getting process fails.
	GetSources() ([]source.Source, error)
}

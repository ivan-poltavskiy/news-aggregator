package news

import (
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	sourceStorage "news-aggregator/storage/source"
)

//go:generate mockgen -source=storage.go -destination=mock_aggregator\mock_news_storage.go -package=mock_aggregator news-aggregator/storage/news NewsStorage
type NewsStorage interface {
	// SaveNews saves the provided news to the storage.
	//Returns the error, if the save fails.
	SaveNews(currentSource source.Source, news []news.News) (source.Source, error)
	// GetNews returns the slice of the news which are provided in the storage.
	// Returns the empty slice an error if the getting process fails.
	GetNews(path string) ([]news.News, error)
	// GetNewsBySourceName returns the slice of news of provided source
	GetNewsBySourceName(sourceName source.Name, sourceStorage sourceStorage.Storage) ([]news.News, error)
}

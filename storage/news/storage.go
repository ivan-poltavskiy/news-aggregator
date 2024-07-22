package news

import (
	"news-aggregator/entity/news"
)

//go:generate mockgen -source=storage.go -destination=mock_aggregator\mock_news_storage.go -package=mock_aggregator news-aggregator/storage/news NewsStorage
type NewsStorage interface {
	// SaveNews saves the provided news to the storage.
	//Returns the error, if the save fails.
	SaveNews(jsonFilePath string, news []news.News) error
	// GetNews returns the slice of the news which are provided in the storage.
	// Returns the empty slice an error if the getting process fails.
	GetNews(jsonFilePath string) ([]news.News, error)
}

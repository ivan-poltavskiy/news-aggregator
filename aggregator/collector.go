package aggregator

import (
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
)

// Collector is using for fetching news from source by source name.
//
//go:generate mockgen -source=C:\Users\dange\GolandProjects\news-aggregator\aggregator\collector.go -destination=mock_aggregator\mock_collector.go -package=mock_aggregator news-aggregator/aggregator Collector
type Collector interface {
	FindNewsByResourcesName(sourcesNames []source.Name) ([]news.News, error)
}

package aggregator

import (
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
)

// Collector is using for fetching news from source by source name.
//
//go:generate mockgen -source=collector.go -destination=mock/mock_collector.go -package=mock news-aggregator/aggregator Collector
type Collector interface {
	FindNewsByResourcesName(sourcesNames []source.Name) ([]news.News, error)
}

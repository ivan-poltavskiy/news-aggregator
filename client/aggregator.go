package client

import (
	"news-aggregator/entity/news"
	"news-aggregator/filter"
)

// Aggregator defines an interface for aggregating collector news.
//
//go:generate mockgen -source=aggregator.go -destination=mock_aggregator/mock_aggregator.go -package=client  news-aggregator/aggregator Aggregator
type Aggregator interface {
	// Aggregate fetches news from the provided sources,
	//applies the given filters, and returns the filtered news.
	Aggregate(sources []string, filters ...filter.NewsFilter) ([]news.News, error)
}

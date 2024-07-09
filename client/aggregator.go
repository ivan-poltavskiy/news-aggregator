package client

import (
	"news-aggregator/entity/news"
	"news-aggregator/filter"
)

// Aggregator defines an interface for aggregating collector news.
type Aggregator interface {
	// Aggregate fetches news from the provided sources,
	//applies the given filters, and returns the filtered news.
	Aggregate(sources []string, filters ...filter.NewsFilter) ([]news.News, error)
}

package client

import (
	"NewsAggregator/entity/article"
	"NewsAggregator/service/filter"
)

// Aggregator defines an interface for aggregating news articles.
type Aggregator interface {
	// Aggregate fetches articles from the provided sources,
	//applies the given filters, and returns the filtered articles.
	Aggregate(sources []string, filters ...filter.Service) ([]article.Article, string)
}

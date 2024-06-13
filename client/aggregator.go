package client

import (
	"news_aggregator/entity/article"
	"news_aggregator/filter"
)

// Aggregator defines an interface for aggregating collector articles.
type Aggregator interface {
	// Aggregate fetches articles from the provided sources,
	//applies the given filters, and returns the filtered articles.
	Aggregate(sources []string, filters ...filter.ArticleFilter) ([]article.Article, error)
}

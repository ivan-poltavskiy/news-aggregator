package client

import (
	"news-aggregator/entity/article"
	"news-aggregator/filter"
)

// Aggregator defines an interface for aggregating collector articles.
type Aggregator interface {
	// Aggregate fetches articles from the provided sources,
	//applies the given filters, and returns the filtered articles.
	Aggregate(sources []string, filters ...filter.ArticleFilter) ([]article.Article, error)
}

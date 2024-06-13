package aggregator

import (
	"fmt"
	"news_aggregator/client"
	"news_aggregator/collector"
	"news_aggregator/entity/article"
	"news_aggregator/entity/source"
	"news_aggregator/filter"
	"news_aggregator/validator"
)

// news provides methods for aggregating collector articles from various sources.
type news struct{}

func New() client.Aggregator {
	return &news{}
}

// Aggregate fetches articles from the provided sources, applies the given
// filters, and returns the filtered articles.
// Parameters:
// - sources: a slice of strings representing the names of the sources to fetch articles from.
// - filters: a variadic parameter of filter.Service to apply filters to the fetched articles.
//
// Returns:
// - A slice of articles that have been fetched and filtered.
// - An error message string if any errors occurred during the process.
//
//go:generate mockgen -destination=mock_aggregator/mock_aggregator.go -package=mock_aggregator news-aggregator/client Aggregator
func (aggregator *news) Aggregate(sources []string, filters ...filter.ArticleFilter) ([]article.Article, string) {
	var sourceNameObjects []source.Name

	for _, name := range sources {
		sourceNameObjects = append(sourceNameObjects, source.Name(name))
	}

	articles, errorMessage := collector.FindByResourcesName(sourceNameObjects)
	if errorMessage != "" {
		return nil, errorMessage
	}

	fmt.Println(validator.CheckSource(articles))

	for _, filter := range filters {
		articles = filter.Filter(articles)
	}

	return articles, ""
}

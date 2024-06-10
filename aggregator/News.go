package aggregator

import (
	"NewsAggregator/client"
	"NewsAggregator/collector"
	"NewsAggregator/entity/article"
	"NewsAggregator/entity/source"
	"NewsAggregator/filter"
	"fmt"
)

// News provides methods for aggregating collector articles from various sources.
type News struct{}

func New() *News {
	return &News{}
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
//go:generate mockgen -destination=mock_aggregator/mock_aggregator.go -package=mock_aggregator NewsAggregator/client Aggregator
func (na *News) Aggregate(sources []string, filters ...filter.ArticleFilter) ([]article.Article, string) {
	var sourceNameObjects []source.Name

	for _, name := range sources {
		sourceNameObjects = append(sourceNameObjects, source.Name(name))
	}

	articles, errorMessage := collector.FindByResourcesName(sourceNameObjects)
	if errorMessage != "" {
		return nil, errorMessage
	}

	fmt.Println(client.CheckSource(articles))

	for _, filter := range filters {
		articles = filter.Filter(articles)
	}

	return articles, ""
}

package aggregator

import (
	"news_aggregator/client"
	"news_aggregator/collector"
	"news_aggregator/entity/article"
	"news_aggregator/entity/source"
	"news_aggregator/filter"
	"news_aggregator/validator"
)

// news provides methods for aggregating collector articles from various sources.
type news struct {
	Collector collector.ArticleCollector
}

func New(articleCollector *collector.ArticleCollector) client.Aggregator {
	news := &news{Collector: *articleCollector}
	return news
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
func (aggregator *news) Aggregate(sources []string, filters ...filter.ArticleFilter) ([]article.Article, error) {
	var sourceNames []source.Name

	for _, name := range sources {
		sourceNames = append(sourceNames, source.Name(name))
	}

	validateSource, err := validator.ValidateSource(sources)
	if !validateSource {
		return nil, err
	}

	articles, err := aggregator.Collector.FindNewsByResourcesName(sourceNames)
	if err != nil {
		return nil, err
	}

	for _, f := range filters {
		articles = f.Filter(articles)
	}

	return articles, nil
}

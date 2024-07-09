package aggregator

import (
	"news-aggregator/client"
	"news-aggregator/entity/article"
	"news-aggregator/entity/source"
	"news-aggregator/filter"
	"news-aggregator/validator"
)

// news provides methods for aggregating newsCollector articles from various sources.
type news struct {
	newsCollector Collector
}

func New(articleCollector Collector) client.Aggregator {
	news := &news{newsCollector: articleCollector}
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

	articles, err := aggregator.newsCollector.FindNewsByResourcesName(sourceNames)
	if err != nil {
		return nil, err
	}

	for _, f := range filters {
		articles = f.Filter(articles)
	}

	return articles, nil
}

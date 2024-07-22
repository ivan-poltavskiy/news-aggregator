package aggregator

import (
	"news-aggregator/client"
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	"news-aggregator/filter"
	"news-aggregator/validator"
)

// newsAggregator provides methods for aggregating news from various sources.
type newsAggregator struct {
	newsCollector Collector
}

func New(newsCollector Collector) client.Aggregator {
	newsAggregator := &newsAggregator{newsCollector: newsCollector}
	return newsAggregator
}

// Aggregate fetches news from the provided sources, applies the given
// filters, and returns the filtered news.
// Parameters:
// - sources: a slice of strings representing the names of the sources to fetch news from.
// - filters: a variadic parameter of filter.Service to apply filters to the fetched news.
//
// Returns:
// - A slice of news that have been fetched and filtered.
// - An error message string if any errors occurred during the process.

func (aggregator *newsAggregator) Aggregate(sources []string, filters ...filter.NewsFilter) ([]news.News, error) {
	var sourceNames []source.Name

	for _, name := range sources {
		sourceNames = append(sourceNames, source.Name(name))
	}

	validateSource, err := validator.ValidateSource(sources)
	if !validateSource {
		return nil, err
	}

	news, err := aggregator.newsCollector.FindNewsByResourcesName(sourceNames)
	if err != nil {
		return nil, err
	}

	for _, f := range filters {
		news = f.Filter(news)
	}

	return news, nil
}

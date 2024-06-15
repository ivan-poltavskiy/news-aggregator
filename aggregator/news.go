package aggregator

import (
	"errors"
	"fmt"
	"news_aggregator/client"
	"news_aggregator/collector"
	"news_aggregator/constants"
	"news_aggregator/entity/article"
	"news_aggregator/entity/source"
	"news_aggregator/filter"
	"news_aggregator/validator"
	"strings"
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
//go:generate mockgen -destination=mock_aggregator/mock_aggregator.go -package=mock_aggregator news_aggregator/client Aggregator
func (aggregator *news) Aggregate(sources []string, filters ...filter.ArticleFilter) ([]article.Article, error) {
	var sourceNames []source.Name

	for _, name := range sources {
		sourceNames = append(sourceNames, source.Name(name))
	}

	articles, err := collector.FindByResourcesName(sourceNames)
	if err != nil {
		return nil, err
	}

	if !validator.ValidateSource(articles) {
		return nil, errors.New(fmt.Sprintf("Please, specify at least one "+
			"news source. The program supports such news resources:\n%s.",
			strings.Join(constants.NewsSources, ", ")))
	}

	for _, f := range filters {
		articles = f.Filter(articles)
	}

	return articles, nil
}

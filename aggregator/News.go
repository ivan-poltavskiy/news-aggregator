package aggregator

import (
	"NewsAggregator/entity/article"
	"NewsAggregator/entity/source"
	"NewsAggregator/service/filter"
	"NewsAggregator/service/news"
)

// News provides methods for aggregating news articles from various sources.
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
func (na *News) Aggregate(sources []string, filters ...filter.Service) ([]article.Article, string) {
	sourceNames := filterUnique(sources)
	var sourceNameObjects []source.Name

	for _, name := range sourceNames {
		sourceNameObjects = append(sourceNameObjects, source.Name(name))
	}

	articles, errorMessage := news.FindByResourcesName(sourceNameObjects)
	if errorMessage != "" {
		return nil, errorMessage
	}

	if len(articles) == 0 {
		return nil, "Sources not found."
	}

	for _, filter := range filters {
		articles = filter.Filter(articles)
	}

	return articles, ""
}

// filterUnique returns a slice containing only unique strings from the input slice.
func filterUnique(input []string) []string {
	uniqueMap := make(map[string]struct{})
	var uniqueList []string
	for _, item := range input {
		if _, ok := uniqueMap[item]; !ok {
			uniqueMap[item] = struct{}{}
			uniqueList = append(uniqueList, item)
		}
	}
	return uniqueList
}

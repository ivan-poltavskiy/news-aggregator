package aggregator

import (
	"NewsAggregator/entity/article"
	"NewsAggregator/entity/source"
	"NewsAggregator/filter_service"
	"NewsAggregator/news_service"
)

type NewsAggregator struct{}

func NewNewsAggregator() *NewsAggregator {
	return &NewsAggregator{}
}

func (na *NewsAggregator) Aggregate(sources []string, filters ...filter_service.FilterService) ([]article.Article, string) {
	sourceNames := filterUnique(sources)
	var sourceNameObjects []source.Name

	for _, name := range sourceNames {
		sourceNameObjects = append(sourceNameObjects, source.Name(name))
	}

	articles, errorMessage := news.FindNewsByResources(sourceNameObjects)
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

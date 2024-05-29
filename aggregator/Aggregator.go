package aggregator

import (
	"NewsAggregator/entity/article"
	"NewsAggregator/filter_service"
)

type Aggregator interface {
	Aggregate(sources []string, filters ...filter_service.FilterService) ([]article.Article, string)
}

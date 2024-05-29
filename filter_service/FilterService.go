package filter_service

import . "NewsAggregator/entity/article"

type FilterService interface {
	Filter(articles []Article) []Article
}

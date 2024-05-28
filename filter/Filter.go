package filter

import . "NewsAggregator/entity/article"

type Filter interface {
	Filter(articles []Article) []Article
}

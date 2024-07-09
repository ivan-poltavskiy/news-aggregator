package filter

import "news-aggregator/entity/news"

// NewsFilter defines a filtering service that can filter a slice of news based on specific criteria.
type NewsFilter interface {
	// Filter filters the given slice of news and returns a filtered slice of news.
	Filter(articles []news.News) []news.News
}

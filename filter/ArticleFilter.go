package filter

import "NewsAggregator/entity/article"

// ArticleFilter defines a filtering service that can filter a slice of articles based on specific criteria.
type ArticleFilter interface {
	// Filter filters the given slice of articles and returns a filtered slice of articles.
	Filter(articles []article.Article) []article.Article
}

package client

import (
	"NewsAggregator/entity/article"
)

type Client interface {
	FetchArticles() []article.Article
	Print(articles []article.Article)
}

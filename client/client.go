package client

import "news-aggregator/entity/article"

type Client interface {
	//FetchArticles collect the articles by some rules defined in the implementations.
	FetchArticles() ([]article.Article, error)
	//Print outputs the transferred articles.
	Print(articles []article.Article)
}

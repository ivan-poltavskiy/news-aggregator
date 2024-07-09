package client

import "news-aggregator/entity/news"

// Sorter sorts input data according to specified rules.
type Sorter interface {
	//SortNews allocates input articles.
	SortNews(articles []news.News, sortBy string) ([]news.News, error)
}

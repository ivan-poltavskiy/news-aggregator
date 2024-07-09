package client

import "news-aggregator/entity/news"

type Client interface {
	//FetchNews collect the news by some rules defined in the implementations.
	FetchNews() ([]news.News, error)
	//Print outputs the transferred news.
	Print(news []news.News)
}

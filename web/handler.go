package main

import (
	"news-aggregator/web/news"
	"news-aggregator/web/source"
)

type Handler interface {
	source.SourceHandler
	news.NewsHandler
}

func NewHandler() *Handler {
	return &Handler{}
}

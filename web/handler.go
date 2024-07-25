package main

import (
	"news-aggregator/storage"
	"news-aggregator/web/news"
	"news-aggregator/web/source"
)

type Handler struct {
	*source.SourceHandler
	*news.NewsHandler
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{
		SourceHandler: source.NewSourceHandler(storage),
		NewsHandler:   news.NewNewsHandler(storage),
	}
}

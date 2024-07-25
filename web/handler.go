package main

import (
	"news-aggregator/storage"
	"news-aggregator/web/news"
	"news-aggregator/web/source"
)

type Handler struct {
	SourceHandler *source.HandlerForSources
	NewsHandler   *news.HandlerForNews
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{
		SourceHandler: source.NewSourceHandler(storage),
		NewsHandler:   news.NewNewsHandler(storage),
	}
}

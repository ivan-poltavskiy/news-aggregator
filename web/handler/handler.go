package handler

import (
	"news-aggregator/storage"
	"news-aggregator/web/news"
	"news-aggregator/web/source"
)

// Handler is an abstract interface for work with different resources
type Handler interface {
	GetSourceHandler() *source.HandlerForSources
	GetNewsHandler() *news.HandlerForNews
}

// handler is an implementation of the Handler interface for news and sources handlers
type handler struct {
	SourceHandler *source.HandlerForSources
	NewsHandler   *news.HandlerForNews
}

// NewHandler returns a new instance of the Handler interface
func NewHandler(storage storage.Storage) Handler {
	return &handler{
		SourceHandler: source.NewSourceHandler(storage),
		NewsHandler:   news.NewNewsHandler(storage),
	}
}

// GetSourceHandler returns the SourceHandler
func (h *handler) GetSourceHandler() *source.HandlerForSources {
	return h.SourceHandler
}

// GetNewsHandler returns the NewsHandler
func (h *handler) GetNewsHandler() *news.HandlerForNews {
	return h.NewsHandler
}

package client

import (
	"NewsAggregator/entity/article"
)

// Client defines an interface for a client that can fetch and print collector articles.
type Client interface {
	// FetchArticles fetches and returns a slice of articles.
	//
	// Returns:
	// - A slice of articles fetched by the client.
	FetchArticles() []article.Article
	// Print prints the provided articles.
	//
	// Parameters:
	// - articles: a slice of articles to be printed.
	Print(articles []article.Article)
}

package parser

import (
	"NewsAggregator/entity/article"
	"NewsAggregator/entity/source"
)

// A Parser to analyze a source and retrieve a list of articles from that source.
type Parser interface {

	// ParseSource returns a list of the source's articles by his path.
	ParseSource(path source.PathToFile) []article.Article
}

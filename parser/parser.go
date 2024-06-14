package parser

import (
	"errors"
	"news_aggregator/entity/article"
	"news_aggregator/entity/source"
	"news_aggregator/parser/html"
)

// Parsers stores the mapping of source types to their corresponding parsers.
var Parsers map[source.Type]Parser

// Initialize initializes a parser map with available parsers for different file types.
func Initialize() {
	Parsers = map[source.Type]Parser{
		"RSS":  Rss{},
		"JSON": Json{},
		"Html": html.UsaToday{},
	}
}

// GetParserBySourceType returns the parser that is required for parsing files of the passed type.
func GetParserBySourceType(typeOfSource source.Type) (Parser, error) {
	parser, exist := Parsers[typeOfSource]
	if !exist {
		return nil, errors.New("parser not exist")
	}
	return parser, nil
}

// A Parser to analyze a source and retrieve a list of articles from that source.
type Parser interface {

	// ParseSource returns a list of the source's articles by his path.
	ParseSource(path source.PathToFile, name source.Name) ([]article.Article, error)
}

package collector

import (
	"errors"
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	"news-aggregator/parser"
	"news-aggregator/parser/html"
)

// Parsers manages the mapping of source types to their corresponding parsers.
type Parsers struct {
	parsers map[source.Type]Parser
}

// GetDefaultParsers initializes and returns a new Parsers with available parsers for different file types.
func GetDefaultParsers() *Parsers {
	return &Parsers{
		parsers: map[source.Type]Parser{
			source.RSS:      parser.Rss{},
			source.JSON:     parser.Json{},
			source.UsaToday: html.UsaToday{},
		},
	}
}

// GetParserBySourceType returns the parser that is required for parsing files of the passed type.
func (pm *Parsers) GetParserBySourceType(typeOfSource source.Type) (Parser, error) {
	parser, exist := pm.parsers[typeOfSource]
	if !exist {
		return nil, errors.New("parser not exist")
	}
	return parser, nil
}

// A Parser analyzes a source and retrieves a list of articles from that source.
type Parser interface {

	// Parse returns a list of the source's news by its path.
	Parse(path source.PathToFile, name source.Name) ([]news.News, error)
}

package parser

import (
	"NewsAggregator/entity/article"
	"NewsAggregator/entity/source"
	"fmt"
)

// ParserMap stores the mapping of source types to their corresponding parsers.
var ParserMap map[source.Type]Parser

// InitializeParserMap initializes the parser map with available parsers.
func InitializeParserMap() {
	ParserMap = map[source.Type]Parser{
		"RSS":  RssParser{},
		"JSON": JsonParser{},
		"Html": HtmlParser{},
	}
}

// GetParserBySourceType returns the parser that is required for parsing files of the passed type.
func GetParserBySourceType(typeOfSource source.Type) Parser {
	parser, exist := ParserMap[typeOfSource]
	if !exist {
		fmt.Println("Wrong Source", typeOfSource)
		return nil
	}
	return parser
}

// A Parser to analyze a source and retrieve a list of articles from that source.
type Parser interface {

	// ParseSource returns a list of the source's articles by his path.
	ParseSource(path source.PathToFile) []article.Article
}

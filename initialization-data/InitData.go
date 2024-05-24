package initialization_data

import (
	. "NewsAggregator/entity/source"
	"NewsAggregator/parser"
)

var Sources []Source

var ParserMap map[Type]parser.Parser

// InitializeSource creates the necessary data for the program to run.
func InitializeSource() {

	ParserMap = map[Type]parser.Parser{
		"RSS":  parser.RssParser{},
		"JSON": parser.JsonParser{},
		"Html": parser.HtmlParser{},
	}

	Sources = []Source{
		{Name: "bbc", PathToFile: "resources/bbc-world-category-19-05-24.xml", SourceType: "RSS", Id: 1},
		{Name: "nbc", PathToFile: "resources/nbc-news.json", SourceType: "JSON", Id: 2},
		{Name: "abc", PathToFile: "resources/abcnews-international-category-19-05-24.xml", SourceType: "RSS", Id: 3},
		{Name: "washington", PathToFile: "resources/washingtontimes-world-category-19-05-24.xml", SourceType: "RSS", Id: 4},
		{Name: "usatoday", PathToFile: "resources/usatoday-world-news.html", SourceType: "Html", Id: 5},
	}
}

package main

import (
	. "NewsAggregator/aggregator"
	. "NewsAggregator/client"
	. "NewsAggregator/entity/source"
	"NewsAggregator/news_service"
	"NewsAggregator/parser"
)

func main() {

	news.InitializeSource([]Source{
		{Name: "bbc", PathToFile: "resources/bbc-world-category-19-05-24.xml", SourceType: "RSS", Id: 1},
		{Name: "nbc", PathToFile: "resources/nbc-news.json", SourceType: "JSON", Id: 2},
		{Name: "abc", PathToFile: "resources/abcnews-international-category-19-05-24.xml", SourceType: "RSS", Id: 3},
		{Name: "washington", PathToFile: "resources/washingtontimes-world-category-19-05-24.xml", SourceType: "RSS", Id: 4},
		{Name: "usatoday", PathToFile: "resources/usatoday-world-news.html", SourceType: "Html", Id: 5},
	})
	parser.InitializeParserMap()

	aggregator := NewNewsAggregator()
	cli := NewCommandLineClient(aggregator)
	articles := cli.FetchArticles()
	if articles != nil {
		cli.Print(articles)
	}
}

package main

import (
	"NewsAggregator/aggregator"
	"NewsAggregator/client"
	"NewsAggregator/collector"
	. "NewsAggregator/entity/source"
	"NewsAggregator/parser"
)

func main() {

	collector.InitializeSource([]Source{
		{Name: "bbc", PathToFile: "resources/bbc-world-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "nbc", PathToFile: "resources/nbc-collector.json", SourceType: "JSON"},
		{Name: "abc", PathToFile: "resources/abcnews-international-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "washington", PathToFile: "resources/washingtontimes-world-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "usatoday", PathToFile: "resources/usatoday-world-collector.html", SourceType: "Html"},
	})
	parser.InitializeParserMap()

	newsAggregator := aggregator.New()
	cli := client.NewCommandLine(newsAggregator)
	articles := cli.FetchArticles()
	if articles != nil {
		cli.Print(articles)
	}
}

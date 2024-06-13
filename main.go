package main

import (
	"news_aggregator/aggregator"
	"news_aggregator/client"
	"news_aggregator/collector"
	"news_aggregator/entity/source"
	"news_aggregator/parser"
)

func main() {

	collector.InitializeSource([]source.Source{
		{Name: "bbc", PathToFile: "resources/bbc-world-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "nbc", PathToFile: "resources/nbc-news.json", SourceType: "JSON"},
		{Name: "abc", PathToFile: "resources/abcnews-international-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "washington", PathToFile: "resources/washingtontimes-world-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "usatoday", PathToFile: "resources/usatoday-world-news.html", SourceType: "Html"},
	})
	parser.Initialize()

	newsAggregator := aggregator.New()
	cli := client.NewCommandLine(newsAggregator)
	articles := cli.FetchArticles()
	cli.Print(articles)

}

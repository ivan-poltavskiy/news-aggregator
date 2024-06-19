package main

import (
	"news_aggregator/aggregator"
	"news_aggregator/client"
	"news_aggregator/collector"
	"news_aggregator/entity/source"
)

func main() {

	sources := []source.Source{
		{Name: "bbc", PathToFile: "resources/bbc-world-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "nbc", PathToFile: "resources/nbc-news.json", SourceType: "JSON"},
		{Name: "abc", PathToFile: "resources/abcnews-international-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "washington", PathToFile: "resources/washingtontimes-world-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "usatoday", PathToFile: "resources/usatoday-world-news.html", SourceType: "UsaToday"},
	}
	articleCollector := collector.New(sources)

	newsAggregator := aggregator.New(articleCollector)
	cli := client.NewCommandLine(newsAggregator)
	articles, err := cli.FetchArticles()
	if err != nil {
		println(err.Error())
	}
	cli.Print(articles)

}

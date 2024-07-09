package main

import (
	"news-aggregator/aggregator"
	"news-aggregator/client"
	"news-aggregator/collector"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
)

func main() {

	sources, err := source.LoadExistingSourcesFromStorage(constant.PathToStorage)
	if err != nil {
		println(err)
		return
	}
	articleCollector := collector.New(sources)

	newsAggregator := aggregator.New(articleCollector)
	cli := client.NewCommandLine(newsAggregator)
	articles, err := cli.FetchNews()
	if err != nil {
		println(err.Error())
	}
	cli.Print(articles)

}

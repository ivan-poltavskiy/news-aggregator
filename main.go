package main

import (
	"news-aggregator/aggregator"
	"news-aggregator/client"
	"news-aggregator/collector"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
)

func main() {

	sources := []source.Source{
		constant.BbcSource,
		constant.NbcSource,
		constant.AbcSource,
		constant.WashingtonSource,
		constant.UsaTodaySource,
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

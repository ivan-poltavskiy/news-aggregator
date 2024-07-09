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
	newsCollector := collector.New(sources)

	newsAggregator := aggregator.New(newsCollector)
	cli := client.NewCommandLine(newsAggregator)
	news, err := cli.FetchNews()
	if err != nil {
		println(err.Error())
	}
	cli.Print(news)

}

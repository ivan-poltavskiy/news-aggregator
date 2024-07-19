package main

import (
	"news-aggregator/aggregator"
	"news-aggregator/client"
	"news-aggregator/collector"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
	"news-aggregator/storage"
)

func main() {

	sourceStorage := storage.NewJsonSourceStorage(source.PathToFile(constant.PathToStorage))
	newsCollector := collector.New(sourceStorage)
	newsAggregator := aggregator.New(newsCollector)
	cli := client.NewCommandLine(newsAggregator)
	articles, err := cli.FetchNews()
	if err != nil {
		println(err.Error())
	}
	cli.Print(articles)

}

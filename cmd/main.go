package main

import (
	"github.com/sirupsen/logrus"
	"news-aggregator/aggregator"
	"news-aggregator/client"
	"news-aggregator/collector"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
	"news-aggregator/storage"
	newsStorage "news-aggregator/storage/news"
	sourceStorage "news-aggregator/storage/source"
)

func main() {
	newsJsonStorage, err := newsStorage.NewJsonNewsStorage(source.PathToFile(constant.PathToResources))
	if err != nil {
		logrus.Fatal(err)
	}
	jsonSourceStorage, err := sourceStorage.NewJsonSourceStorage(source.PathToFile(constant.PathToStorage))
	if err != nil {
		panic(err)
	}
	newStorage := storage.NewStorage(
		newsJsonStorage,
		jsonSourceStorage,
	)
	newsCollector := collector.New(newStorage)
	newsAggregator := aggregator.New(newsCollector)
	cli := client.NewCommandLine(newsAggregator)
	articles, err := cli.FetchNews()
	if err != nil {
		println(err.Error())
	}
	cli.Print(articles)

}

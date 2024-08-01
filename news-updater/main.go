package main

import (
	"github.com/sirupsen/logrus"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
	"news-aggregator/storage"
	newsStorage "news-aggregator/storage/news"
	sourceStorage "news-aggregator/storage/source"
	"news-updater/updater"
)

func main() {
	logrus.Info("res: " + constant.PathToResources)
	logrus.Info("str: " + constant.PathToStorage)
	newsJsonStorage, err := newsStorage.NewJsonStorage(source.PathToFile(constant.PathToResources))
	if err != nil {
		logrus.Fatal(err)
	}
	sourceJsonStorage, err := sourceStorage.NewJsonStorage(source.PathToFile(constant.PathToStorage))
	if err != nil {
		logrus.Fatal(err)
	}

	resourcesStorage := storage.NewStorage(newsJsonStorage, sourceJsonStorage)

	service := updater.Service{Storage: resourcesStorage}
	service.UpdateNews()
}

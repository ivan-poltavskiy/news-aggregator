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
	logrus.Info("Path to resources: " + constant.PathToResources)
	logrus.Info("Path to storage: " + constant.PathToStorage)

	newsJsonStorage, newsStorageErr := newsStorage.NewJsonStorage(source.PathToFile(constant.PathToResources))

	if newsStorageErr != nil {
		logrus.Fatal(newsStorageErr)
	}

	sourceJsonStorage, sourceStorageErr := sourceStorage.NewJsonStorage(source.PathToFile(constant.PathToStorage))

	if sourceStorageErr != nil {
		logrus.Fatal(sourceStorageErr)
	}

	resourcesStorage := storage.NewStorage(newsJsonStorage, sourceJsonStorage)

	service := updater.Service{Storage: resourcesStorage}
	service.UpdateNews()

}

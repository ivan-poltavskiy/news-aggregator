package news

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	"os"
)

type jsonResourcesStorage struct {
	pathToStorage source.PathToFile
}

// NewJsonResourcesStorage create new instance of storage in JSON file
func NewJsonResourcesStorage(pathToStorage source.PathToFile) NewsStorage {
	if pathToStorage == "" {
		panic("pathToStorage is empty")
	}
	return &jsonResourcesStorage{pathToStorage}
}

func (jsonStorage *jsonResourcesStorage) SaveNews(jsonFilePath string, news []news.News) error {
	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		logrus.Error("Failed to create JSON file: ", err)
		return fmt.Errorf("failed to create JSON file")
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			logrus.Error("Failed to close the JSON file: ", err)
		}
	}(jsonFile)

	if err := json.NewEncoder(jsonFile).Encode(news); err != nil {
		logrus.Error("Failed to encode articles to JSON file: ", err)
		return fmt.Errorf("failed to encode articles to JSON file")
	}

	return nil
}

func (jsonStorage *jsonResourcesStorage) GetNews(jsonFilePath string) ([]news.News, error) {
	var existingArticles []news.News

	if _, err := os.Stat(jsonFilePath); err == nil {
		jsonFile, err := os.Open(jsonFilePath)
		if err != nil {
			logrus.Error("Failed to open existing JSON file: ", err)
			return nil, fmt.Errorf("failed to open existing JSON file")
		}
		defer func(jsonFile *os.File) {
			err := jsonFile.Close()
			if err != nil {
				logrus.Error("Failed to close the existing JSON file: ", err)
			}
		}(jsonFile)

		if err := json.NewDecoder(jsonFile).Decode(&existingArticles); err != nil {
			logrus.Error("Failed to decode existing articles from JSON file: ", err)
			return nil, fmt.Errorf("failed to decode existing articles from JSON file")
		}
	}

	return existingArticles, nil
}

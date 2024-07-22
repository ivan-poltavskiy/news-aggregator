package news

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"news-aggregator/constant"
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	sourceStorage "news-aggregator/storage/source"
	"os"
	"path/filepath"
)

type jsonNewsStorage struct {
	pathToStorage source.PathToFile
}

// NewJsonNewsStorage create new instance of storage in JSON file
func NewJsonNewsStorage(pathToStorage source.PathToFile) NewsStorage {
	if pathToStorage == "" {
		panic("pathToStorage is empty")
	}
	return &jsonNewsStorage{pathToStorage}
}

// SaveNews saves the provided news articles to the specified JSON file.
func (jsonStorage *jsonNewsStorage) SaveNews(currentSource source.Source, news []news.News) (source.Source, error) {

	directoryPath := filepath.ToSlash(filepath.Join(constant.PathToResources, string(currentSource.Name)))

	if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
		logrus.Error("Failed to create directory: ", err)
		return source.Source{}, fmt.Errorf("failed to create directory")
	}

	jsonFilePath := filepath.ToSlash(filepath.Join(constant.PathToResources, string(currentSource.Name), string(currentSource.Name)+".json"))

	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		logrus.Error("Failed to create JSON file: ", err)
		return source.Source{}, fmt.Errorf("failed to create JSON file")
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			logrus.Error("Failed to close the JSON file: ", err)
		}
	}(jsonFile)

	if err := json.NewEncoder(jsonFile).Encode(news); err != nil {
		logrus.Error("Failed to encode articles to JSON file: ", err)
		return source.Source{}, fmt.Errorf("failed to encode articles to JSON file")
	}

	logrus.Info("jsonNewsStorage: Articles successfully parsed and saved to: ", jsonFilePath)
	currentSource.PathToFile = source.PathToFile(jsonFilePath)
	return currentSource, nil
}

// GetNews retrieves news articles from the specified JSON file.
func (jsonStorage *jsonNewsStorage) GetNews(jsonFilePath string) ([]news.News, error) {
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

func (jsonStorage *jsonNewsStorage) GetNewsBySourceName(sourceName source.Name, sourceStorage sourceStorage.Storage) ([]news.News, error) {
	source, err := sourceStorage.GetSourceByName(sourceName)
	if err != nil {
		logrus.Error("Failed to get source by name: ", err)
		return nil, err
	}
	news, err := jsonStorage.GetNews(string(source.PathToFile))
	if err != nil {
		logrus.Error("Failed to get source by path: ", err)
		return nil, err
	}
	return news, nil
}

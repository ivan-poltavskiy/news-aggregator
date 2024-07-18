package storage

import (
	"bufio"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"news-aggregator/entity/source"
	"os"
)

type JsonStorage struct {
	pathToStorage source.PathToFile
}

func NewJsonStorage(pathToStorage source.PathToFile) *JsonStorage {
	if pathToStorage == "" {
		panic("pathToStorage is empty")
	}
	return &JsonStorage{pathToStorage}
}

func (storage *JsonStorage) SaveSource() {

}

// GetSources returns the all sources from the JSON storage
func (storage *JsonStorage) GetSources() ([]source.Source, error) {
	logrus.Info("JsonStorage: Starting loading the existing sources from storage")
	file, err := os.Open(string(storage.pathToStorage))
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.Error("Source: Error closing file: ", err)
		}
	}(file)

	reader := bufio.NewReader(file)
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var sources []source.Source
	err = json.Unmarshal(content, &sources)
	if err != nil {
		return nil, err
	}

	return sources, nil
}
func (storage *JsonStorage) DeleteSource() {

}

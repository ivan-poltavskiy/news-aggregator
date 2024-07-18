package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
	"os"
	"path/filepath"
	"strings"
)

type JsonSourceStorage struct {
	pathToStorage source.PathToFile
}

func NewJsonSourceStorage(pathToStorage source.PathToFile) *JsonSourceStorage {
	if pathToStorage == "" {
		panic("pathToStorage is empty")
	}
	return &JsonSourceStorage{pathToStorage}
}

func (storage *JsonSourceStorage) SaveSource(source source.Source) error {
	logrus.Info("JsonSourceStorage: Starting to save the source to storage")

	// Read existing sources
	existingSources, err := storage.GetSources()
	if err != nil && !os.IsNotExist(err) {
		logrus.Error("JsonSourceStorage: Failed to read existing sources: ", err)
		return err
	}

	// Check if the source already exists
	for _, existingSource := range existingSources {
		if existingSource.Name == source.Name {
			logrus.Info("JsonSourceStorage: Source already exists, skipping save")
			return nil
		}
	}

	// Add the new source
	existingSources = append(existingSources, source)

	// Save the updated sources list
	file, err := os.Create(string(storage.pathToStorage))
	if err != nil {
		logrus.Error("JsonSourceStorage: Failed to create storage file: ", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.Error("JsonSourceStorage: Error closing storage file: ", err)
		}
	}(file)

	err = json.NewEncoder(file).Encode(existingSources)
	if err != nil {
		logrus.Error("JsonSourceStorage: Failed to encode sources to JSON: ", err)
		return err
	}

	logrus.Info("JsonSourceStorage: Source successfully saved to storage")
	return nil
}

// GetSources returns the all sources from the JSON storage
func (storage *JsonSourceStorage) GetSources() ([]source.Source, error) {
	logrus.Info("JsonSourceStorage: Starting loading the existing sources from storage")
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
func (storage *JsonSourceStorage) DeleteSourceByName(name string) error {
	var updatedSources []source.Source
	found := false
	definedSources, err := storage.GetSources()
	if err != nil {
		return err
	}
	for _, currentSource := range definedSources {
		if strings.ToLower(string(currentSource.Name)) != strings.ToLower(name) {
			updatedSources = append(updatedSources, currentSource)
		} else {
			found = true
			directoryPath := filepath.Join(constant.PathToResources, strings.ToLower(name))
			err := os.RemoveAll(directoryPath)
			if err != nil {
				logrus.Errorf("Failed to delete source directory %s: %v", directoryPath, err)
				return err
			}
			logrus.Infof("Deleted source directory: %s", directoryPath)
		}
	}

	if !found {
		return fmt.Errorf("source not found: %s", name)
	}
	for _, s := range updatedSources {
		err := storage.SaveSource(s)
		if err != nil {
			logrus.Errorf("Failed to write updated sources to file: %v", err)
			return err
		}

	}
	logrus.Info("Updated sources written to file successfully")

	return nil
}

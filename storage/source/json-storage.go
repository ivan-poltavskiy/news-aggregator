package source

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

type jsonSourceStorage struct {
	pathToStorage source.PathToFile
}

// NewJsonSourceStorage create new instance of storage in JSON file
func NewJsonSourceStorage(pathToStorage source.PathToFile) Storage {
	if pathToStorage == "" {
		panic("pathToStorage is empty")
	}
	return &jsonSourceStorage{pathToStorage}
}

// SaveSource load the input source to the storage
func (storage *jsonSourceStorage) SaveSource(source source.Source) error {
	logrus.Info("jsonSourceStorage: Starting to save the source to storage")

	// Read existing sources
	existingSources, err := storage.GetSources()
	if err != nil && !os.IsNotExist(err) {
		logrus.Error("jsonSourceStorage: Failed to read existing sources: ", err)
		return err
	}

	// Check if the source already exists
	for _, existingSource := range existingSources {
		if existingSource.Name == source.Name {
			logrus.Info("jsonSourceStorage: Source already exists, skipping save")
			return nil
		}
	}

	// Add the new source
	existingSources = append(existingSources, source)

	// Save the updated sources list
	file, err := os.Create(string(storage.pathToStorage))
	if err != nil {
		logrus.Error("jsonSourceStorage: Failed to create storage file: ", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.Error("jsonSourceStorage: Error closing storage file: ", err)
		}
	}(file)

	err = json.NewEncoder(file).Encode(existingSources)
	if err != nil {
		logrus.Error("jsonSourceStorage: Failed to encode sources to JSON: ", err)
		return err
	}

	logrus.Info("jsonSourceStorage: Source successfully saved to storage")
	return nil
}

// GetSources returns the all sources from the JSON storage
func (storage *jsonSourceStorage) GetSources() ([]source.Source, error) {
	logrus.Info("jsonSourceStorage: Starting loading the existing sources from storage")
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

// DeleteSourceByName remove the source from JSON storage by the name of this source
func (storage *jsonSourceStorage) DeleteSourceByName(name string) error {
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
	file, err := os.Create(string(storage.pathToStorage))
	if err != nil {
		logrus.Error("jsonSourceStorage: Failed to create storage file: ", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.Error("jsonSourceStorage: Error closing storage file: ", err)
		}
	}(file)

	err = json.NewEncoder(file).Encode(updatedSources)
	if err != nil {
		logrus.Error("jsonSourceStorage: Failed to encode sources to JSON: ", err)
		return err
	}

	logrus.Info("jsonSourceStorage: Source successfully deleted from storage")
	return nil
}
func (storage *jsonSourceStorage) IsSourceExists(name source.Name) bool {
	sources, err := storage.GetSources()
	if err != nil {
		logrus.Error("IsSourceExists: ", err)
		return false
	}
	for _, s := range sources {
		if s.Name == name {
			logrus.Info("IsSourceExists: Source exists: ", name)
			return true
		}
	}
	logrus.Info("IsSourceExists: Source does not exist: ", name)
	return false
}

func (storage *jsonSourceStorage) GetSourceByName(name source.Name) (source.Source, error) {
	logrus.Info("jsonSourceStorage: Starting to get source by name from storage")

	sources, err := storage.GetSources()
	if err != nil {
		logrus.Error("jsonSourceStorage: Failed to get sources: ", err)
		return source.Source{}, err
	}

	for _, s := range sources {
		if strings.ToLower(string(s.Name)) == strings.ToLower(string(name)) {
			logrus.Info("jsonSourceStorage: Source found: ", name)
			return s, nil
		}
	}
	logrus.Info("jsonSourceStorage: source not found", name)
	return source.Source{}, nil
}

package source

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
	"news-aggregator/storage"
	"os"
	"path/filepath"
	"strings"
)

type jsonStorage struct {
	pathToStorage source.PathToFile
}

// NewJsonStorage create new instance of storage in JSON file
func NewJsonStorage(pathToStorage source.PathToFile) (storage.Source, error) {
	if pathToStorage == "" {
		return nil, fmt.Errorf("NewJsonStorage: pathToStorage is empty")
	}
	return &jsonStorage{pathToStorage}, nil
}

// SaveSource load the input source to the storage
func (storage *jsonStorage) SaveSource(source source.Source) error {
	logrus.Info("jsonStorage: Starting to save the source to storage")

	existingSources, err := storage.GetSources()
	if err != nil && !os.IsNotExist(err) {
		logrus.Error("jsonStorage: Failed to read existing sources: ", err)
		return err
	}

	for _, existingSource := range existingSources {
		if existingSource.Name == source.Name {
			logrus.Info("jsonStorage: Source already exists, skipping save")
			return nil
		}
	}

	existingSources = append(existingSources, source)

	file, err := os.Create(string(storage.pathToStorage)) // Use os.Create to create or truncate the file
	if err != nil {
		logrus.Error("jsonStorage: Failed to create storage file: ", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.Error("jsonStorage: Error closing file: ", err)
		}
	}(file)

	err = json.NewEncoder(file).Encode(existingSources)
	if err != nil {
		logrus.Error("jsonStorage: Failed to encode sources to JSON: ", err)
		return err
	}

	logrus.Info("jsonStorage: Source successfully saved to storage")
	return nil
}

func (storage *jsonStorage) GetSources() ([]source.Source, error) {
	logrus.Info("jsonStorage: Starting loading the existing sources from storage")
	file, err := os.Open(string(storage.pathToStorage))
	if err != nil {
		if os.IsNotExist(err) {
			return []source.Source{}, nil // Return empty slice if file does not exist
		}
		logrus.Error("jsonStorage: Failed to open storage file: ", err)
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.Error("jsonStorage: Error closing storage file: ", err)
		}
	}(file)

	reader := bufio.NewReader(file)
	content, err := io.ReadAll(reader)
	if err != nil {
		logrus.Error("jsonStorage: Failed to read storage file: ", err)
		return nil, err
	}
	var sources []source.Source
	if len(content) != 0 {
		err = json.Unmarshal(content, &sources)
		if err != nil {
			logrus.Error("jsonStorage: Failed to unmarshal sources: ", err)
			return nil, err
		}
	}
	return sources, nil
}

// DeleteSourceByName remove the source from JSON storage by the name of this source
func (storage *jsonStorage) DeleteSourceByName(name source.Name) error {
	var updatedSources []source.Source
	found := false
	definedSources, err := storage.GetSources()
	if err != nil {
		return err
	}
	for _, currentSource := range definedSources {
		if strings.ToLower(string(currentSource.Name)) != strings.ToLower(string(name)) {
			updatedSources = append(updatedSources, currentSource)
		} else {
			found = true
			directoryPath := filepath.Join(constant.PathToResources, strings.ToLower(string(name)))
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
		logrus.Error("jsonStorage: Failed to create storage file: ", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.Error("jsonStorage: Error closing storage file: ", err)
		}
	}(file)

	err = json.NewEncoder(file).Encode(updatedSources)
	if err != nil {
		logrus.Error("jsonStorage: Failed to encode sources to JSON: ", err)
		return err
	}

	logrus.Info("jsonStorage: Source successfully deleted from storage")
	return nil
}
func (storage *jsonStorage) IsSourceExists(name source.Name) bool {
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

func (storage *jsonStorage) GetSourceByName(name source.Name) (source.Source, error) {
	logrus.Info("jsonStorage: Starting to get source by name from storage")

	sources, err := storage.GetSources()
	if err != nil {
		logrus.Error("jsonStorage: Failed to get sources: ", err)
		return source.Source{}, err
	}

	for _, s := range sources {
		if strings.ToLower(string(s.Name)) == strings.ToLower(string(name)) {
			logrus.Info("jsonStorage: Source found: ", name)
			return s, nil
		}
	}
	logrus.Info("jsonStorage: source not found: ", name)
	return source.Source{}, nil
}

// UpdateSource updates existing source in the JSON storage
func (storage *jsonStorage) UpdateSource(updatedSource source.Source, currentName string) error {
	logrus.Info("jsonStorage: Starting to update the source in storage")

	existingSources, err := storage.GetSources()
	if err != nil {
		logrus.Error("jsonStorage: Failed to get sources: ", err)
		return err
	}

	sourceFound := false

	for i, existingSource := range existingSources {
		if strings.ToLower(string(existingSource.Name)) == strings.ToLower(currentName) {
			existingSources[i] = updatedSource
			sourceFound = true
			logrus.Info("jsonStorage: Source updated: ", updatedSource.Name)
			break
		}
	}

	if !sourceFound {
		logrus.Error("jsonStorage: Source not found: ", currentName)
		return fmt.Errorf("source with name '%s' not found", currentName)
	}

	file, err := os.Create(string(storage.pathToStorage))
	if err != nil {
		logrus.Error("jsonStorage: Failed to create storage file: ", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.Error("jsonStorage: Error closing file: ", err)
		}
	}(file)

	err = json.NewEncoder(file).Encode(existingSources)
	if err != nil {
		logrus.Error("jsonStorage: Failed to encode sources to JSON: ", err)
		return err
	}

	logrus.Info("jsonStorage: Source successfully updated in storage")
	return nil
}

package service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
	"os"
	"path/filepath"
	"strings"
)

func DeleteAndUpdateSources(name string) error {
	var updatedSources []source.Source
	found := false
	for _, currentSource := range ReadSourcesFromFile() {
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

	if err := WriteSourcesToFile(updatedSources); err != nil {
		logrus.Errorf("Failed to write updated sources to file: %v", err)
		return err
	}
	logrus.Info("Updated sources written to file successfully")

	return nil
}

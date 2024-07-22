package service

import (
	"github.com/sirupsen/logrus"
	"news-aggregator/storage/source"
)

// DeleteSourceByName removes the source from storage by name.
func DeleteSourceByName(name string, sourceStorage source.Storage) error {
	err := sourceStorage.DeleteSourceByName(name)
	if err != nil {
		logrus.Error("Error deleting source:", err)
		return err
	}
	return nil
}

package service

import (
	"github.com/sirupsen/logrus"
	"news-aggregator/storage"
)

func DeleteSourceByName(name string, sourceStorage storage.Storage) error {
	err := sourceStorage.DeleteSourceByName(name)
	if err != nil {
		logrus.Error("Error deleting source:", err)
		return err
	}
	return nil
}

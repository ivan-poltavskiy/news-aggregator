package service

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
	"os"
	"regexp"
	"strings"
)

func WriteSourcesToFile(sources []source.Source) error {
	file, err := os.Create(constant.PathToStorage)
	if err != nil {
		logrus.Error("WriteSourcesToFile: Error creating file ", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.Error("WriteSourcesToFile: Error closing file ", err)
		}
	}(file)

	encoder := json.NewEncoder(file)
	err = encoder.Encode(&sources)
	if err != nil {
		logrus.Error("WriteSourcesToFile: Error encoding sources ", err)
		return err
	}

	logrus.Info("WriteSourcesToFile: Sources were successfully written to file")
	return nil
}

func ReadSourcesFromFile() []source.Source {
	file, err := os.Open(constant.PathToStorage)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Warn("ReadSourcesFromFile: Sources file does not exist")
			return []source.Source{}
		}
		logrus.Error("ReadSourcesFromFile: Error opening sources file ", err)
		return nil
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.Error("ReadSourcesFromFile: Error closing file ", err)
		}
	}(file)

	var sources []source.Source
	if err := json.NewDecoder(file).Decode(&sources); err != nil {
		logrus.Error("ReadSourcesFromFile: Error decoding sources file ", err)
		return nil
	}

	logrus.Info("ReadSourcesFromFile: Sources were successfully read from file")
	return sources
}

func ExtractDomainName(url string) string {
	re := regexp.MustCompile(`https?://(www\.)?([^/]+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 3 {
		logrus.Warn("ExtractDomainName: Failed to extract domain name from URL: ", url)
		return "unknown"
	}
	domain := matches[2]
	domain = strings.Split(domain, ".")[0]
	logrus.Info("ExtractDomainName: Extracted domain name: ", domain)
	return domain
}

func IsSourceExists(name source.Name) bool {
	sources := ReadSourcesFromFile()
	for _, s := range sources {
		if s.Name == name {
			logrus.Info("IsSourceExists: Source exists: ", name)
			return true
		}
	}
	logrus.Info("IsSourceExists: Source does not exist: ", name)
	return false
}

func AddSourceToStorage(newSource source.Source) {
	sources := append(ReadSourcesFromFile(), newSource)

	file, err := os.Create(constant.PathToStorage)
	if err != nil {
		logrus.Error("AddSourceToStorage: Error creating sources file ", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.Error("AddSourceToStorage: Error closing file ", err)
		}
	}(file)

	if err := json.NewEncoder(file).Encode(sources); err != nil {
		logrus.Error("AddSourceToStorage: Error encoding sources to file ", err)
	} else {
		logrus.Info("AddSourceToStorage: Source added to storage: ", newSource.Name)
	}
}

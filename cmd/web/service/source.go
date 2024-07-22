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

// ReadSourcesFromStorage returns the entities of sources from the storage
func ReadSourcesFromStorage() []source.Source {
	file, err := os.Open(constant.PathToStorage)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Warn("ReadSourcesFromStorage: Sources file does not exist")
			return []source.Source{}
		}
		logrus.Error("ReadSourcesFromStorage: Error opening sources file ", err)
		return nil
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.Error("ReadSourcesFromStorage: Error closing file ", err)
		}
	}(file)

	var sources []source.Source
	if err := json.NewDecoder(file).Decode(&sources); err != nil {
		logrus.Error("ReadSourcesFromStorage: Error decoding sources file ", err)
		return nil
	}

	logrus.Info("ReadSourcesFromStorage: Sources were successfully read from file")
	return sources
}

// ExtractDomainName parse the url to get the resource domain
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
	sources := ReadSourcesFromStorage()
	for _, s := range sources {
		if s.Name == name {
			logrus.Info("IsSourceExists: Source exists: ", name)
			return true
		}
	}
	logrus.Info("IsSourceExists: Source does not exist: ", name)
	return false
}

// AddSourceToStorage add the entity of source to the storage
func AddSourceToStorage(newSource source.Source) {
	sources := append(ReadSourcesFromStorage(), newSource)

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

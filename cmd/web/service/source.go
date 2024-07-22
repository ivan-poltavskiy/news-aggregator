package service

import (
	"github.com/sirupsen/logrus"
	"news-aggregator/entity/source"
	"news-aggregator/storage"
	"regexp"
	"strings"
)

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

func IsSourceExists(name source.Name, sourceStorage storage.Storage) bool {
	sources, err := sourceStorage.GetSources()
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

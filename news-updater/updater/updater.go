package updater

import (
	"github.com/sirupsen/logrus"
	"news-aggregator/entity/source"
	"news-aggregator/storage"
	"news-aggregator/web/feed"
	"sync"
)

// Service represents the service for updating news.
type Service struct {
	Storage storage.Storage
}

// UpdateNews updates news for all sources.
func (service Service) UpdateNews() {
	logrus.Info("Starting update of news")
	sources, err := service.Storage.GetSources()
	if err != nil {
		logrus.Error("Failed to retrieve sources: ", err)
		return
	}

	var wg sync.WaitGroup
	for _, src := range sources {
		wg.Add(1)
		go func(src source.Source) {
			defer wg.Done()
			if src.SourceType == source.STORAGE {
				err := updateSourceNews(src, service.Storage)
				if err != nil {
					logrus.Error("Failed to update news for source: ", src.Name)
				}
			}
		}(src)
	}
	wg.Wait()
	logrus.Info("Update of news completed")
}

// updateSourceNews updates the news of the input source
func updateSourceNews(inputSource source.Source, storage storage.Storage) error {
	domainName := feed.ExtractDomainName(string(inputSource.Link))
	rssURL, err := feed.GetRssFeedLink(string(inputSource.Link))
	if err != nil {
		return err
	}

	currentNews, err := feed.ParseRssFeed(rssURL, domainName)
	if err != nil {
		return err
	}

	_, err = storage.SaveNews(inputSource, currentNews)
	if err != nil {
		return err
	}

	return nil
}

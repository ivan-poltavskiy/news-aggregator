package service

import (
	"github.com/sirupsen/logrus"
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	newsStorage "news-aggregator/storage/news"
	sourceStorage "news-aggregator/storage/source"
	"sync"
	"time"
)

// SaveNews saves the news to the storage
func SaveNews(sourceEntity source.Source, newsStorage newsStorage.NewsStorage, sourceStorage sourceStorage.Storage, parsedNews []news.News) (source.Source, error) {

	existingNews, err := newsStorage.GetNewsBySourceName(sourceEntity.Name, sourceStorage)
	if err != nil {
		return source.Source{}, err
	}

	newArticles := newsUnification(parsedNews, existingNews)
	if len(newArticles) == 0 {
		logrus.Info("No new parsed news to add")
		return sourceEntity, nil
	}

	existingNews = append(existingNews, newArticles...)

	sourceEntity, err = newsStorage.SaveNews(sourceEntity, existingNews)

	return sourceEntity, nil
}

// PeriodicallyUpdateNews updates news for all sources.
func PeriodicallyUpdateNews(sourceStorage sourceStorage.Storage, newsUpdatePeriod time.Duration, newsStorage newsStorage.NewsStorage) {
	ticker := time.NewTicker(newsUpdatePeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logrus.Info("Starting periodic update of news")
			sources, err := sourceStorage.GetSources()
			if err != nil {
				logrus.Error("Failed to retrieve sources: ", err)
				continue
			}

			var wg sync.WaitGroup
			for _, src := range sources {
				wg.Add(1)
				go func(src source.Source) {
					defer wg.Done()
					err := updateSourceNews(src, newsStorage, sourceStorage)
					if err != nil {
						logrus.Error("Failed to update news for source: ", src.Name, err)
					}
				}(src)
			}
			wg.Wait()
			logrus.Info("Periodic update of news completed")
		}
	}
}

// newsUnification checks whether there are articles from the new feed in the existing news, and if so, removes them
func newsUnification(articles []news.News, existingArticles []news.News) []news.News {
	existingTitles := make(map[string]struct{})
	for _, existingArticle := range existingArticles {
		existingTitles[existingArticle.Title.String()] = struct{}{}
	}

	var newArticles []news.News
	for _, newArticle := range articles {
		if _, exists := existingTitles[newArticle.Title.String()]; !exists {
			newArticles = append(newArticles, newArticle)
		}
	}

	return newArticles
}

// updateSourceNews updating the news of the input source
func updateSourceNews(inputSource source.Source, newsStorage newsStorage.NewsStorage, sourceStorage sourceStorage.Storage) error {
	domainName := ExtractDomainName(string(inputSource.Link))
	rssURL, err := getRssFeedLink(string(inputSource.Link))
	if err != nil {
		return err
	}

	currentNews, err := parseRssFeed(rssURL, domainName)
	if err != nil {
		return err
	}

	_, err = SaveNews(inputSource, newsStorage, sourceStorage, currentNews)
	if err != nil {
		return err
	}

	return nil
}

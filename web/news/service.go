package news

import (
	"github.com/sirupsen/logrus"
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	"news-aggregator/storage"
	"news-aggregator/web/feed"
	"sync"
	"time"
)

type Service struct {
	storage storage.Storage
}

func NewService(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

// SaveNews saves the news to the storage
func (service Service) SaveNews(sourceEntity source.Source, parsedNews []news.News) (source.Source, error) {

	existingNews, err := service.storage.GetNewsBySourceName(sourceEntity.Name, service.storage)
	if err != nil {
		return source.Source{}, err
	}

	newArticles := newsUnification(parsedNews, existingNews)
	if len(newArticles) == 0 {
		logrus.Info("No new parsed news to add")
		return sourceEntity, nil
	}

	existingNews = append(existingNews, newArticles...)

	sourceEntity, err = service.storage.SaveNews(sourceEntity, existingNews)

	return sourceEntity, nil
}

// PeriodicallyUpdateNews updates news for all sources.
func (service Service) PeriodicallyUpdateNews(newsUpdatePeriod time.Duration) {
	ticker := time.NewTicker(newsUpdatePeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logrus.Info("Starting periodic update of news")
			sources, err := service.storage.GetSources()
			if err != nil {
				logrus.Error("Failed to retrieve sources: ", err)
				continue
			}

			var wg sync.WaitGroup
			for _, src := range sources {
				wg.Add(1)
				go func(src source.Source) {
					defer wg.Done()
					if src.SourceType == source.STORAGE {
						err := updateSourceNews(src, service.storage)
						if err != nil {
							logrus.Error("Failed to update news for source: ", src.Name, err)
						}
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

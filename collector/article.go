package collector

import (
	"github.com/sirupsen/logrus"
	"news-aggregator/entity/article"
	"news-aggregator/entity/source"
	"strings"
)

type ArticleCollector struct {
	Sources []source.Source
	Parsers *Parsers
}

// New create new instance of collector
func New(sources []source.Source) *ArticleCollector {
	logrus.Info("Article Collector initialized")
	return &ArticleCollector{
		Sources: sources,
		Parsers: InitParsers(),
	}
}

// FindNewsByResourcesName returns the list of news from the passed sources.
func (articleCollector *ArticleCollector) FindNewsByResourcesName(sourcesNames []source.Name) ([]article.Article, error) {
	logrus.Info("Article collector: Start searching for articles by sources names: ", sourcesNames)

	var foundArticles []article.Article

	for _, sourceName := range sourcesNames {
		for _, currentSource := range articleCollector.Sources {
			articles, err := articleCollector.findNewsForCurrentSource(currentSource, sourceName)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"sourceName":    sourceName,
					"currentSource": currentSource,
				}).Error("Article collector: Error finding news for current source: ", err)
				return nil, err
			}
			foundArticles = append(foundArticles, articles...)
		}
	}
	logrus.Info("Article collector: Completed finding news by resource names: ", sourcesNames)
	return foundArticles, nil
}

// Returns the list of news from the passed source.
func (articleCollector *ArticleCollector) findNewsForCurrentSource(currentSource source.Source, name source.Name) ([]article.Article, error) {
	if strings.ToLower(string(currentSource.Name)) != strings.ToLower(string(name)) {
		return nil, nil
	}

	sourceParser, err := articleCollector.Parsers.GetParserBySourceType(currentSource.SourceType)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"sourceType": currentSource.SourceType,
		}).Error("Article collector: Error getting parser by source type: ", err)
		return []article.Article{}, err
	}

	logrus.WithFields(logrus.Fields{
		"source": currentSource,
		"name":   name,
	}).Info("Article collector: Parsing articles")

	articles, err := sourceParser.Parse(currentSource.PathToFile, name)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"source": currentSource,
			"name":   name,
		}).Error("Article collector: Error parsing articles: ", err)
		return nil, err
	}

	return articles, nil
}

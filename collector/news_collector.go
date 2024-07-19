package collector

import (
	"news-aggregator/aggregator"
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	"news-aggregator/storage"
	"strings"
)

type newsCollector struct {
	sourceStorage storage.Storage
	parsers       *Parsers
}

// New create new instance of collector
func New(sourceStorage storage.Storage) aggregator.Collector {
	return &newsCollector{sourceStorage: sourceStorage, parsers: GetDefaultParsers()}
}

// FindNewsByResourcesName returns the list of news from the passed sources.
func (newsCollector *newsCollector) FindNewsByResourcesName(sourcesNames []source.Name) ([]news.News, error) {
	var foundNews []news.News
	sources, err := newsCollector.sourceStorage.GetSources()
	if err != nil {
		return nil, err
	}
	for _, sourceName := range sourcesNames {
		for _, currentSource := range sources {
			if strings.ToLower(string(currentSource.Name)) == strings.ToLower(string(sourceName)) {
				news, err := newsCollector.findNewsForCurrentSource(currentSource, sourceName)
				if err != nil {
					return nil, err
				}
				foundNews = append(foundNews, news...)
			}
		}
	}
	return foundNews, nil
}

// Returns the list of news from the passed source.
func (newsCollector *newsCollector) findNewsForCurrentSource(currentSource source.Source, name source.Name) ([]news.News, error) {

	sourceParser, err := newsCollector.parsers.GetParserBySourceType(currentSource.SourceType)
	if err != nil {
		return []news.News{}, err
	}

	foundNews, err := sourceParser.Parse(currentSource.PathToFile, name)
	if err != nil {
		return nil, err
	}

	return foundNews, nil
}

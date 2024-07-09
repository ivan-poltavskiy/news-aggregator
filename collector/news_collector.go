package collector

import (
	"news-aggregator/aggregator"
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	"strings"
)

type newsCollector struct {
	sources []source.Source
	parsers *Parsers
}

// New create new instance of collector
func New(sources []source.Source) aggregator.Collector {
	return &newsCollector{sources: sources, parsers: GetDefaultParsers()}
}

// FindNewsByResourcesName returns the list of news from the passed sources.
func (newsCollector *newsCollector) FindNewsByResourcesName(sourcesNames []source.Name) ([]news.News, error) {

	var foundNews []news.News

	for _, sourceName := range sourcesNames {
		for _, currentSource := range newsCollector.sources {
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

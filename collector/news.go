package collector

import (
	"news-aggregator/aggregator"
	"news-aggregator/entity/article"
	"news-aggregator/entity/source"
	"strings"
)

type news struct {
	sources []source.Source
	parsers *Parsers
}

// New create new instance of collector
func New(sources []source.Source) aggregator.Collector {
	return &news{sources: sources, parsers: GetDefaultParsers()}
}

// FindNewsByResourcesName returns the list of news from the passed sources.
func (newsCollector *news) FindNewsByResourcesName(sourcesNames []source.Name) ([]article.Article, error) {

	var foundNews []article.Article

	for _, sourceName := range sourcesNames {
		for _, currentSource := range newsCollector.sources {
			if strings.ToLower(string(currentSource.Name)) == strings.ToLower(string(sourceName)) {
				articles, err := newsCollector.findNewsForCurrentSource(currentSource, sourceName)
				if err != nil {
					return nil, err
				}
				foundNews = append(foundNews, articles...)
			}
		}
	}
	return foundNews, nil
}

// Returns the list of news from the passed source.
func (newsCollector *news) findNewsForCurrentSource(currentSource source.Source, name source.Name) ([]article.Article, error) {

	sourceParser, err := newsCollector.parsers.GetParserBySourceType(currentSource.SourceType)
	if err != nil {
		return []article.Article{}, err
	}

	foundNews, err := sourceParser.Parse(currentSource.PathToFile, name)
	if err != nil {
		return nil, err
	}

	return foundNews, nil
}

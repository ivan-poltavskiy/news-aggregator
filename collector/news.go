package collector

import (
	"news-aggregator/aggregator"
	"news-aggregator/entity/article"
	"news-aggregator/entity/source"
	"strings"
)

type news struct {
	Sources []source.Source
	Parsers *Parsers
}

// New create new instance of collector
func New(sources []source.Source) aggregator.Collector {
	return &news{Sources: sources, Parsers: GetDefaultParsers()}
}

// FindNewsByResourcesName returns the list of news from the passed sources.
func (newsCollector *news) FindNewsByResourcesName(sourcesNames []source.Name) ([]article.Article, error) {

	var foundNews []article.Article

	for _, sourceName := range sourcesNames {
		for _, currentSource := range newsCollector.Sources {
			articles, err := newsCollector.findNewsForCurrentSource(currentSource, sourceName)
			if err != nil {
				return nil, err
			}
			foundNews = append(foundNews, articles...)
		}
	}
	return foundNews, nil
}

// Returns the list of news from the passed source.
func (newsCollector *news) findNewsForCurrentSource(currentSource source.Source, name source.Name) ([]article.Article, error) {

	if strings.ToLower(string(currentSource.Name)) != strings.ToLower(string(name)) {
		return nil, nil
	}

	sourceParser, err := newsCollector.Parsers.GetParserBySourceType(currentSource.SourceType)
	if err != nil {
		return []article.Article{}, err
	}

	foundNews, err := sourceParser.Parse(currentSource.PathToFile, name)
	if err != nil {
		return nil, err
	}

	return foundNews, nil
}

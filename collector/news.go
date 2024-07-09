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
	return &news{Sources: sources, Parsers: InitParsers()}
}

// FindNewsByResourcesName returns the list of news from the passed sources.
func (articleCollector *news) FindNewsByResourcesName(sourcesNames []source.Name) ([]article.Article, error) {

	var foundArticles []article.Article

	for _, sourceName := range sourcesNames {
		for _, currentSource := range articleCollector.Sources {
			articles, err := articleCollector.findNewsForCurrentSource(currentSource, sourceName)
			if err != nil {
				return nil, err
			}
			foundArticles = append(foundArticles, articles...)
		}
	}
	return foundArticles, nil
}

// Returns the list of news from the passed source.
func (articleCollector *news) findNewsForCurrentSource(currentSource source.Source, name source.Name) ([]article.Article, error) {

	if strings.ToLower(string(currentSource.Name)) != strings.ToLower(string(name)) {
		return nil, nil
	}

	sourceParser, err := articleCollector.Parsers.GetParserBySourceType(currentSource.SourceType)
	if err != nil {
		return []article.Article{}, err
	}

	articles, err := sourceParser.Parse(currentSource.PathToFile, name)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

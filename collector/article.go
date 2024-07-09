package collector

import (
	"news-aggregator/aggregator"
	"news-aggregator/entity/article"
	"news-aggregator/entity/source"
	"strings"
)

type articleCollector struct {
	Sources []source.Source
	Parsers *Parsers
}

// New create new instance of collector
func New(sources []source.Source) aggregator.Collector {
	return &articleCollector{Sources: sources, Parsers: InitParsers()}
}

// FindNewsByResourcesName returns the list of news from the passed sources.
func (articleCollector *articleCollector) FindNewsByResourcesName(sourcesNames []source.Name) ([]article.Article, error) {

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
func (articleCollector *articleCollector) findNewsForCurrentSource(currentSource source.Source, name source.Name) ([]article.Article, error) {

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

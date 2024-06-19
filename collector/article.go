package collector

import (
	"news_aggregator/entity/article"
	"news_aggregator/entity/source"
	"strings"
)

type ArticleCollector struct {
	Sources []source.Source
}

var parserManager *Parsers

// New create new instance of collector
func New(sources []source.Source) *ArticleCollector {
	parserManager = InitParsers()
	return &ArticleCollector{Sources: sources}
}

// FindNewsByResourcesName returns the list of news from the passed sources.
func (articleCollector *ArticleCollector) FindNewsByResourcesName(sourcesNames []source.Name) ([]article.Article, error) {

	var foundArticles []article.Article

	for _, name := range sourcesNames {
		for _, currentSourceType := range articleCollector.Sources {
			articles, err := articleCollector.findNewsForCurrentSource(currentSourceType, name)
			if err != nil {
				return nil, err
			}
			foundArticles = append(foundArticles, articles...)
		}
	}
	return foundArticles, nil
}

// Returns the list of news from the passed source.
func (articleCollector *ArticleCollector) findNewsForCurrentSource(currentSource source.Source, name source.Name) ([]article.Article, error) {

	if strings.ToLower(string(currentSource.Name)) != strings.ToLower(string(name)) {
		return nil, nil
	}

	sourceParser, err := parserManager.GetParserBySourceType(currentSource.SourceType)
	if err != nil {
		return []article.Article{}, err
	}

	articles, err := sourceParser.Parse(currentSource.PathToFile, name)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

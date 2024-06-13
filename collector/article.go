package collector

import (
	"news_aggregator/entity/article"
	"news_aggregator/entity/source"
	"news_aggregator/parser"
	"strings"
)

var Sources []source.Source

// FindByResourcesName returns the list of news from the passed sources.
func FindByResourcesName(sourcesNames []source.Name) ([]article.Article, error) {

	var foundNews []article.Article

	for _, name := range sourcesNames {
		for _, currentSourceType := range Sources {
			var err error
			foundNews, err = findForCurrentSource(currentSourceType, name, foundNews)
			if err != nil {
				return nil, err
			}
		}
	}
	return foundNews, nil
}

// InitializeSource initializes the resources that will be available for parsing.
func InitializeSource(sources []source.Source) {
	Sources = sources
}

// Returns the list of news from the passed source.
func findForCurrentSource(currentSource source.Source, name source.Name,
	allArticles []article.Article) ([]article.Article, error) {

	if strings.ToLower(string(currentSource.Name)) == strings.ToLower(string(name)) {
		currentParser, err := parser.GetParserBySourceType(currentSource.SourceType)
		if err != nil {
			return nil, err
		}
		articles, err := currentParser.ParseSource(currentSource.PathToFile)
		if err != nil {
			return nil, err
		}
		allArticles = append(allArticles, articles...)
	}
	return allArticles, nil
}

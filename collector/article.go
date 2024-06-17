package collector

import (
	"news_aggregator/entity/article"
	"news_aggregator/entity/source"
	"strings"
)

var Sources []source.Source

// FindNewsByResourcesName returns the list of news from the passed sources.
func FindNewsByResourcesName(sourcesNames []source.Name) ([]article.Article, error) {

	var foundArticles []article.Article

	for _, name := range sourcesNames {
		for _, currentSourceType := range Sources {
			articles, err := findNewsForCurrentSource(currentSourceType, name)
			if err != nil {
				return nil, err
			}
			foundArticles = append(foundArticles, articles...)
		}
	}
	return foundArticles, nil
}

// InitializeSource initializes the resources that will be available for parsing.
func InitializeSource(sources []source.Source) {
	Sources = sources
}

// Returns the list of news from the passed source.
func findNewsForCurrentSource(currentSource source.Source, name source.Name) ([]article.Article, error) {

	if strings.ToLower(string(currentSource.Name)) != strings.ToLower(string(name)) {
		return nil, nil
	}

	currentParser, err := GetParserBySourceType(currentSource.SourceType)
	if err != nil {
		return nil, err
	}

	articles, err := currentParser.ParseSource(currentSource.PathToFile, name)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

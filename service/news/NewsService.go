package news

import (
	"NewsAggregator/entity/article"
	"NewsAggregator/entity/source"
	"NewsAggregator/parser"
	"strings"
)

var Sources []source.Source

// FindByResourcesName returns the list of news from the passed sources.
func FindByResourcesName(sourcesNames []source.Name) ([]article.Article, string) {

	var foundNews []article.Article

	for _, name := range sourcesNames {
		for _, currentSourceType := range Sources {
			foundNews = findForCurrentSource(currentSourceType, name, foundNews)
		}
	}
	return foundNews, ""
}

// Returns the list of news from the passed source.
func findForCurrentSource(currentSourceType source.Source,
	name source.Name, allArticles []article.Article) []article.Article {

	if strings.ToLower(string(currentSourceType.Name)) == strings.ToLower(string(name)) {
		articles := parser.GetParserBySourceType(currentSourceType.SourceType).ParseSource(currentSourceType.PathToFile)
		allArticles = append(allArticles, articles...)

	}
	return allArticles
}

// InitializeSource initializes the resources that will be available for parsing.
func InitializeSource(sources []source.Source) {
	Sources = sources
}

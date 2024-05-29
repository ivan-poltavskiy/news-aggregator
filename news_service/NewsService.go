package news

import (
	. "NewsAggregator/entity/article"
	. "NewsAggregator/entity/source"
	"NewsAggregator/parser"
	"strings"
)

var Sources []Source

// FindNewsByResources returns the list of news from the passed sources.
func FindNewsByResources(sourcesNames []Name) ([]Article, string) {

	var foundNews []Article

	for _, name := range sourcesNames {
		for _, currentSourceType := range Sources {
			foundNews = findNewsForCurrentSource(currentSourceType, name, foundNews)
		}
	}
	return foundNews, ""
}

// Returns the list of news from the passed source.
func findNewsForCurrentSource(currentSourceType Source,
	name Name, allArticles []Article) []Article {

	if strings.ToLower(string(currentSourceType.Name)) == strings.ToLower(string(name)) {
		articles := parser.GetParserBySourceType(currentSourceType.SourceType).ParseSource(currentSourceType.PathToFile)
		allArticles = append(allArticles, articles...)

	}
	return allArticles
}

func InitializeSource(sources []Source) {
	Sources = sources
}

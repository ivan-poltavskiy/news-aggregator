package filter

import (
	. "NewsAggregator/entity/article"
	. "NewsAggregator/entity/source"
	. "NewsAggregator/initialization-data"
	"fmt"
	"strings"
	"time"
)

// FindNewsForAllResources returns the list of news from the passed sources.
func FindNewsForAllResources(sourcesNames []Name) ([]Article, string) {

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
		articles := GetParserBySourceType(currentSourceType.SourceType).ParseSource(currentSourceType.PathToFile)
		allArticles = append(allArticles, articles...)

	}
	return allArticles
}

// filterNewsByKeyword filters the incoming news list by keyword and returns the filtered list.
func filterNewsByKeyword(keyword string, articles []Article) []Article {
	var matchingOptions []Article

	for _, article := range articles {
		if strings.Contains(strings.ToLower(string(article.Title)), strings.ToLower(keyword)) {
			matchingOptions = append(matchingOptions, article)
		}
	}

	if len(matchingOptions) == 0 {
		fmt.Println("No matches found for this keyword.")
	}

	return matchingOptions
}

func NewsByKeywords(keywords []string, articles []Article) []Article {
	var matchingArticles []Article

	for _, keyword := range keywords {
		matchingArticles = append(matchingArticles, filterNewsByKeyword(keyword, articles)...)
	}

	return matchingArticles
}

// ByDate filters the list of incoming news by date and returns news whose publication date is between startDate and endDate
func ByDate(startDate time.Time, endDate time.Time, articles []Article) []Article {
	var matchingOptions []Article

	for _, article := range articles {
		if article.Date.After(startDate) && article.Date.Before(endDate) {
			matchingOptions = append(matchingOptions, article)
		}
	}
	if len(matchingOptions) == 0 {
		fmt.Println("No articles were found in this time period.")
	}

	return matchingOptions
}

package filterService

import (
	. "NewsAggregator/entity/article"
	. "NewsAggregator/entity/source"
	. "NewsAggregator/initialization-data"
	"fmt"
	"strings"
	"time"
)

// FindNewsForAllResources returns the list of news from the passed sources.
func FindNewsForAllResources(names []Name) ([]Article, string) {

	var foundNews []Article

	for _, name := range names {
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

// FilterNewsByKeyword filters the incoming news list by keyword and returns the filtered list.
func FilterNewsByKeyword(keyword string, articles []Article) []Article {
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

// FilterByDate filters the list of incoming news by date and returns news whose publication date is between startDate and endDate
func FilterByDate(startDate time.Time, endDate time.Time, articles []Article) []Article {
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

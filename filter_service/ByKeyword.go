package filter_service

import (
	. "NewsAggregator/entity/article"
	"fmt"
	"strings"
)

type ByKeyword struct {
	Keywords []string
}

func (f ByKeyword) Filter(articles []Article) []Article {
	var matchingArticles []Article
	for _, keyword := range f.Keywords {
		matchingArticles = append(matchingArticles, filterNewsByKeyword(keyword, articles)...)
	}
	return matchingArticles
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

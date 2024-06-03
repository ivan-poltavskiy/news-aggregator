package filter

import (
	"NewsAggregator/entity/article"
	"fmt"
	"strings"
)

// ByKeyword filters the slice of articles by provided keyword and returns
// the slice of matching articles.
type ByKeyword struct {
	Keywords []string
}

// Filter filters the incoming collector list from different sources by keywords.
func (f ByKeyword) Filter(articles []article.Article) []article.Article {
	var matchingArticles []article.Article
	for _, keyword := range f.Keywords {
		matchingArticles = append(matchingArticles, filterNewsByKeyword(keyword, articles)...)
	}
	return matchingArticles
}

// filterNewsByKeyword filters the incoming collector list by keyword and returns the filtered list.
func filterNewsByKeyword(keyword string, articles []article.Article) []article.Article {
	var matchingArticles []article.Article

	for _, article := range articles {
		if strings.Contains(strings.ToLower(string(article.Title)), strings.ToLower(keyword)) ||
			strings.Contains(strings.ToLower(string(article.Description)), strings.ToLower(keyword)) {
			matchingArticles = append(matchingArticles, article)
		}
	}

	if len(matchingArticles) == 0 {
		fmt.Println("No matches found for this keyword.")
	}

	return matchingArticles
}

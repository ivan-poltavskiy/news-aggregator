package filter

import (
	"fmt"
	"github.com/reiver/go-porterstemmer"
	"news_aggregator/entity/article"
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
		stemmedKeyword := porterstemmer.StemString(strings.ToLower(keyword))
		matchingArticles = append(matchingArticles, filterNewsByKeyword(stemmedKeyword, articles)...)
	}
	return matchingArticles
}

// filterNewsByKeyword filters the incoming collector list by keyword and returns the filtered list.
func filterNewsByKeyword(keyword string, articles []article.Article) []article.Article {
	var matchingArticles []article.Article

	for _, a := range articles {
		titleStemmed := porterstemmer.StemString(strings.ToLower(string(a.Title)))
		descriptionStemmed := porterstemmer.StemString(strings.ToLower(string(a.Description)))

		if strings.Contains(titleStemmed, keyword) || strings.Contains(descriptionStemmed, keyword) {
			matchingArticles = append(matchingArticles, a)
		}
	}

	if len(matchingArticles) == 0 {
		fmt.Println("No matches found for this keyword.")
	}

	return matchingArticles
}

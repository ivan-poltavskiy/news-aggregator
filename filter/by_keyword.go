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

	for _, article := range articles {

		if matchesConditions(article, keyword) {
			matchingArticles = append(matchingArticles, article)
		}
	}

	if len(matchingArticles) == 0 {
		fmt.Println("No matches found for this keyword.")
	}

	return matchingArticles
}

// matchesConditions checks if an article matches at least one condition.
func matchesConditions(a article.Article, keyword string) bool {
	conditions := []string{
		porterstemmer.StemString(strings.ToLower(string(a.Title))),
		porterstemmer.StemString(strings.ToLower(string(a.Description))),
	}

	for _, condition := range conditions {
		if strings.Contains(condition, keyword) {
			return true
		}
	}
	return false
}

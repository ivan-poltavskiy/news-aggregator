package filter

import (
	"fmt"
	"github.com/reiver/go-porterstemmer"
	"news-aggregator/entity/article"
	"strings"
)

// ByKeyword filters the slice of articles by provided keyword and returns
// the slice of matching articles.
type ByKeyword struct {
	Keywords []string
}

// Filter filters the incoming collector list from different sources by keywords.
func (keywordFilter ByKeyword) Filter(newsArticles []article.Article) []article.Article {
	var matchingNewsArticles []article.Article
	for _, keyword := range keywordFilter.Keywords {
		stemmedKeyword := porterstemmer.StemString(strings.ToLower(keyword))
		matchingNewsArticles = append(matchingNewsArticles, filterNewsByKeyword(stemmedKeyword, newsArticles)...)
	}
	return matchingNewsArticles
}

// filterNewsByKeyword filters the incoming collector list by keyword and returns the filtered list.
func filterNewsByKeyword(keyword string, newsArticles []article.Article) []article.Article {
	var matchingNewsArticles []article.Article

	for _, article := range newsArticles {

		if matchesConditions(article, keyword) {
			matchingNewsArticles = append(matchingNewsArticles, article)
		}
	}

	if len(matchingNewsArticles) == 0 {
		fmt.Println("No matches found for this keyword.")
	}

	return matchingNewsArticles
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

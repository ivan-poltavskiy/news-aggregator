package filter

import (
	"fmt"
	"github.com/reiver/go-porterstemmer"
	"news-aggregator/entity/news"
	"strings"
)

// ByKeyword filters the slice of news by provided keyword and returns
// the slice of matching news.
type ByKeyword struct {
	Keywords []string
}

// Filter filters the incoming news from different sources by keywords.
func (keywordFilter ByKeyword) Filter(newsArticles []news.News) []news.News {
	var matchingNews []news.News
	for _, keyword := range keywordFilter.Keywords {
		stemmedKeyword := porterstemmer.StemString(strings.ToLower(keyword))
		matchingNews = append(matchingNews, filterNewsByKeyword(stemmedKeyword, newsArticles)...)
	}
	return matchingNews
}

// filterNewsByKeyword filters the incoming news by keyword and returns the filtered list.
func filterNewsByKeyword(keyword string, newsArticles []news.News) []news.News {
	var matchingNews []news.News

	for _, article := range newsArticles {

		if matchesConditions(article, keyword) {
			matchingNews = append(matchingNews, article)
		}
	}

	if len(matchingNews) == 0 {
		fmt.Println("No matches found for this keyword.")
	}

	return matchingNews
}

// matchesConditions checks if news matches at least one condition.
func matchesConditions(a news.News, keyword string) bool {
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

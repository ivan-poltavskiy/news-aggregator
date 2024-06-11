package filter

import (
	"NewsAggregator/entity/article"
	"fmt"
	"time"
)

// ByDate filters the slice of articles by a provided date range and returns
// the slice of matching articles.
type ByDate struct {
	StartDate time.Time
	EndDate   time.Time
}

// Filter filters the incoming list of articles by the date range specified in the ByDate struct.
func (f ByDate) Filter(articles []article.Article) []article.Article {
	var matchingArticles []article.Article
	for _, article := range articles {
		if article.Date.After(f.StartDate) && article.Date.Before(f.EndDate) {
			matchingArticles = append(matchingArticles, article)
		}
	}
	if len(matchingArticles) == 0 {
		fmt.Println("No articles were found in this time period.")
	}
	return matchingArticles
}

package filter

import (
	"fmt"
	"news-aggregator/entity/article"
	"time"
)

// ByDate filters the slice of articles by a provided date range and returns
// the slice of matching articles.
type ByDate struct {
	StartDate time.Time
	EndDate   time.Time
}

// Filter filters the incoming list of articles by the date range specified in the ByDate struct.
func (dateFilter ByDate) Filter(articles []article.Article) []article.Article {
	var matchingNewsArticles []article.Article
	for _, a := range articles {
		if a.Date.After(dateFilter.StartDate) && a.Date.Before(dateFilter.EndDate) {
			matchingNewsArticles = append(matchingNewsArticles, a)
		}
	}
	if len(matchingNewsArticles) == 0 {
		fmt.Println("No articles were found in this time period.")
	}
	return matchingNewsArticles
}

package filter_service

import (
	. "NewsAggregator/entity/article"
	"fmt"
	"time"
)

type ByDate struct {
	StartDate time.Time
	EndDate   time.Time
}

func (f ByDate) Filter(articles []Article) []Article {
	var matchingArticles []Article
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

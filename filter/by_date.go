package filter

import (
	"fmt"
	"news-aggregator/entity/news"
	"time"
)

// ByDate filters the slice of news by a provided date range and returns
// the slice of matching news.
type ByDate struct {
	StartDate time.Time
	EndDate   time.Time
}

// Filter filters the incoming list of news by the date range specified in the ByDate struct.
func (dateFilter ByDate) Filter(articles []news.News) []news.News {
	var matchingNews []news.News
	for _, a := range articles {
		if a.Date.After(dateFilter.StartDate) && a.Date.Before(dateFilter.EndDate) {
			matchingNews = append(matchingNews, a)
		}
	}
	if len(matchingNews) == 0 {
		fmt.Println("No articles were found in this time period.")
	}
	return matchingNews
}

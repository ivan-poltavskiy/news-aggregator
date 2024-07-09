package sorter

import (
	"errors"
	"news-aggregator/entity/news"
	"sort"
	"strings"
)

type DateSorter struct {
}

// SortNews sorts news by ASC or DESC
func (DateSorter) SortNews(news []news.News, sortBy string) ([]news.News, error) {

	lowerCaseSortParameter := strings.ToLower(sortBy)
	if lowerCaseSortParameter == "" {
		return news, nil
	}

	var sortingFunctions = map[string]func(i, j int) bool{
		"asc": func(i, j int) bool {
			return news[i].Date.Before(news[j].Date)
		},
		"desc": func(i, j int) bool {
			return news[i].Date.After(news[j].Date)
		},
	}
	if sortingFunctions[lowerCaseSortParameter] != nil {
		sort.Slice(news, sortingFunctions[lowerCaseSortParameter])
		return news, nil
	}

	return nil, errors.New("wrong sorting parameter: " + sortBy)
}

package client

import (
	"errors"
	"github.com/sirupsen/logrus"
	"news-aggregator/entity/article"
	"sort"
	"strings"
)

type DateSorter struct {
}

// SortArticle sorts news by ASC or DESC
func (DateSorter) SortArticle(articles []article.Article, sortBy string) ([]article.Article, error) {

	logrus.Info("DateSorter: sorting articles by " + sortBy)
	lowerCaseSortParameter := strings.ToLower(sortBy)
	if lowerCaseSortParameter == "" {
		return articles, nil
	}

	var sortingFunctions = map[string]func(i, j int) bool{
		"asc": func(i, j int) bool {
			return articles[i].Date.Before(articles[j].Date)
		},
		"desc": func(i, j int) bool {
			return articles[i].Date.After(articles[j].Date)
		},
	}
	if sortingFunctions[lowerCaseSortParameter] != nil {
		sort.Slice(articles, sortingFunctions[lowerCaseSortParameter])
		return articles, nil
	}

	return nil, errors.New("wrong sorting parameter: " + sortBy)
}

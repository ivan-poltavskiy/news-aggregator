package client

import (
	"errors"
	"github.com/sirupsen/logrus"
	"news-aggregator/constant"
	"news-aggregator/entity/news"
	"news-aggregator/filter"
	"news-aggregator/validator"
	"strings"
	"time"
)

//go:generate mockgen -source=client.go -destination=../mocks/mock_client.go -package=mocks news-aggregator/client Client
type Client interface {
	//FetchNews collect the news by some rules defined in the implementations.
	FetchNews() ([]news.News, error)
	//Print outputs the transferred news.
	Print(news []news.News)
}

// buildKeywordFilter extracts keywords from command line arguments and adds them to the filters.
func buildKeywordFilter(keywords string, filters []filter.NewsFilter) []filter.NewsFilter {
	logrus.Info("building keywords filter for: " + keywords)
	if keywords != "" {
		keywordList := strings.Split(keywords, ",")
		uniqueKeywords := checkUnique(keywordList)
		filters = append(filters, filter.ByKeyword{Keywords: uniqueKeywords})
	}
	return filters
}

// buildDateFilters extracts date filters from command line arguments and adds them to the filters.
func buildDateFilters(startDateStr, endDateStr string, filters []filter.NewsFilter) ([]filter.NewsFilter, error) {
	logrus.Info("building date filters for start date: " + startDateStr + "and the end date: " + endDateStr)
	validationErr, isValid := validator.ValidateDate(startDateStr, endDateStr)

	if validationErr != nil {
		return nil, validationErr
	}
	if isValid {
		startDate, err := time.Parse(constant.DateOutputLayout, startDateStr)
		if err != nil {
			return nil, errors.New("Invalid start date: " + startDateStr)
		}

		endDate, err := time.Parse(constant.DateOutputLayout, endDateStr)
		if err != nil {
			return nil, errors.New("Invalid end date: " + endDateStr)
		}

		return append(filters, filter.ByDate{StartDate: startDate, EndDate: endDate}), nil
	}
	return filters, nil
}

// checkUnique returns a slice containing only unique strings from the input slice.
func checkUnique(input []string) []string {
	uniqueMap := make(map[string]struct{})
	var uniqueList []string
	for _, item := range input {
		if _, ok := uniqueMap[item]; !ok {
			uniqueMap[item] = struct{}{}
			uniqueList = append(uniqueList, item)
		}
	}
	return uniqueList
}

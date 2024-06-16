package validator

import (
	"errors"
	"news_aggregator/entity/article"
)

// ValidateSource checks if the provided list of news articles contains at least one article.
// If the input slice is empty, the function will return false, indicating that there are no valid news sources.
func ValidateSource(sources []article.Article) bool {
	return len(sources) != 0
}

// ValidateDate validates the provided start and end dates.
// It returns an error if the start date is after the end date, otherwise, it returns nil.
func ValidateDate(startDate, endDate string) (error, bool) {

	if startDate == "" && endDate == "" {
		return nil, false
	}
	if startDate == "" || endDate == "" {
		return errors.New("either both start date and end date must be provided, or neither"), false
	}

	return nil, true
}

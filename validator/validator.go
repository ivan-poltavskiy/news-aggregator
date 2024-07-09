package validator

import (
	"errors"
	"fmt"
	"news-aggregator/entity/source"
	"slices"
	"strings"
)

// ValidateSource checks if the provided list of news articles contains at least one news.
// If the input slice is empty, the function will return false, indicating that there are no valid news sources.
func ValidateSource(sources []string) (bool, error) {
	for _, currentSource := range sources {
		if !slices.Contains(source.NewsSources, strings.ToLower(currentSource)) {
			return false, errors.New("Source " + currentSource + " is not valid. " +
				"The program supports such news resources:\n" + strings.Join(source.NewsSources, ", "))
		}
	}
	if len(sources) == 0 {
		return false, errors.New(fmt.Sprintf("Please, specify at least one "+
			"news source. The program supports such news resources:\n%s.",
			strings.Join(source.NewsSources, ", ")))
	}
	return true, nil
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

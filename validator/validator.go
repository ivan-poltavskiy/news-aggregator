package validator

import (
	"fmt"
	"news_aggregator/entity/article"
	"time"
)

// ValidateSource validates the provided list of news articles.
// It returns a message prompting the user to specify at least one correct news source
// if the input slice is empty. The supported news sources are listed in the message.
func ValidateSource(sources []article.Article) string {
	if len(sources) == 0 {
		return "Please, specify at least one news source. " +
			"The program supports such news resources:\nABC, BBC, NBC, USA " +
			"Today and Washington Times."
	}
	return ""
}

// ValidateDate validates the provided start and end date strings.
// It prints an error message and returns false if either date string is empty
// or if there is an error parsing the dates.
// It returns the parsed start and end dates if validation is successful.
func ValidateDate(startDateStr, endDateStr string) (bool, time.Time, time.Time) {
	if startDateStr == "" || endDateStr == "" {
		fmt.Println("Please specify both start date and end date or omit them." +
			"Date format - yyyy-mm-dd")
		return false, time.Time{}, time.Time{}
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		fmt.Println("Error parsing start date:", err)
		return false, time.Time{}, time.Time{}
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		fmt.Println("Error parsing end date:", err)
		return false, time.Time{}, time.Time{}
	}

	if startDate.After(endDate) {
		return false, time.Time{}, time.Time{}
	}

	return true, startDate, endDate
}

// CheckUnique returns a slice containing only unique strings from the input slice.
func CheckUnique(input []string) []string {
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

package validator

import (
	"errors"
	"fmt"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
	"strings"

	"github.com/sirupsen/logrus"
	"slices"
)

// ValidateSource checks if the provided list of news articles contains at least one article.
// If the input slice is empty, the function will return false, indicating that there are no valid news sources.
func ValidateSource(sources []string) (bool, error) {
	logrus.Info("Validator: Starting source validation for sources:", sources)

	storage, err := source.LoadExistingSourcesFromStorage(constant.PathToStorage)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"path":  constant.PathToStorage,
			"error": err,
		}).Error("Validator: Failed to load existing sources from storage")
		return false, err
	}

	var storageNames []string
	for _, s := range storage {
		storageNames = append(storageNames, string(s.Name))
	}

	for _, currentSource := range sources {
		if !slices.Contains(storageNames, strings.ToLower(currentSource)) {
			errMessage := fmt.Sprintf("Source %s is not valid. The program supports such news resources:\n%s",
				currentSource, strings.Join(storageNames, ", "))
			logrus.WithFields(logrus.Fields{
				"current_source": currentSource,
				"valid_sources":  storageNames,
			}).Error("Validator: Invalid source")
			return false, errors.New(errMessage)
		}
	}
	if len(sources) == 0 {
		errMessage := fmt.Sprintf("Please, specify at least one news source. The program supports such news resources:\n%s.",
			strings.Join(storageNames, ", "))
		logrus.WithFields(logrus.Fields{
			"valid_sources": storageNames,
		}).Error("Validator: No sources specified")
		return false, errors.New(errMessage)
	}

	logrus.Info("Validator: Source validation successful:", sources)
	return true, nil
}

// ValidateDate validates the provided start and end dates.
// It returns an error if the start date is after the end date, otherwise, it returns nil.
func ValidateDate(startDate, endDate string) (error, bool) {
	logrus.WithFields(logrus.Fields{
		"start_date": startDate,
		"end_date":   endDate,
	}).Info("Validator: Starting date validation")

	if startDate == "" && endDate == "" {
		logrus.Info("Validator: No dates provided")
		return nil, false
	}
	if startDate == "" || endDate == "" {
		errMessage := "either both start date and end date must be provided, or neither"
		logrus.WithFields(logrus.Fields{
			"start_date": startDate,
			"end_date":   endDate,
		}).Error("Validator:", errMessage)
		return errors.New(errMessage), false
	}

	logrus.Info("Validator: Date validation successful:", startDate, endDate)
	return nil, true
}

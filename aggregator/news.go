package aggregator

import (
	"github.com/sirupsen/logrus"
	"news-aggregator/client"
	"news-aggregator/collector"
	"news-aggregator/entity/article"
	"news-aggregator/entity/source"
	"news-aggregator/filter"
	"news-aggregator/validator"
)

// News provides methods for aggregating articles from various sources.
type News struct {
	articleCollector collector.ArticleCollector
}

func New(articleCollector *collector.ArticleCollector) client.Aggregator {
	logrus.Info("News Aggregator Initialized")
	news := &News{articleCollector: *articleCollector}
	return news
}

// Aggregate fetches articles from the provided sources, applies the given
// filters, and returns the filtered articles.
// Parameters:
// - sources: a slice of strings representing the names of the sources to fetch articles from.
// - filters: a variadic parameter of filter.Service to apply filters to the fetched articles.
//
// Returns:
// - A slice of articles that have been fetched and filtered.
// - An error message string if any errors occurred during the process.
func (aggregator *News) Aggregate(sources []string, filters ...filter.ArticleFilter) ([]article.Article, error) {
	logrus.Info("News Aggregator: Starting article aggregation")

	var sourceNames []source.Name
	for _, name := range sources {
		sourceNames = append(sourceNames, source.Name(name))
	}
	logrus.WithFields(logrus.Fields{
		"sources": sources,
	}).Info("News Aggregator: Source names prepared for validation")

	validateSource, err := validator.ValidateSource(sources)
	if !validateSource {
		logrus.WithFields(logrus.Fields{
			"sources": sources,
			"error":   err,
		}).Error("News Aggregator: Source validation failed for sources: ", sources)
		return nil, err
	}
	logrus.Info("News Aggregator: Source validation successful for sources: ", sources)

	articles, err := aggregator.articleCollector.FindNewsByResourcesName(sourceNames)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"sources": sourceNames,
			"error":   err,
		}).Error("News Aggregator: Error finding news by resource names")
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"articles_found": len(articles),
	}).Info("News Aggregator: Articles successfully found")

	for _, f := range filters {
		articles = f.Filter(articles)
	}
	logrus.WithFields(logrus.Fields{
		"remaining_articles": len(articles),
	}).Info("News Aggregator: Filter applied:", filters)

	logrus.WithFields(logrus.Fields{
		"total_articles": len(articles),
	}).Info("News Aggregator: Article aggregation completed successfully")
	return articles, nil
}

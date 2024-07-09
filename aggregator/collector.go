package aggregator

import (
	"news-aggregator/entity/article"
	"news-aggregator/entity/source"
)

// Collector is using for fetching news from source by source name.
type Collector interface {
	FindNewsByResourcesName(sourcesNames []source.Name) ([]article.Article, error)
}

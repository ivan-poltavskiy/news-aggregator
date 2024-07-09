package aggregator

import (
	"news-aggregator/entity/article"
	"news-aggregator/entity/source"
)

type Collector interface {
	FindNewsByResourcesName(sourcesNames []source.Name) ([]article.Article, error)
}

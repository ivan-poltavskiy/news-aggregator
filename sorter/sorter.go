package sorter

import "news-aggregator/entity/article"

// Sorter sorts input data according to specified rules.
type Sorter interface {
	//SortArticle allocates input articles.
	SortArticle(articles []article.Article, sortBy string) ([]article.Article, error)
}

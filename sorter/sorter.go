package sorter

import "news_aggregator/entity/article"

// Sorter sorts input data according to specified rules.
type Sorter interface {
	//SortArticle allocates input articles.
	SortArticle(article []article.Article) []article.Article
}

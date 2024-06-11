// Package collector provides functionality for gathering news from specific sources.
// FindByResourcesName(sourcesNames []source.Name) ([]article.Article, string)
// is used to receive all news from sources passed to it, if these sources are
// correct and present in the system
// findForCurrentSource(currentSource source.Source,
//
//	name source.Name, allArticles []article.Article) []article.Article returns
//	the list of collector from the passed source.
//
// InitializeSource(sources []source.Source) initializes the resources that
// will be available for parsing.
package collector

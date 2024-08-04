// Package collector provides functionality for gathering newsCollector from specific sources.
// FindNewsByResourcesName(sourcesNames []source.Name) ([]newsCollector.News, string)
// is used to receive all newsCollector from sources passed to it, if these sources are
// correct and present in the system
// findNewsForCurrentSource(currentSource source.Source,
//
//	name source.Name, allArticles []newsCollector.News) []newsCollector.News returns
//	the list of collector from the passed source.
//
// InitializeSource(sources []source.Source) initializes the news that
// will be available for parsing.
package collector

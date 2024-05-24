package source

import "NewsAggregator/entity"

// Source is the set of information about source of article.
type Source struct {
	Name       Name
	PathToFile PathToFile
	SourceType Type
	Id         entity.Id
}

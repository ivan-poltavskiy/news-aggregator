package article

import (
	"news-aggregator/entity/source"
	"time"
)

// Article is the set of information about articles in the system.
type Article struct {
	Title       Title       `json:"title"`
	Description Description `json:"description"`
	Link        Link        `json:"url"`
	Date        time.Time   `json:"publishedAt"`
	SourceName  source.Name
}

// Description provides brief information about the article.
type Description string

func (d Description) String() string {
	return string(d)
}

// Link contains the url of the article.
type Link string

type Title string

func (t Title) String() string {
	return string(t)
}

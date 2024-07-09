package news

import (
	"news-aggregator/entity/source"
	"time"
)

// News is the set of information about news articles in the system.
type News struct {
	Title       Title       `json:"title"`
	Description Description `json:"description"`
	Link        Link        `json:"url"`
	Date        time.Time   `json:"publishedAt"`
	SourceName  source.Name
}

// Description provides brief information about the news.
type Description string

func (d Description) String() string {
	return string(d)
}

// Link contains the url of the news.
type Link string

type Title string

func (t Title) String() string {
	return string(t)
}

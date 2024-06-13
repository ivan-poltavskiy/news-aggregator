package article

import (
	"news_aggregator/entity/source"
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

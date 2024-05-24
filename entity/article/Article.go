package article

import (
	"NewsAggregator/entity"
	"time"
)

// Article is the set of information about articles in the system.
type Article struct {
	Id          entity.Id
	Title       Title       `json:"title"`
	Description Description `json:"description"`
	Link        Link        `json:"url"`
	Date        time.Time   `json:"publishedAt"`
}

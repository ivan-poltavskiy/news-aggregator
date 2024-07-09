package parser

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	"os"
)

// Rss analyzes RSS sources.
type Rss struct {
}

// Parse reads and parses a XML (RSS) file specified by the path and returns a slice of articles.
func (rss Rss) Parse(path source.PathToFile, name source.Name) ([]news.News, error) {

	parser := gofeed.NewParser()
	filename := fmt.Sprintf(string(path))
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = fmt.Errorf("cannot close RSS file: %w", cerr)
		}
	}()

	feed, err := parser.Parse(file)
	if err != nil {
		return nil, err
	}

	var newsData []news.News
	for _, item := range feed.Items {
		newsData = append(newsData, news.News{
			Title:       news.Title(item.Title),
			Description: news.Description(item.Description),
			Link:        news.Link(item.Link),
			Date:        *item.PublishedParsed,
			SourceName:  name,
		})
	}
	return newsData, nil
}

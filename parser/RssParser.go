package parser

import (
	. "NewsAggregator/entity"
	. "NewsAggregator/entity/article"
	"NewsAggregator/entity/source"
	"fmt"
	"github.com/mmcdole/gofeed"
	"os"
)

// RssParser for analyze RSS sources.
type RssParser struct {
}

func (rssParser RssParser) ParseSource(path source.PathToFile) []Article {

	parser := gofeed.NewParser()
	filename := fmt.Sprintf(string(path))
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Cannot open the file:", err)
		return nil

	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	feed, err := parser.Parse(file)
	if err != nil {
		fmt.Println("Error with parse RSS content:", err)
		return nil
	}

	var articles []Article
	for i, item := range feed.Items {
		articles = append(articles, Article{
			Id:          Id(i + 1),
			Title:       Title(item.Title),
			Description: Description(item.Description),
			Link:        Link(item.Link),
			Date:        *item.PublishedParsed,
		})
	}
	return articles
}

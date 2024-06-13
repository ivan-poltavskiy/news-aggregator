package parser

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"news_aggregator/entity/article"
	"news_aggregator/entity/source"
	"os"
)

// Rss analyzes RSS sources.
type Rss struct {
}

// ParseSource reads and parses a XML (RSS) file specified by the path and returns a slice of articles.
func (rss Rss) ParseSource(path source.PathToFile) []article.Article {

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

	var articles []article.Article
	for _, item := range feed.Items {
		articles = append(articles, article.Article{
			Title:       article.Title(item.Title),
			Description: article.Description(item.Description),
			Link:        article.Link(item.Link),
			Date:        *item.PublishedParsed,
		})
	}
	return articles
}

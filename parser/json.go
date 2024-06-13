package parser

import (
	"encoding/json"
	"fmt"
	"news_aggregator/entity/article"
	"news_aggregator/entity/source"
	"os"
)

// Json analyzes JSON sources.
type Json struct {
}

// ParseSource reads and parses a JSON file specified by the path and returns a slice of articles.
func (jsonFile Json) ParseSource(path source.PathToFile) []article.Article {
	filename := fmt.Sprintf(string(path))

	byteValue, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error with parse JSON content:", err)
		return nil
	}

	var articles struct {
		Articles []article.Article `json:"articles"`
	}

	err = json.Unmarshal(byteValue, &articles)
	if err != nil {
		fmt.Println("Error unmarshalling JSON content:", err)
		return nil
	}

	return articles.Articles
}

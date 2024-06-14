package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"news_aggregator/entity/article"
	"news_aggregator/entity/source"
	"os"
)

// Json analyzes JSON sources.
type Json struct {
}

// ParseSource reads and parses a JSON file specified by the path and returns a slice of articles.
func (jsonFile Json) ParseSource(path source.PathToFile, name source.Name) ([]article.Article, error) {
	filename := fmt.Sprintf(string(path))

	byteValue, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.New("Error with parse JSON content: " + err.Error())
	}

	var articles struct {
		Articles []article.Article `json:"articles"`
	}

	err = json.Unmarshal(byteValue, &articles)
	if err != nil {
		return nil, errors.New("Error with parse JSON content: " + err.Error())
	}

	for i := range articles.Articles {
		articles.Articles[i].SourceName = name
	}

	return articles.Articles, nil
}

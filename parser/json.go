package parser

import (
	"encoding/json"
	"errors"
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	"os"
)

// Json analyzes JSON sources.
type Json struct {
}

// Parse reads and parses a JSON file specified by the path and returns a slice of news.
func (jsonFile Json) Parse(path source.PathToFile, name source.Name) ([]news.News, error) {

	newsContent, err := os.ReadFile(string(path))
	if err != nil {
		return nil, errors.New("Error with parse JSON content: " + err.Error())
	}

	var newsData struct {
		News []news.News `json:"articles"`
	}

	err = json.Unmarshal(newsContent, &newsData)
	if err != nil {
		return nil, errors.New("Error with parse JSON content: " + err.Error())
	}

	for i := range newsData.News {
		newsData.News[i].SourceName = name
	}

	return newsData.News, nil
}

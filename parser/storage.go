package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	"os"
)

type Storage struct {
}

func (storage Storage) Parse(path source.PathToFile, name source.Name) ([]news.News, error) {

	if _, err := os.Stat(string(path)); os.IsNotExist(err) {
		return nil, fmt.Errorf("JSON file not found for source: %s", path)
	}

	jsonFile, err := os.Open(string(path))
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	var articles []news.News
	if err := json.Unmarshal(byteValue, &articles); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}

	for i := range articles {
		articles[i].SourceName = name
	}

	return articles, nil
}

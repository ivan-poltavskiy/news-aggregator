package parser

import (
	. "NewsAggregator/entity/article"
	"NewsAggregator/entity/source"
	"encoding/json"
	"fmt"
	"os"
)

// JsonParser for analyze JSON sources.
type JsonParser struct {
}

func (jsonFile JsonParser) ParseSource(path source.PathToFile) []Article {
	filename := fmt.Sprintf(string(path))

	byteValue, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error with parse JSON content:", err)
		return nil
	}

	var articles struct {
		Articles []Article `json:"articles"`
	}

	err = json.Unmarshal(byteValue, &articles)
	if err != nil {
		fmt.Println("Error unmarshalling JSON content:", err)
		return nil
	}

	return articles.Articles
}

package source

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// PathToFile describes the path to a specific source in the system
type PathToFile string

// Name describes the name of source
type Name string

// Type describes of the type of document for a particular source
type Type string

// Source is the set of information about source of article.
type Source struct {
	Name       Name
	PathToFile PathToFile
	SourceType Type
}

// Stores all types of sources provided.
const (
	RSS      Type = "RSS"
	JSON     Type = "JSON"
	UsaToday Type = "UsaToday"
	STORAGE  Type = "STORAGE"
)

// LoadExistingSourcesFromStorage loads sources from a JSON file
func LoadExistingSourcesFromStorage(filename string) ([]Source, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	value, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var sources []Source
	err = json.Unmarshal(value, &sources)
	if err != nil {
		return nil, err
	}

	return sources, nil
}

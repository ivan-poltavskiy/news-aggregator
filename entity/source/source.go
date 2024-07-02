package source

import (
	"bufio"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
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
	logrus.Info("Source: Starting loading the existing sources from storage")
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.Error("Source: Error closing file: ", err)
		}
	}(file)

	reader := bufio.NewReader(file)
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var sources []Source
	err = json.Unmarshal(content, &sources)
	if err != nil {
		return nil, err
	}

	return sources, nil
}

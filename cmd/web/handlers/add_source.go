package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
	"os"
	"path/filepath"
	"time"
)

// addSourceRequest is a structure containing the fields required to add a new source.
type addSourceRequest struct {
	URL string `json:"url"`
}

// AddSourceHandler is a handler for add the new source to the storage.
func AddSourceHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var requestBody addSourceRequest

	err = json.Unmarshal(body, &requestBody)
	if err != nil || requestBody.URL == "" {
		http.Error(w, "Invalid request body or URL parameter is missing", http.StatusBadRequest)
		return
	}

	url := requestBody.URL

	err, rssURL := GetRssFeedLink(w, url)
	if rssURL == "" {
		return
	}

	rssResp, err := http.Get(rssURL)
	if err != nil || rssResp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to download RSS feed", http.StatusInternalServerError)
		return
	}
	defer rssResp.Body.Close()

	currentDate := time.Now().Format(constant.DateOutputLayout)

	directoryPath := filepath.Join("resources", currentDate)
	if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		return
	}

	fileName := ExtractDomainName(url) + ".xml"

	filePath := filepath.Join(directoryPath, fileName)
	outputFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, rssResp.Body)
	if err != nil {
		http.Error(w, "Failed to save RSS feed", http.StatusInternalServerError)
		return
	}

	sourceEntity := source.Source{
		Name:       source.Name(ExtractDomainName(url)),
		PathToFile: source.PathToFile(filePath),
		SourceType: source.RSS,
	}

	if !SourceExists(sourceEntity.Name) {
		AddSourceToFile(sourceEntity)
		fmt.Fprintf(w, "RSS feed downloaded and saved to %s and source added", filePath)
	} else {
		fmt.Fprintf(w, "RSS feed downloaded and saved to %s but source already exists", filePath)
	}
}

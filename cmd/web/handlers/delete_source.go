package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"news-aggregator/cmd/web"
	"news-aggregator/entity/source"
	"os"
	"strings"
)

type deleteSourceRequest struct {
	Name string `json:"name"`
}

func DeleteSourceByNameHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var request deleteSourceRequest
	err = json.Unmarshal(body, &request)
	if err != nil || request.Name == "" {
		http.Error(w, "Invalid request body or name parameter is missing", http.StatusBadRequest)
		return
	}

	sources := web.ReadSourcesFromFile()
	var updatedSources []source.Source
	found := false
	for _, currentSource := range sources {
		if strings.ToLower(string(currentSource.Name)) != strings.ToLower(request.Name) {
			updatedSources = append(updatedSources, currentSource)
		} else {
			found = true
			// Удаление файла
			err := os.Remove(string(currentSource.PathToFile))
			if err != nil {
				http.Error(w, "Failed to delete source file", http.StatusInternalServerError)
				return
			}
		}
	}

	if !found {
		http.Error(w, "Source not found", http.StatusNotFound)
		return
	}

	err = web.WriteSourcesToFile(updatedSources)
	if err != nil {
		http.Error(w, "Failed to write updated sources to file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Source deleted successfully"))
}

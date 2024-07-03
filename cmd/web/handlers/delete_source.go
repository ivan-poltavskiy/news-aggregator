package handlers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"news-aggregator/entity/source"
	"os"
	"path/filepath"
	"strings"
)

type deleteSourceRequest struct {
	Name string `json:"name"`
}

// DeleteSourceByNameHandler is a handler for removing the source from the storage.
func DeleteSourceByNameHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Error("DeleteSourceByNameHandler: Failed to read request body ", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Error("DeleteSourceByNameHandler: Failed to close request body ", err)
		}
	}(r.Body)

	var request deleteSourceRequest
	err = json.Unmarshal(body, &request)
	if err != nil || request.Name == "" {
		logrus.Error("DeleteSourceByNameHandler: Invalid request body or name parameter is missing")
		http.Error(w, "Invalid request body or name parameter is missing", http.StatusBadRequest)
		return
	}
	logrus.Info("DeleteSourceByNameHandler: The name from the request to delete the source was successfully retrieved: ", request.Name)

	var updatedSources []source.Source
	found := false
	for _, currentSource := range ReadSourcesFromFile() {
		if strings.ToLower(string(currentSource.Name)) != strings.ToLower(request.Name) {
			updatedSources = append(updatedSources, currentSource)
		} else {
			found = true
			directoryPath := filepath.Join("resources", strings.ToLower(request.Name))
			err := os.RemoveAll(directoryPath)
			if err != nil {
				logrus.Error("DeleteSourceByNameHandler: Failed to delete source directory ", err)
				http.Error(w, "Failed to delete source directory", http.StatusInternalServerError)
				return
			}
			logrus.Info("DeleteSourceByNameHandler: The source directory was successfully deleted: ", directoryPath)
		}
	}

	if !found {
		logrus.Warn("DeleteSourceByNameHandler: Source not found: ", request.Name)
		http.Error(w, "Source not found", http.StatusNotFound)
		return
	}

	err = WriteSourcesToFile(updatedSources)
	if err != nil {
		logrus.Error("DeleteSourceByNameHandler: Failed to write updated sources to file ", err)
		http.Error(w, "Failed to write updated sources to file", http.StatusInternalServerError)
		return
	}
	logrus.Info("DeleteSourceByNameHandler: The updated sources were successfully written to the file")

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Source deleted successfully"))
	if err != nil {
		logrus.Error("DeleteSourceByNameHandler: Failed to write response for delete source ", err)
		return
	}
	logrus.Info("DeleteSourceByNameHandler: Response for delete source was successfully written")
}

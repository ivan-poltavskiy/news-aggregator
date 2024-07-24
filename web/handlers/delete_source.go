package handlers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"news-aggregator/storage"
	sourceService "news-aggregator/web/source"
	"strings"
)

type deleteSourceRequest struct {
	Name string `json:"name"`
}

// DeleteSourceByNameHandler is a handler for removing the source from the storage.
func DeleteSourceByNameHandler(w http.ResponseWriter, r *http.Request, sourceStorage storage.Storage) {

	var request deleteSourceRequest
	body, err := io.ReadAll(r.Body)

	if err != nil {
		logrus.Error("Failed to read request body: ", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Error("Failed to close request body: ", err)
		}
	}(r.Body)

	err = json.Unmarshal(body, &request)
	if err != nil || request.Name == "" {
		logrus.Error("Invalid request body or name parameter is missing")
		http.Error(w, "Invalid request body or name parameter is missing", http.StatusBadRequest)
		return
	}

	logrus.Infof("Request to delete source received: %s", request.Name)

	service := sourceService.NewSourceService(sourceStorage)
	err = service.DeleteSourceByName(request.Name)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			logrus.Warnf("Source not found: %s", request.Name)
			http.Error(w, "Source not found", http.StatusNotFound)
		} else {
			logrus.Error("Failed to delete source: ", err)
			http.Error(w, "Failed to delete source", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Source deleted successfully"))
	if err != nil {
		logrus.Error("Failed to write response for delete source: ", err)
		return
	}
	logrus.Info("Response for delete source written successfully")
}

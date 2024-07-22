package handlers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"news-aggregator/cmd/web/service"
	"news-aggregator/storage/source"
)

// addSourceRequest is a structure containing the fields required to add a new source.
type addSourceRequest struct {
	URL string `json:"url"`
}

// AddSourceHandler is a handler for adding the new source to the storage.
func AddSourceHandler(w http.ResponseWriter, r *http.Request, storage source.Storage) {
	var requestBody addSourceRequest

	if err := getUrlFromRequest(r, &requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.Info("AddSourceHandler: The URL from the request to add the source was successfully retrieved: ", requestBody.URL)

	sourceName, err := service.SaveSource(requestBody.URL, storage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("News saved successfully. Name of source: " + string(sourceName))); err != nil {
		logrus.Error("Failed to write response: ", err)
	}
}

// get the URL of the source from the request
func getUrlFromRequest(r *http.Request, requestBody *addSourceRequest) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Error("Failed to read request body: ", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Error("Failed to close body: ", err)
		}
	}(r.Body)

	if err := json.Unmarshal(body, requestBody); err != nil || requestBody.URL == "" {
		logrus.Error("Invalid request body or URL parameter is missing")
		return err
	}

	logrus.Info("getUrlFromRequest: Successfully parsed request body")
	return nil
}

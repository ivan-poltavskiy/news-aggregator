package handlers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"news-aggregator/cmd/web/service"
)

// addSourceRequest is a structure containing the fields required to add a new source.
type addSourceRequest struct {
	URL string `json:"url"`
}

// AddSourceHandler is a handler for adding the new source to the storage.
func AddSourceHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody addSourceRequest

	if err := getUrlFromRequest(r, &requestBody); err != nil {
		httpResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.Info("AddSourceHandler: The URL from the request to add the source was successfully retrieved: ", requestBody.URL)

	sourceEntity, err := service.AddSource(requestBody.URL)
	if err != nil {
		httpResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpResponse(w, "Articles saved successfully. Name of source: "+string(sourceEntity.Name), http.StatusOK)
}

// get the URL of the source from the request
func getUrlFromRequest(r *http.Request, requestBody *addSourceRequest) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Error("Failed to read request body: ", err)
		return err
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, requestBody); err != nil || requestBody.URL == "" {
		logrus.Error("Invalid request body or URL parameter is missing")
		return err
	}

	logrus.Info("getUrlFromRequest: Successfully parsed request body")
	return nil
}

// Helper function to write HTTP responses
func httpResponse(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	if _, err := w.Write([]byte(message)); err != nil {
		logrus.Error("Failed to write response: ", err)
	}
}

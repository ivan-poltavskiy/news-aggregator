package source

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"news-aggregator/entity/source"
	"news-aggregator/storage"
	"strings"
)

type deleteSourceRequest struct {
	Name string `json:"name"`
}

type addSourceRequest struct {
	URL string `json:"url"`
}

type SourceHandler struct {
	service *SourcesService
}

func NewSourceHandler(storage storage.Storage) *SourceHandler {
	return &SourceHandler{
		service: NewSourceService(storage),
	}
}

func (h *SourceHandler) DeleteSourceByNameHandler(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.DeleteSourceByName(source.Name(request.Name))
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

func (h *SourceHandler) AddSourceHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody addSourceRequest

	if err := parseRequest(r, &requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logrus.Info("AddSourceHandler: The URL from the request to add the source was successfully retrieved: ", requestBody.URL)

	sourceName, err := h.service.SaveSource(requestBody.URL)
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
func parseRequest(r *http.Request, requestBody *addSourceRequest) error {
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

	logrus.Info("parseRequest: Successfully parsed request body")
	return nil
}

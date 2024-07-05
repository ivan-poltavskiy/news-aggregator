package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"news-aggregator/cmd/web/service"
	"news-aggregator/entity/source"
)

// mock the AddSource function
func mockAddSource(url string) (source.Name, error) {
	if url == "" {
		return "", fmt.Errorf("passed url is empty")
	}
	if url == "https://www.pravda.com.ua/" {
		return "pravda", nil
	}
	return "", fmt.Errorf("unknown error")
}

func TestAddSourceHandler(t *testing.T) {
	patch := monkey.Patch(service.AddSource, mockAddSource)
	defer patch.Unpatch()

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "ValidRequest",
			requestBody:    addSourceRequest{URL: "https://www.pravda.com.ua/"},
			expectedStatus: http.StatusOK,
			expectedBody:   "Articles saved successfully. Name of source: pravda",
		},
		{
			name:           "EmptyURL",
			requestBody:    addSourceRequest{URL: ""},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "passed url is empty",
		},
		{
			name:           "UnknownURL",
			requestBody:    addSourceRequest{URL: "https://unknown.com/"},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "unknown error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, err := http.NewRequest(http.MethodPost, "/source", bytes.NewBuffer(body))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(AddSourceHandler)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}

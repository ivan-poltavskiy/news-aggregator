package handlers

import (
	"bou.ke/monkey"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"news-aggregator/cmd/web/service"
	"news-aggregator/entity/source"
	newsStorage "news-aggregator/storage/news"
	newsStorage_mock "news-aggregator/storage/news/mock_aggregator"
	sourceStorage "news-aggregator/storage/source"
	"news-aggregator/storage/source/mock_aggregator"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mock the Save function
func mockSaveSource(url string, sourceStorage sourceStorage.Storage, newsStorage newsStorage.NewsStorage) (source.Name, error) {
	if url == "" {
		return "", fmt.Errorf("passed url is empty")
	}
	if url == "https://www.pravda.com.ua/" {
		return "pravda", nil
	}
	if sourceStorage == nil && newsStorage == nil {
		return "", fmt.Errorf("passed sourceStorage or newsStorage is empty")
	}
	return "", fmt.Errorf("unknown error")
}

func TestAddSourceHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorage := mock_aggregator.NewMockStorage(ctrl)
	newsMockStorage := newsStorage_mock.NewMockNewsStorage(ctrl)

	patch := monkey.Patch(service.SaveSource, mockSaveSource)
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
			expectedBody:   "News saved successfully. Name of source: pravda",
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
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				AddSourceHandler(w, r, mockStorage, newsMockStorage)
			})

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}

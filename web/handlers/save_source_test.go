package handlers

import (
	"bou.ke/monkey"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"news-aggregator/entity/source"
	"news-aggregator/storage/mock_aggregator"
	sourceService "news-aggregator/web/source"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockSaveSource mocks the SaveSource method
func mockSaveSource(_ *sourceService.SourcesService, url string) (source.Name, error) {
	if url == "" {
		return "", fmt.Errorf("passed url is empty")
	}
	if url == "https://www.pravda.com.ua/" {
		return "pravda", nil
	}
	return "", fmt.Errorf("unknown error")
}

func TestAddSourceHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorage := mock_aggregator.NewMockStorage(ctrl)

	service := sourceService.NewSourceService(mockStorage)

	// Patch the SaveSource method
	patch := monkey.PatchInstanceMethod(reflect.TypeOf(service), "SaveSource", mockSaveSource)
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
				AddSourceHandler(w, r, mockStorage)
			})

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}

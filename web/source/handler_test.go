package source

import (
	"bou.ke/monkey"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"news-aggregator/entity/source"
	"news-aggregator/storage/mock_aggregator"
	"reflect"
	"testing"
)

func TestDeleteSourceByNameHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock_aggregator.NewMockStorage(ctrl)
	handler := NewSourceHandler(mockStorage)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedBody   string
		mockFunc       func()
	}{
		{
			name:           "ValidRequest",
			requestBody:    map[string]string{"name": "ExistingSource"},
			expectedStatus: http.StatusOK,
			expectedBody:   "Source deleted successfully",
			mockFunc: func() {
				mockStorage.EXPECT().DeleteSourceByName(source.Name("ExistingSource")).Return(nil)
			},
		},
		{
			name:           "NonExistingSource",
			requestBody:    map[string]string{"name": "NonExistingSource"},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Source not found",
			mockFunc: func() {
				mockStorage.EXPECT().DeleteSourceByName(source.Name("NonExistingSource")).Return(errors.New("source not found"))
			},
		},
		{
			name:           "InvalidRequestBody",
			requestBody:    "InvalidBody",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request body or name parameter is missing",
			mockFunc:       func() {},
		},
		{
			name:           "EmptyName",
			requestBody:    map[string]string{"name": ""},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request body or name parameter is missing",
			mockFunc:       func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, err := http.NewRequest(http.MethodDelete, "/sources", bytes.NewBuffer(body))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()

			tt.mockFunc()

			handler.DeleteSourceByNameHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}

// mockSaveSource mocks the SaveSource method
func mockSaveSource(_ *Service, url string) (source.Name, error) {
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

	service := NewService(mockStorage)
	patch := monkey.PatchInstanceMethod(reflect.TypeOf(service), "SaveSource", mockSaveSource)
	defer patch.Unpatch()

	handler := NewSourceHandler(mockStorage)

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

			handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.AddSourceHandler(w, r)
			})

			handlerFunc.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}

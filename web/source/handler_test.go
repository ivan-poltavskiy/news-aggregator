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
	storage "news-aggregator/storage/mock_aggregator"
	"reflect"
	"testing"
)

func TestDeleteSourceByNameHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storage.NewMockStorage(ctrl)
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
func mockSaveSource(_ *Service, request AddSourceRequest) (source.Name, error) {
	if request.URL == "" {
		return "", fmt.Errorf("passed url is empty")
	}
	if request.URL == "https://www.pravda.com.ua/" {
		return "pravda", nil
	}
	return "", fmt.Errorf("unknown error")
}

func TestAddSourceHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorage := storage.NewMockStorage(ctrl)

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
			requestBody:    AddSourceRequest{URL: "https://www.pravda.com.ua/"},
			expectedStatus: http.StatusOK,
			expectedBody:   "News saved successfully. Name of source: pravda",
		},
		{
			name:           "EmptyURL",
			requestBody:    AddSourceRequest{URL: ""},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "passed url is empty",
		},
		{
			name:           "UnknownURL",
			requestBody:    AddSourceRequest{URL: "https://unknown.com/"},
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

func TestGetAllSources(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := storage.NewMockStorage(ctrl)
	handler := NewSourceHandler(mockService)

	tests := []struct {
		name           string
		mockFunc       func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "SuccessfulFetch",
			mockFunc: func() {
				mockService.EXPECT().GetSources().Return([]source.Source{
					{Name: "Source1", SourceType: source.STORAGE},
					{Name: "Source2", SourceType: source.STORAGE},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `["Source1","Source2"]`,
		},
		{
			name: "ServiceError",
			mockFunc: func() {
				mockService.EXPECT().GetSources().Return(nil, errors.New("internal server error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			rr := httptest.NewRecorder()

			handler.GetAllSources(rr)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestGetAllSourcesWithWriteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := storage.NewMockStorage(ctrl)
	handler := NewSourceHandler(mockService)

	tests := []struct {
		name           string
		mockFunc       func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "WriteError",
			mockFunc: func() {
				mockService.EXPECT().GetSources().Return([]source.Source{
					{Name: "Source1", SourceType: source.STORAGE},
					{Name: "Source2", SourceType: source.STORAGE},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			rr := httptest.NewRecorder()
			writer := &CustomResponseWriter{ResponseWriter: rr, err: errors.New("write error")}

			handler.GetAllSources(writer)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Empty(t, rr.Body.String())
		})
	}
}

// CustomResponseWriter simulates an error during the Write operation.
type CustomResponseWriter struct {
	http.ResponseWriter
	err error
}

func (c *CustomResponseWriter) Write(b []byte) (int, error) {
	return 0, c.err
}

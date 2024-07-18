package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"news-aggregator/cmd/web/handlers"
	"news-aggregator/storage/mock_aggregator"
)

func TestDeleteSourceByNameHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock_aggregator.NewMockStorage(ctrl)

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
				mockStorage.EXPECT().DeleteSourceByName("ExistingSource").Return(nil)
			},
		},
		{
			name:           "NonExistingSource",
			requestBody:    map[string]string{"name": "NonExistingSource"},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Source not found",
			mockFunc: func() {
				mockStorage.EXPECT().DeleteSourceByName("NonExistingSource").Return(errors.New("source not found"))
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
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlers.DeleteSourceByNameHandler(w, r, mockStorage)
			})

			tt.mockFunc()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}

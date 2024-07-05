package handlers_test

import (
	"bou.ke/monkey"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"news-aggregator/cmd/web/handlers"
	"news-aggregator/cmd/web/service"
)

// mock the DeleteAndUpdateSources function
func mockDeleteAndUpdateSources(name string) error {
	if name == "ExistingSource" {
		return nil
	}
	if name == "NonExistingSource" {
		return fmt.Errorf("source not found")
	}
	return fmt.Errorf("unknown error")
}

func TestDeleteSourceByNameHandler(t *testing.T) {
	patch := monkey.Patch(service.DeleteAndUpdateSources, mockDeleteAndUpdateSources)
	defer patch.Unpatch()

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "ValidRequest",
			requestBody:    map[string]string{"name": "ExistingSource"},
			expectedStatus: http.StatusOK,
			expectedBody:   "Source deleted successfully",
		},
		{
			name:           "NonExistingSource",
			requestBody:    map[string]string{"name": "NonExistingSource"},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Source not found",
		},
		{
			name:           "InvalidRequestBody",
			requestBody:    "InvalidBody",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request body or name parameter is missing",
		},
		{
			name:           "EmptyName",
			requestBody:    map[string]string{"name": ""},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request body or name parameter is missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, err := http.NewRequest(http.MethodDelete, "/sources", bytes.NewBuffer(body))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlers.DeleteSourceByNameHandler)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}

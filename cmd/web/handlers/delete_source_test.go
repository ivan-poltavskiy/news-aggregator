package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
	"os"
	"path/filepath"
	"testing"
)

type testDeleteSourceRequest struct {
	Name string `json:"name"`
}

func setupTestEnvironment(t *testing.T) {
	// Set up a test environment
	originalPath := constant.PathToStorage
	originalPathToResources := constant.PathToResources
	constant.PathToResources = "../../../resources/testdata/handlers"
	constant.PathToStorage = "../../../resources/testdata/handlers/sources.json"
	t.Cleanup(func() {
		constant.PathToResources = originalPathToResources
		constant.PathToStorage = originalPath
	})

	// Create directories and files in the 'resources' directory
	err := os.MkdirAll("../../../resources/testdata/handlers", os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create resources directory: %v", err)
	}

	// Create sample sources.json
	sources := []source.Source{
		{Name: "TestSource1"},
		{Name: "TestSource2"},
	}
	data, err := json.Marshal(sources)
	if err != nil {
		t.Fatalf("Failed to marshal sources: %v", err)
	}
	err = os.WriteFile("../../../resources/testdata/handlers/sources.json", data, os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to write test sources.json: %v", err)
	}

	// Create sample source directory
	err = os.MkdirAll(filepath.Join("../../../resources/testdata/handlers", "testsource1"), os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create test source directory: %v", err)
	}
}

func cleanupTestEnvironment(t *testing.T) {
	// Remove the created directories and files after test
	err := os.RemoveAll("../../../resources/testdata/handlers")
	if err != nil {
		t.Fatalf("Failed to clean up resources directory: %v", err)
	}
}

func TestDeleteSourceByNameHandler(t *testing.T) {
	tests := []struct {
		name         string
		setupFunc    func(t *testing.T) (*http.Request, *httptest.ResponseRecorder)
		validateFunc func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		{
			name: "Delete existing source",
			setupFunc: func(t *testing.T) (*http.Request, *httptest.ResponseRecorder) {
				setupTestEnvironment(t)
				requestBody, _ := json.Marshal(testDeleteSourceRequest{Name: "TestSource1"})
				req, err := http.NewRequest("DELETE", "/sources", bytes.NewBuffer(requestBody))
				if err != nil {
					t.Fatalf("Failed to create request: %v", err)
				}
				rr := httptest.NewRecorder()
				return req, rr
			},
			validateFunc: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)

				// Validate the storage
				data, err := os.ReadFile(constant.PathToStorage)
				if err != nil {
					t.Fatalf("Failed to read updated sources.json: %v", err)
				}
				var updatedSources []source.Source
				err = json.Unmarshal(data, &updatedSources)
				if err != nil {
					t.Fatalf("Failed to unmarshal updated sources.json: %v", err)
				}
				assert.Equal(t, 1, len(updatedSources))
				assert.Equal(t, "TestSource2", string(updatedSources[0].Name))

				// Validate the directory removal
				sourceDir := filepath.Join("../../../resources/testdata/handlers", "testsource1")
				if _, err := os.Stat(sourceDir); !os.IsNotExist(err) {
					t.Errorf("Handler did not remove the source directory")
				}
			},
		},
		{
			name: "Delete non existing source",
			setupFunc: func(t *testing.T) (*http.Request, *httptest.ResponseRecorder) {
				setupTestEnvironment(t)
				requestBody, _ := json.Marshal(testDeleteSourceRequest{Name: "TestSource111"})
				req, err := http.NewRequest("DELETE", "/sources", bytes.NewBuffer(requestBody))
				if err != nil {
					t.Fatalf("Failed to create request: %v", err)
				}
				rr := httptest.NewRecorder()
				return req, rr
			},
			validateFunc: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, rr.Code, "Handler returned wrong status code: got %v want %v", rr.Code, http.StatusNotFound)

				assert.Equal(t, "Source not found\n", rr.Body.String(), "Handler returned unexpected body: got %v want %v", rr.Body.String(), "Source not found\n")

				// Validate the storage
				data, err := os.ReadFile(constant.PathToStorage)
				if err != nil {
					t.Fatalf("Failed to read updated sources.json: %v", err)
				}
				var updatedSources []source.Source
				err = json.Unmarshal(data, &updatedSources)
				if err != nil {
					t.Fatalf("Failed to unmarshal updated sources.json: %v", err)
				}
				assert.Equal(t, 2, len(updatedSources))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, rr := tt.setupFunc(t)
			defer cleanupTestEnvironment(t)
			DeleteSourceByNameHandler(rr, req)
			tt.validateFunc(t, rr)
		})
	}
}

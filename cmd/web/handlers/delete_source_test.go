package handlers

import (
	"bytes"
	"encoding/json"
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
	err = os.WriteFile("../../../resources/testdata/handlers/sources.json", data, 0644)
	if err != nil {
		t.Fatalf("Failed to write test sources.json: %v", err)
	}

	// Create sample source directory
	err = os.MkdirAll(filepath.Join("../../../resources/testdata/handlers", "testsource1"), 0755)
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
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	// Prepare the request
	requestBody, _ := json.Marshal(testDeleteSourceRequest{Name: "TestSource1"})
	req, err := http.NewRequest("DELETE", "/sources", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()

	// Call the handler
	DeleteSourceByNameHandler(rr, req)

	// Check the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := "Source deleted successfully"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

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
	if len(updatedSources) != 1 || updatedSources[0].Name != "TestSource2" {
		t.Errorf("Handler did not update the sources correctly: %v", updatedSources)
	}

	// Validate the directory removal
	sourceDir := filepath.Join("../../../resources/testdata/handlers", "testsource1")
	if _, err := os.Stat(sourceDir); !os.IsNotExist(err) {
		t.Errorf("Handler did not remove the source directory")
	}
}

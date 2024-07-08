package service

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
)

func setupTestData(t *testing.T, resources string, initialSources []source.Source) {
	data, err := json.Marshal(initialSources)
	if err != nil {
		t.Fatalf("Failed to marshal sources: %v", err)
	}
	sourcesFilePath := filepath.Join(resources, "sources.json")
	err = os.WriteFile(sourcesFilePath, data, os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to write test sources.json: %v", err)
	}
	constant.PathToStorage = sourcesFilePath
}

func TestDeleteAndUpdateSources(t *testing.T) {
	tests := []struct {
		name            string
		initialSources  []source.Source
		sourceToDelete  string
		expectedSources []source.Source
		expectError     bool
		errorMessage    string
	}{
		{
			name: "Delete existing source",
			initialSources: []source.Source{
				{Name: "Source1"},
				{Name: "Source2"},
				{Name: "Source3"},
			},
			sourceToDelete: "Source2",
			expectedSources: []source.Source{
				{Name: "Source1"},
				{Name: "Source3"},
			},
			expectError: false,
		},
		{
			name: "Delete non-existing source",
			initialSources: []source.Source{
				{Name: "Source1"},
				{Name: "Source2"},
				{Name: "Source3"},
			},
			sourceToDelete: "NonExistingSource",
			expectedSources: []source.Source{
				{Name: "Source1"},
				{Name: "Source2"},
				{Name: "Source3"},
			},
			expectError:  true,
			errorMessage: "source not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resources := t.TempDir()
			constant.PathToResources = resources

			setupTestData(t, resources, tt.initialSources)

			err := DeleteAndUpdateSources(tt.sourceToDelete)
			if tt.expectError {
				assert.Error(t, err, "Expected error when deleting non-existing source")
			}

			updatedSources := ReadSourcesFromStorage()
			assert.Equal(t, len(tt.expectedSources), len(updatedSources), "Number of sources should be as expected")
			for i, source := range tt.expectedSources {
				assert.Equal(t, source.Name, updatedSources[i].Name, "Source name should match")
			}

			dirPath := filepath.Join(resources, "source2")
			_, err = os.Stat(dirPath)
		})
	}
}

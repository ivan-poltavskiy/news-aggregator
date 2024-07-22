package source

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"news-aggregator/entity/source"
)

func setupTempFile(t *testing.T, content []byte) string {
	tmpDir := os.TempDir()

	tmpFile := filepath.Join(tmpDir, "storage.json")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatal(err)
	}

	return tmpFile
}

func teardownTempFile(filePath string) {
	os.RemoveAll(filepath.Dir(filePath))
}

func TestIsSourceExists(t *testing.T) {
	sources := []source.Source{
		{Name: "Source1"},
		{Name: "Source2"},
	}
	data, err := json.Marshal(sources)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		input    source.Name
		expected bool
	}{
		{
			name:     "source exists",
			input:    "Source1",
			expected: true,
		},
		{
			name:     "source does not exist",
			input:    "Source3",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := setupTempFile(t, data)
			defer teardownTempFile(filePath)

			storage := &jsonSourceStorage{pathToStorage: source.PathToFile(filePath)}
			exists := storage.IsSourceExists(tt.input)
			assert.Equal(t, tt.expected, exists)
		})
	}
}

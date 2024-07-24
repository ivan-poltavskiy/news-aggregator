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

func TestSaveSource(t *testing.T) {
	sources := []source.Source{
		{Name: "Source1"},
	}
	_, err := json.Marshal(sources)
	if err != nil {
		t.Fatal(err)
	}

	newSource := source.Source{Name: "Source2"}

	tests := []struct {
		name      string
		existing  []source.Source
		newSource source.Source
		expected  []source.Source
		expectErr bool
	}{
		{
			name:      "add new source",
			existing:  sources,
			newSource: newSource,
			expected:  append(sources, newSource),
			expectErr: false,
		},
		{
			name:      "source already exists",
			existing:  append(sources, newSource),
			newSource: newSource,
			expected:  append(sources, newSource),
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			existingData, err := json.Marshal(tt.existing)
			if err != nil {
				t.Fatal(err)
			}

			filePath := setupTempFile(t, existingData)
			defer teardownTempFile(filePath)

			storage := &jsonSourceStorage{pathToStorage: source.PathToFile(filePath)}
			err = storage.SaveSource(tt.newSource)
			if (err != nil) != tt.expectErr {
				t.Fatalf("SaveSource() error = %v, wantErr %v", err, tt.expectErr)
			}

			if err != nil {
				return
			}

			sources, err := storage.GetSources()
			if err != nil {
				t.Fatal(err)
			}

			assert.ElementsMatch(t, tt.expected, sources)
		})
	}
}

func TestDeleteSourceByName(t *testing.T) {
	sources := []source.Source{
		{Name: "Source1"},
		{Name: "Source2"},
	}
	data, err := json.Marshal(sources)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		inputName string
		expected  []source.Source
		expectErr bool
	}{
		{
			name:      "delete existing source",
			inputName: "Source1",
			expected:  []source.Source{{Name: "Source2"}},
			expectErr: false,
		},
		{
			name:      "delete non-existing source",
			inputName: "Source3",
			expected:  sources,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := setupTempFile(t, data)
			defer teardownTempFile(filePath)

			storage := &jsonSourceStorage{pathToStorage: source.PathToFile(filePath)}
			err = storage.DeleteSourceByName(source.Name(tt.inputName))
			if (err != nil) != tt.expectErr {
				t.Fatalf("DeleteSourceByName() error = %v, wantErr %v", err, tt.expectErr)
			}

			if err != nil {
				return
			}

			sources, err := storage.GetSources()
			if err != nil {
				t.Fatal(err)
			}

			assert.ElementsMatch(t, tt.expected, sources)
		})
	}
}

func TestGetSourceByName(t *testing.T) {
	sources := []source.Source{
		{Name: "Source1"},
		{Name: "Source2"},
	}
	data, err := json.Marshal(sources)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		inputName source.Name
		expected  source.Source
		expectErr bool
	}{
		{
			name:      "get existing source",
			inputName: "Source1",
			expected:  sources[0],
			expectErr: false,
		},
		{
			name:      "get non-existing source",
			inputName: "Source3",
			expected:  source.Source{},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := setupTempFile(t, data)
			defer teardownTempFile(filePath)

			storage := &jsonSourceStorage{pathToStorage: source.PathToFile(filePath)}
			source, err := storage.GetSourceByName(tt.inputName)
			if (err != nil) != tt.expectErr {
				t.Fatalf("GetSourceByName() error = %v, wantErr %v", err, tt.expectErr)
			}

			if err == nil {
				assert.Equal(t, tt.expected, source)
			}
		})
	}
}

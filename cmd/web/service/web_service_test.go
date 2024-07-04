package service_test

import (
	"encoding/json"
	"news-aggregator/cmd/web/service"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
)

func TestWriteSourcesToFile(t *testing.T) {
	sources := []source.Source{
		{Name: "Source1"},
		{Name: "Source2"},
	}

	resources := t.TempDir()
	constant.PathToStorage = filepath.Join(resources, "sources.json")

	err := service.WriteSourcesToFile(sources)
	assert.NoError(t, err, "Expected no error when writing sources to file")

	file, err := os.Open(constant.PathToStorage)
	assert.NoError(t, err, "Expected no error when opening sources file")
	defer file.Close()

	var readSources []source.Source
	err = json.NewDecoder(file).Decode(&readSources)
	assert.NoError(t, err, "Expected no error when decoding sources file")
	assert.Equal(t, sources, readSources, "Sources should match")
}

func TestReadSourcesFromFile(t *testing.T) {
	sources := []source.Source{
		{Name: "Source1"},
		{Name: "Source2"},
	}

	resources := t.TempDir()
	constant.PathToStorage = filepath.Join(resources, "sources.json")
	data, err := json.Marshal(sources)
	assert.NoError(t, err, "Expected no error when marshaling sources")
	err = os.WriteFile(constant.PathToStorage, data, os.ModePerm)
	assert.NoError(t, err, "Expected no error when writing sources to file")

	readSources := service.ReadSourcesFromFile()
	assert.Equal(t, sources, readSources, "Sources should match")
}

func TestReadSourcesFromFile_NonExistentFile(t *testing.T) {
	resources := t.TempDir()
	constant.PathToStorage = filepath.Join(resources, "sources.json")

	readSources := service.ReadSourcesFromFile()
	assert.Empty(t, readSources, "Expected empty sources when file does not exist")
}
func TestExtractDomainName(t *testing.T) {
	tests := []struct {
		url          string
		expectedName string
	}{
		{"https://www.example.com/path", "example"},
		{"http://example.com/path", "example"},
		{"https://example.com", "example"},
		{"https://sub.example.com/path", "sub"},
		{"invalid-url", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			domain := service.ExtractDomainName(tt.url)
			assert.Equal(t, tt.expectedName, domain, "Expected domain name to match")
		})
	}
}
func TestIsSourceExists(t *testing.T) {
	sources := []source.Source{
		{Name: "Source1"},
		{Name: "Source2"},
	}

	resources := t.TempDir()
	constant.PathToStorage = filepath.Join(resources, "sources.json")
	data, err := json.Marshal(sources)
	assert.NoError(t, err, "Expected no error when marshaling sources")
	err = os.WriteFile(constant.PathToStorage, data, os.ModePerm)
	assert.NoError(t, err, "Expected no error when writing sources to file")

	tests := []struct {
		name     source.Name
		expected bool
	}{
		{"Source1", true},
		{"Source3", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.name), func(t *testing.T) {
			exists := service.IsSourceExists(tt.name)
			assert.Equal(t, tt.expected, exists, "Expected source existence to match")
		})
	}
}
func TestAddSourceToStorage(t *testing.T) {
	initialSources := []source.Source{
		{Name: "Source1"},
		{Name: "Source2"},
	}
	newSource := source.Source{Name: "Source3"}

	resources := t.TempDir()
	constant.PathToStorage = filepath.Join(resources, "sources.json")
	data, err := json.Marshal(initialSources)
	assert.NoError(t, err, "Expected no error when marshaling sources")
	err = os.WriteFile(constant.PathToStorage, data, os.ModePerm)
	assert.NoError(t, err, "Expected no error when writing sources to file")

	service.AddSourceToStorage(newSource)

	file, err := os.Open(constant.PathToStorage)
	assert.NoError(t, err, "Expected no error when opening sources file")
	defer file.Close()

	var readSources []source.Source
	err = json.NewDecoder(file).Decode(&readSources)
	assert.NoError(t, err, "Expected no error when decoding sources file")
	expectedSources := append(initialSources, newSource)
	assert.Equal(t, expectedSources, readSources, "Sources should match")
}

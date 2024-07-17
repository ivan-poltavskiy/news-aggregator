package service

import (
	"encoding/json"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func setupTestEnvironment(t *testing.T) string {
	// Create a temporary directory for resources
	resources := t.TempDir()

	// Set the paths to the temporary directory
	constant.PathToResources = resources
	constant.PathToStorage = filepath.Join(resources, "sources.json")

	// Create sample sources.json
	var sources []source.Source
	data, err := json.Marshal(sources)
	if err != nil {
		t.Fatalf("Failed to marshal sources: %v", err)
	}
	err = os.WriteFile(constant.PathToStorage, data, os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to write test sources.json: %v", err)
	}

	return resources
}

func TestAddSource(t *testing.T) {
	setupTestEnvironment(t)

	tests := []struct {
		name    string
		url     string
		want    source.Name
		wantErr bool
	}{
		{
			name:    "Add pravda rrs source",
			url:     "https://www.pravda.com.ua/",
			want:    source.Name("pravda"),
			wantErr: false,
		},
		{
			name:    "Add not rrs source",
			url:     "https://www.bbc.com",
			want:    "",
			wantErr: true,
		},
		{
			name:    "Add empty source",
			url:     "",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SaveSource(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SaveSource() got = %v, want %v", got, tt.want)
			}
		})
	}
}

package service

import (
	"encoding/json"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
	"os"
	"reflect"
	"testing"
)

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
	err := os.MkdirAll(constant.PathToResources, os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create resources directory: %v", err)
	}

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
}

func cleanupTestEnvironment(t *testing.T) {

	err := os.RemoveAll(constant.PathToResources)
	if err != nil {
		t.Fatalf("Failed to clean up resources directory: %v", err)
	}
}

func TestAddSource(t *testing.T) {
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    source.Name
		wantErr bool
	}{
		{name: "Add pravda rrs source",
			args:    args{url: "https://www.pravda.com.ua/"},
			want:    source.Name("pravda"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddSource(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddSource() got = %v, want %v", got, tt.want)
			}
		})
	}
}

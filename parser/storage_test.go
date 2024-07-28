package parser

import (
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	"os"
	"reflect"
	"testing"
)

// Helper function to create a temporary JSON file for testing
func createTempJSONFile(content string) (string, func(), error) {
	tmpFile, err := os.CreateTemp("", "test_*.json")
	if err != nil {
		return "", nil, err
	}
	if _, err := tmpFile.Write([]byte(content)); err != nil {
		return "", nil, err
	}
	if err := tmpFile.Close(); err != nil {
		return "", nil, err
	}
	return tmpFile.Name(), func() { os.Remove(tmpFile.Name()) }, nil
}

func TestParse(t *testing.T) {
	type args struct {
		path source.PathToFile
		name source.Name
	}

	tests := []struct {
		name    string
		args    args
		content string
		want    []news.News
		wantErr bool
	}{
		{
			name: "Parse valid JSON file",
			args: args{
				path: "test_valid.json",
				name: "testjson",
			},
			content: `[{"title": "Test News 1", "description": "Description 1", "url": "http://example.com/1"},
					  {"title": "Test News 2", "description": "Description 2", "url": "http://example.com/2"}]`,
			want: []news.News{
				{Title: "Test News 1", Description: "Description 1", Link: "http://example.com/1", SourceName: "testjson"},
				{Title: "Test News 2", Description: "Description 2", Link: "http://example.com/2", SourceName: "testjson"},
			},
			wantErr: false,
		},
		{
			name: "File not found",
			args: args{
				path: "non_existent_file.json",
				name: "testjson",
			},
			content: "",
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid JSON format",
			args: args{
				path: "test_invalid.json",
				name: "testjson",
			},
			content: `[{ "title": "Invalid JSON" }`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFilePath, cleanup, err := createTempJSONFile(tt.content)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer cleanup()

			tt.args.path = source.PathToFile(tmpFilePath)

			storage := Storage{}
			got, err := storage.Parse(tt.args.path, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

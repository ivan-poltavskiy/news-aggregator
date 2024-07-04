package service

import (
	"news-aggregator/entity/source"
	"os"
	"reflect"
	"testing"
)

func setupTestEnvironment(t *testing.T) {
	storageDir := "./storage"
	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		err = os.Mkdir(storageDir, os.ModePerm)
		if err != nil {
			t.Fatalf("Failed to create storage directory: %v", err)
		}
	}
}

func cleanupTestEnvironment(t *testing.T) {

	err := os.RemoveAll("resources")
	if err != nil {
		t.Fatalf("Failed to clean up resources directory: %v", err)
	}

	err = os.RemoveAll("./storage")
	if err != nil {
		t.Fatalf("Failed to clean up storage directory: %v", err)
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
		want    source.Source
		wantErr bool
	}{
		{name: "Add pravda source",
			args:    args{url: "https://www.pravda.com.ua/"},
			want:    source.Source{Name: "pravda", PathToFile: "resources\\pravda\\pravda.json", SourceType: "STORAGE"},
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

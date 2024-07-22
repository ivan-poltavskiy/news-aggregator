package service_test

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"news-aggregator/cmd/web/service"
	"news-aggregator/constant"
	"news-aggregator/entity/source"
	"news-aggregator/storage/source/mock_aggregator"
	"os"
	"path/filepath"
	"testing"
)

func setupTestEnvironment(t *testing.T) string {
	// Create a temporary directory for resources
	resources := t.TempDir()

	// Set the paths to the temporary directory
	constant.PathToResources = resources
	constant.PathToStorage = filepath.Join(resources, "sources.json")

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

func TestSaveSource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	setupTestEnvironment(t)

	mockStorage := mock_aggregator.NewMockStorage(ctrl)

	tests := []struct {
		name    string
		url     string
		want    source.Name
		wantErr bool
		setup   func()
	}{
		{
			name:    "Add pravda rrs source",
			url:     "https://www.pravda.com.ua/",
			want:    "pravda",
			wantErr: false,
			setup: func() {
				mockStorage.EXPECT().GetSources().Return([]source.Source{}, nil)
				mockStorage.EXPECT().SaveSource(gomock.AssignableToTypeOf(source.Source{})).Return(nil)
			},
		},
		{
			name:    "Add not rss source",
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
			if tt.setup != nil {
				tt.setup()
			}

			got, err := service.SaveSource(tt.url, mockStorage)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("SaveSource() got = %v, want %v", got, tt.want)
			}
		})
	}
}

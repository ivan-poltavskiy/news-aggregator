package service_test

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"news-aggregator/cmd/web/service"
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	newsStorage_mock "news-aggregator/storage/news/mock_aggregator"
	sourceStorage_mock "news-aggregator/storage/source/mock_aggregator"
)

func TestSaveSource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceStorage := sourceStorage_mock.NewMockStorage(ctrl)
	mockNewsStorage := newsStorage_mock.NewMockNewsStorage(ctrl)

	// Define a path to the test resources directory
	testResourcesDir := "resources"

	// Ensure the test resources directory is removed after tests
	defer func() {
		if err := os.RemoveAll(testResourcesDir); err != nil {
			t.Errorf("Failed to remove test resources directory: %v", err)
		}
	}()

	tests := []struct {
		name    string
		url     string
		want    source.Name
		wantErr bool
		setup   func()
	}{
		{
			name:    "Add valid RSS source",
			url:     "https://www.pravda.com.ua/",
			want:    "pravda",
			wantErr: false,
			setup: func() {
				// Set up expectations for the valid RSS source
				mockSourceStorage.EXPECT().GetSources().Return([]source.Source{}, nil).Times(1)
				mockSourceStorage.EXPECT().SaveSource(gomock.Any()).Return(nil).Times(1)
				mockNewsStorage.EXPECT().GetNews(gomock.Any()).Return([]news.News{
					{
						Title:       "Через ракетну небезпеку у Києві та низці областей оголосили повітряну тривогу",
						Description: "1 липня у Києві та ще низці областей оголосили повітряну тривогу через ракетну небезпеку. \r\nДжерело: Повітряні сили, мапа повітряних тривог \r\nДослівно: \"Увага! Ракетна небезпека для північних областей, де оголошено повітряну тривогу\".",
						Link:        "https://www.pravda.com.ua/news/2024/07/1/7463435/",
						SourceName:  "pravda",
					},
				}, nil).Times(1)
				mockNewsStorage.EXPECT().SaveNews(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},
		{
			name:    "Add invalid RSS source",
			url:     "https://www.test1.com",
			want:    "",
			wantErr: true,
			setup: func() {
			},
		},
		{
			name:    "Add empty source",
			url:     "",
			want:    "",
			wantErr: true,
			setup: func() {
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			got, err := service.SaveSource(tt.url, mockSourceStorage, mockNewsStorage)
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

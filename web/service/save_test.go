package service_test

import (
	"github.com/stretchr/testify/assert"
	"news-aggregator/web/service"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
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

	testResourcesDir := "resources"

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
				returnsNews := []news.News{
					{
						Title:       "Через ракетну небезпеку у Києві та низці областей оголосили повітряну тривогу",
						Description: "1 липня у Києві та ще низці областей оголосили повітряну тривогу через ракетну небезпеку. \r\nДжерело: Повітряні сили, мапа повітряних тривог \r\nДослівно: \"Увага! Ракетна небезпека для північних областей, де оголошено повітряну тривогу\".",
						Link:        "https://www.pravda.com.ua/news/2024/07/1/7463435/",
						SourceName:  "pravda",
					},
				}

				mockSourceStorage.EXPECT().IsSourceExists(gomock.Any()).Return(false).Times(1)
				mockSourceStorage.EXPECT().SaveSource(gomock.Any()).Return(nil).Times(1)
				mockNewsStorage.EXPECT().SaveNews(gomock.Any(), gomock.Any()).Return(
					source.Source{Name: "pravda"},
					nil).Times(1)
				mockNewsStorage.EXPECT().GetNewsBySourceName(gomock.Any(), gomock.Any()).Return(returnsNews, nil).Times(1)
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

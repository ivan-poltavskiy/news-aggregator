package source_test

import (
	"errors"
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	"news-aggregator/mocks"
	sourceService "news-aggregator/web/source"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteSourceByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)

	tests := []struct {
		name       string
		sourceName string
		mockFunc   func()
		expectErr  bool
	}{
		{
			name:       "Success",
			sourceName: "example-source",
			mockFunc: func() {
				mockStorage.EXPECT().DeleteSourceByName(source.Name("example-source")).Return(nil)
			},
			expectErr: false,
		},
		{
			name:       "Failure",
			sourceName: "non-existent-source",
			mockFunc: func() {
				mockStorage.EXPECT().DeleteSourceByName(source.Name("non-existent-source")).Return(errors.New("delete error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			service := sourceService.NewService(mockStorage)
			err := service.DeleteSourceByName(source.Name(tt.sourceName))
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSaveSource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)

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

				mockStorage.EXPECT().IsSourceExists(gomock.Any()).Return(false).Times(1)
				mockStorage.EXPECT().SaveSource(gomock.Any()).Return(nil).Times(1)
				mockStorage.EXPECT().SaveNews(gomock.Any(), gomock.Any()).Return(
					source.Source{Name: "pravda"},
					nil).Times(1)
				mockStorage.EXPECT().GetNewsBySourceName(gomock.Any(), gomock.Any()).Return(returnsNews, nil).Times(1)
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

			service := sourceService.NewService(mockStorage)
			got, err := service.SaveSource(tt.url)
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

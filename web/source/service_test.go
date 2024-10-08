package source_test

import (
	"errors"
	"news-aggregator/entity/news"
	"news-aggregator/entity/source"
	client "news-aggregator/storage/mock_aggregator"
	sourceService "news-aggregator/web/source"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteSourceByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := client.NewMockStorage(ctrl)

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

	mockStorage := client.NewMockStorage(ctrl)

	testResourcesDir := "resources"

	defer func() {
		if err := os.RemoveAll(testResourcesDir); err != nil {
			t.Errorf("Failed to remove test resources directory: %v", err)
		}
	}()

	tests := []struct {
		name       string
		url        string
		sourceName string
		want       source.Name
		wantErr    bool
		setup      func()
	}{
		{
			name:       "Add valid RSS source",
			url:        "https://www.pravda.com.ua/",
			sourceName: "pravda",
			want:       "pravda",
			wantErr:    false,
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
			name:       "Add invalid RSS source",
			url:        "https://www.test1.com",
			sourceName: "test1",
			want:       "",
			wantErr:    true,
			setup: func() {
			},
		},
		{
			name:       "Add empty source",
			sourceName: "",
			url:        "",
			want:       "",
			wantErr:    true,
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
			request := sourceService.AddSourceRequest{Name: tt.sourceName, URL: tt.url}
			got, err := service.SaveSource(request)
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

func TestGetAllSources(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := client.NewMockStorage(ctrl)

	tests := []struct {
		name      string
		mockFunc  func()
		expected  []source.Name
		expectErr bool
	}{
		{
			name: "Success - Get all STORAGE sources",
			mockFunc: func() {
				mockStorage.EXPECT().GetSources().Return([]source.Source{
					{Name: "source1", SourceType: source.STORAGE},
					{Name: "source2", SourceType: source.STORAGE},
					{Name: "source3", SourceType: source.JSON},
				}, nil)
			},
			expected:  []source.Name{"source1", "source2"},
			expectErr: false,
		},
		{
			name: "Failure - Error retrieving sources",
			mockFunc: func() {
				mockStorage.EXPECT().GetSources().Return(nil, errors.New("storage error"))
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name: "Success - No STORAGE sources",
			mockFunc: func() {
				mockStorage.EXPECT().GetSources().Return([]source.Source{
					{Name: "source3", SourceType: source.JSON},
				}, nil)
			},
			expected:  nil,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			service := sourceService.NewService(mockStorage)
			sources, err := service.GetAllSources()

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, sources)
			}
		})
	}
}

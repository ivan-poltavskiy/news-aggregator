package service_test

import (
	"github.com/golang/mock/gomock"
	"news-aggregator/cmd/web/service"
	"news-aggregator/storage/source/mock_aggregator"
	"testing"

	"github.com/stretchr/testify/assert"
	"news-aggregator/entity/source"
)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock_aggregator.NewMockStorage(ctrl)

	sources := []source.Source{
		{Name: "Source1"},
		{Name: "Source2"},
	}

	tests := []struct {
		mockFunc func()
		name     source.Name
		expected bool
	}{
		{
			mockFunc: func() {
				mockStorage.EXPECT().GetSources().Return(sources, nil)
			},
			name:     "Source1",
			expected: true,
		},
		{
			mockFunc: func() {
				mockStorage.EXPECT().GetSources().Return(sources, nil)
			},
			name:     "Source3",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.name), func(t *testing.T) {
			tt.mockFunc()
			exists := service.IsSourceExists(tt.name, mockStorage)
			assert.Equal(t, tt.expected, exists, "Expected source existence to match")
		})
	}
}

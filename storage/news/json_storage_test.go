package news

import (
	"encoding/json"
	"news-aggregator/entity/source"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"news-aggregator/constant"
	"news-aggregator/entity/news"
)

func createTestSource(name string) source.Source {
	return source.Source{
		Name:       source.Name(name),
		PathToFile: source.PathToFile(filepath.Join(constant.PathToResources, name, name+".json")),
	}
}

func TestSaveNews(t *testing.T) {
	tests := []struct {
		name         string
		sourceName   string
		newsArticles []news.News
		expectError  bool
	}{
		{
			name:         "successful save",
			sourceName:   "test_source",
			newsArticles: []news.News{{Title: "Test Article", Description: "Test Description", Link: "http://example.com"}},
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "news-aggregator")
			require.NoError(t, err)
			defer os.RemoveAll(tmpDir)

			constant.PathToResources = tmpDir
			jsonStorage, _ := NewJsonNewsStorage(source.PathToFile(filepath.Join(tmpDir, tt.sourceName)))

			if tt.name == "directory creation error" {
				// Simulate directory creation error
				constant.PathToResources = "/invalid/path"
			}

			currentSource := createTestSource(tt.sourceName)
			_, err = jsonStorage.SaveNews(currentSource, tt.newsArticles)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				_, err := os.Stat(string(currentSource.PathToFile))
				require.NoError(t, err)
			}
		})
	}
}

func TestGetNews(t *testing.T) {
	tests := []struct {
		name           string
		setupFile      func(t *testing.T, filePath string) ([]news.News, error)
		expectError    bool
		expectedLength int
	}{
		{
			name: "successful get",
			setupFile: func(t *testing.T, filePath string) ([]news.News, error) {
				newsArticles := []news.News{
					{Title: "Test Article", Description: "Test Description", Link: "http://example.com"},
				}
				file, err := os.Create(filePath)
				if err != nil {
					return nil, err
				}
				defer file.Close()
				err = json.NewEncoder(file).Encode(newsArticles)
				return newsArticles, err
			},
			expectError:    false,
			expectedLength: 1,
		},
		{
			name: "file does not exist",
			setupFile: func(t *testing.T, filePath string) ([]news.News, error) {
				return nil, nil
			},
			expectError:    false,
			expectedLength: 0,
		},
		{
			name: "invalid JSON format",
			setupFile: func(t *testing.T, filePath string) ([]news.News, error) {
				file, err := os.Create(filePath)
				if err != nil {
					return nil, err
				}
				defer file.Close()
				_, err = file.WriteString("invalid json")
				return nil, err
			},
			expectError:    true,
			expectedLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "news-aggregator")
			require.NoError(t, err)
			defer os.RemoveAll(tmpDir)

			constant.PathToResources = tmpDir
			jsonStorage, _ := NewJsonNewsStorage(source.PathToFile(filepath.Join(tmpDir, "test_source.json")))

			filePath := filepath.Join(tmpDir, "test_source.json")
			_, err = tt.setupFile(t, filePath)
			require.NoError(t, err)

			news, err := jsonStorage.GetNews(filePath)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, news, tt.expectedLength)
			}
		})
	}
}
